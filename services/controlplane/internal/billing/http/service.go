package http

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/33-token/model-api-platform/services/controlplane/internal/billing/orders"
	"github.com/33-token/model-api-platform/services/controlplane/internal/billing/payments"
	"github.com/33-token/model-api-platform/services/controlplane/internal/billing/plans"
	"github.com/33-token/model-api-platform/services/controlplane/internal/billing/subscriptions"
	"github.com/33-token/model-api-platform/services/controlplane/internal/billing/summary"
	"github.com/33-token/model-api-platform/services/controlplane/internal/store"
)

var (
	ErrOrderNotFound     = errors.New("order not found")
	ErrSnapshotNotFound  = errors.New("plan snapshot not found")
	ErrSummaryNotFound   = errors.New("account summary not found")
	ErrPaymentNotSettled = errors.New("payment event is not settled as paid")
)

type OrderView struct {
	ID                  string `json:"id"`
	AccountID           string `json:"accountId"`
	PlanPriceSnapshotID string `json:"planPriceSnapshotId"`
	Status              string `json:"status"`
	PaymentProvider     string `json:"paymentProvider"`
	Amount              int64  `json:"amount"`
	Currency            string `json:"currency"`
}

type PlanView struct {
	ID                string `json:"id"`
	Code              string `json:"code"`
	Name              string `json:"name"`
	BillingPeriodType string `json:"billingPeriodType"`
	QuotaTotal        int64  `json:"quotaTotal"`
	PriceAmount       int64  `json:"priceAmount"`
	Currency          string `json:"currency"`
}

type BillingService struct {
	mu                  sync.Mutex
	planService         plans.Service
	orderService        orders.Service
	paymentService      payments.Service
	subscriptionService *subscriptions.Service
	summaryService      summary.Service
	plans               map[string]store.Plan
	snapshots           map[string]store.PlanPriceSnapshot
	orders              map[string]store.Order
	paymentEvents       map[string]payments.PaymentEvent
	summaries           map[string]summary.CycleSummary
	subscriptions       map[string]string
	nextOrderNumber     int64
}

func NewBillingService() *BillingService {
	return &BillingService{
		planService:         plans.NewService(),
		orderService:        orders.NewService(),
		paymentService:      payments.NewService(),
		subscriptionService: subscriptions.NewService(),
		summaryService:      summary.NewService(),
		plans: map[string]store.Plan{
			"plan_basic": {
				ID:                "plan_basic",
				Code:              "basic",
				Name:              "Basic",
				Status:            "active",
				BillingPeriodType: "monthly",
				QuotaTotal:        1000000,
				CreatedAt:         time.Now().UTC(),
				UpdatedAt:         time.Now().UTC(),
			},
		},
		snapshots: map[string]store.PlanPriceSnapshot{
			"snapshot_basic": {
				ID:                "snapshot_basic",
				PlanID:            "plan_basic",
				PlanCode:          "basic",
				PlanName:          "Basic",
				PriceAmount:       1900,
				Currency:          "USD",
				BillingPeriodType: "monthly",
				QuotaTotal:        1000000,
				CreatedAt:         time.Now().UTC(),
			},
		},
		orders:          map[string]store.Order{},
		paymentEvents:   map[string]payments.PaymentEvent{},
		summaries:       map[string]summary.CycleSummary{},
		subscriptions:   map[string]string{},
		nextOrderNumber: 1,
	}
}

func (s *BillingService) ListPlans() []PlanView {
	s.mu.Lock()
	defer s.mu.Unlock()

	planList := make([]store.Plan, 0, len(s.plans))
	for _, plan := range s.plans {
		planList = append(planList, plan)
	}

	sellable := s.planService.ListSellable(planList)
	views := make([]PlanView, 0, len(sellable))
	for _, plan := range sellable {
		views = append(views, PlanView{
			ID:                plan.ID,
			Code:              plan.Code,
			Name:              plan.Name,
			BillingPeriodType: plan.BillingPeriodType,
			QuotaTotal:        plan.QuotaTotal,
			PriceAmount:       s.snapshotPriceForPlan(plan.ID),
			Currency:          "USD",
		})
	}
	return views
}

func (s *BillingService) CreateOrder(input orders.CreateOrderInput) (OrderView, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	snapshot, ok := s.snapshots[input.PlanPriceSnapshotID]
	if !ok {
		return OrderView{}, ErrSnapshotNotFound
	}

	created, err := s.orderService.CreateOrder(input, snapshot)
	if err != nil {
		return OrderView{}, err
	}

	created.Order.ID = fmt.Sprintf("order_%06d", s.nextOrderNumber)
	s.nextOrderNumber++
	s.orders[created.Order.ID] = created.Order
	return orderToView(created.Order), nil
}

func (s *BillingService) GetOrder(orderID string) (OrderView, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	order, ok := s.orders[orderID]
	if !ok {
		return OrderView{}, ErrOrderNotFound
	}

	return orderToView(order), nil
}

func (s *BillingService) ApplyPaymentEvent(event payments.PaymentEvent) (OrderView, AccountQuotaView, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.paymentService.IsPaid(event) {
		return OrderView{}, AccountQuotaView{}, ErrPaymentNotSettled
	}

	if cached, ok := s.paymentEvents[event.ProviderEventID]; ok {
		order, exists := s.orders[cached.OrderID]
		if !exists {
			return OrderView{}, AccountQuotaView{}, ErrOrderNotFound
		}
		quotaView, err := s.buildAccountQuotaView(order.AccountID)
		if err != nil {
			return OrderView{}, AccountQuotaView{}, err
		}
		return orderToView(order), quotaView, nil
	}

	order, ok := s.orders[event.OrderID]
	if !ok {
		return OrderView{}, AccountQuotaView{}, ErrOrderNotFound
	}

	order.Status = "paid_pending_activation"
	s.orders[order.ID] = order

	result, err := s.subscriptionService.Activate(order.ID)
	if err != nil {
		return OrderView{}, AccountQuotaView{}, err
	}

	snapshot, ok := s.snapshots[order.PlanPriceSnapshotID]
	if !ok {
		return OrderView{}, AccountQuotaView{}, ErrSnapshotNotFound
	}

	if result.Status == "active" {
		order.Status = "completed"
		now := time.Now().UTC()
		order.CompletedAt = &now
		s.orders[order.ID] = order
		s.subscriptions[order.AccountID] = result.Status

		current, exists := s.summaries[order.AccountID]
		if !exists {
			current = summary.CycleSummary{
				AccountID:      order.AccountID,
				QuotaTotal:     snapshot.QuotaTotal,
				QuotaRemaining: snapshot.QuotaTotal,
				Status:         "active",
			}
		} else {
			current.QuotaTotal = snapshot.QuotaTotal
			if current.QuotaRemaining == 0 && current.QuotaUsed == 0 {
				current.QuotaRemaining = snapshot.QuotaTotal
			}
			current.Status = "active"
		}
		s.summaries[order.AccountID] = current
	}
	s.paymentEvents[event.ProviderEventID] = event

	quotaView, err := s.buildAccountQuotaView(order.AccountID)
	if err != nil {
		return OrderView{}, AccountQuotaView{}, err
	}

	return orderToView(order), quotaView, nil
}

func (s *BillingService) QueryAccountQuota(accountID string) (AccountQuotaView, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.buildAccountQuotaView(accountID)
}

func (s *BillingService) SettleUsage(delta summary.UsageDelta) (AccountQuotaView, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	current, ok := s.summaries[delta.AccountID]
	if !ok {
		return AccountQuotaView{}, ErrSummaryNotFound
	}

	result := s.summaryService.Settle(current, delta)
	s.summaries[delta.AccountID] = result.Summary
	s.subscriptions[delta.AccountID] = result.SubscriptionStatus

	return BuildAccountQuotaView(result.SubscriptionStatus, result.Summary), nil
}

func (s *BillingService) buildAccountQuotaView(accountID string) (AccountQuotaView, error) {
	current, ok := s.summaries[accountID]
	if !ok {
		return AccountQuotaView{
			Status:         "inactive",
			QuotaTotal:     0,
			QuotaUsed:      0,
			QuotaRemaining: 0,
		}, nil
	}

	return BuildAccountQuotaView(s.subscriptions[accountID], current), nil
}

func orderToView(order store.Order) OrderView {
	return OrderView{
		ID:                  order.ID,
		AccountID:           order.AccountID,
		PlanPriceSnapshotID: order.PlanPriceSnapshotID,
		Status:              order.Status,
		PaymentProvider:     order.PaymentProvider,
		Amount:              order.Amount,
		Currency:            order.Currency,
	}
}

func (s *BillingService) snapshotPriceForPlan(planID string) int64 {
	for _, snapshot := range s.snapshots {
		if snapshot.PlanID == planID {
			return snapshot.PriceAmount
		}
	}
	return 0
}
