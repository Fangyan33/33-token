package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/33-token/model-api-platform/services/controlplane/internal/billing/orders"
	"github.com/33-token/model-api-platform/services/controlplane/internal/billing/payments"
	"github.com/33-token/model-api-platform/services/controlplane/internal/billing/summary"
	"github.com/33-token/model-api-platform/services/controlplane/internal/httpx"
)

type AccountQuotaView struct {
	Status         string `json:"status"`
	QuotaTotal     int64  `json:"quotaTotal"`
	QuotaUsed      int64  `json:"quotaUsed"`
	QuotaRemaining int64  `json:"quotaRemaining"`
}

func BuildAccountQuotaView(subscriptionStatus string, current summary.CycleSummary) AccountQuotaView {
	status := current.Status
	if subscriptionStatus != "" && subscriptionStatus != "active" {
		status = subscriptionStatus
	}

	return AccountQuotaView{
		Status:         status,
		QuotaTotal:     current.QuotaTotal,
		QuotaUsed:      current.QuotaUsed,
		QuotaRemaining: current.QuotaRemaining,
	}
}

type Handler struct {
	service *BillingService
}

func NewHandler(service *BillingService) *Handler {
	return &Handler{service: service}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/v1/billing")

	switch {
	case r.Method == http.MethodGet && path == "/plans":
		h.handleListPlans(w, r)
	case r.Method == http.MethodGet && path == "/account-quota":
		h.handleAccountQuota(w, r)
	case r.Method == http.MethodPost && path == "/orders":
		h.handleCreateOrder(w, r)
	case r.Method == http.MethodGet && strings.HasPrefix(path, "/orders/"):
		h.handleGetOrder(w, r)
	case r.Method == http.MethodPost && path == "/payment-events":
		h.handlePaymentEvent(w, r)
	case r.Method == http.MethodPost && path == "/settlements":
		h.handleSettlement(w, r)
	default:
		httpx.WriteJSON(w, http.StatusNotFound, map[string]string{
			"error": "billing endpoint not found",
		})
	}
}

func (h *Handler) handleListPlans(w http.ResponseWriter, _ *http.Request) {
	httpx.WriteJSON(w, http.StatusOK, map[string]any{
		"plans": h.service.ListPlans(),
	})
}

type createOrderRequest struct {
	AccountID           string `json:"accountId"`
	PlanPriceSnapshotID string `json:"planPriceSnapshotId"`
	Amount              int64  `json:"amount"`
	Currency            string `json:"currency"`
	PaymentProvider     string `json:"paymentProvider"`
	OrderType           string `json:"orderType"`
}

type paymentEventRequest struct {
	ProviderEventID string `json:"providerEventId"`
	OrderID         string `json:"orderId"`
	EventType       string `json:"eventType"`
	EventStatus     string `json:"eventStatus"`
}

type settlementRequest struct {
	AccountID    string `json:"accountId"`
	InputTokens  int64  `json:"inputTokens"`
	OutputTokens int64  `json:"outputTokens"`
	TotalTokens  int64  `json:"totalTokens"`
	QuotaDelta   int64  `json:"quotaDelta"`
}

func (h *Handler) handleAccountQuota(w http.ResponseWriter, r *http.Request) {
	accountID := r.URL.Query().Get("account_id")
	if accountID == "" {
		httpx.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "account_id is required"})
		return
	}

	view, err := h.service.QueryAccountQuota(accountID)
	if err != nil {
		httpx.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	httpx.WriteJSON(w, http.StatusOK, view)
}

func (h *Handler) handleCreateOrder(w http.ResponseWriter, r *http.Request) {
	var req createOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid json body"})
		return
	}

	view, err := h.service.CreateOrder(orders.CreateOrderInput{
		AccountID:           req.AccountID,
		PlanPriceSnapshotID: req.PlanPriceSnapshotID,
		Amount:              req.Amount,
		Currency:            req.Currency,
		PaymentProvider:     req.PaymentProvider,
		OrderType:           req.OrderType,
	})
	if err != nil {
		status := http.StatusBadRequest
		if errors.Is(err, ErrSnapshotNotFound) {
			status = http.StatusNotFound
		}
		httpx.WriteJSON(w, status, map[string]string{"error": err.Error()})
		return
	}

	httpx.WriteJSON(w, http.StatusCreated, view)
}

func (h *Handler) handleGetOrder(w http.ResponseWriter, r *http.Request) {
	orderID := strings.TrimPrefix(r.URL.Path, "/api/v1/billing/orders/")
	view, err := h.service.GetOrder(orderID)
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, ErrOrderNotFound) {
			status = http.StatusNotFound
		}
		httpx.WriteJSON(w, status, map[string]string{"error": err.Error()})
		return
	}

	httpx.WriteJSON(w, http.StatusOK, view)
}

func (h *Handler) handlePaymentEvent(w http.ResponseWriter, r *http.Request) {
	var req paymentEventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid json body"})
		return
	}

	orderView, quotaView, err := h.service.ApplyPaymentEvent(payments.PaymentEvent{
		ProviderEventID: req.ProviderEventID,
		OrderID:         req.OrderID,
		EventType:       req.EventType,
		EventStatus:     req.EventStatus,
	})
	if err != nil {
		status := http.StatusBadRequest
		if errors.Is(err, ErrOrderNotFound) {
			status = http.StatusNotFound
		}
		httpx.WriteJSON(w, status, map[string]string{"error": err.Error()})
		return
	}

	httpx.WriteJSON(w, http.StatusOK, map[string]any{
		"order": orderView,
		"quota": quotaView,
	})
}

func (h *Handler) handleSettlement(w http.ResponseWriter, r *http.Request) {
	var req settlementRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid json body"})
		return
	}

	view, err := h.service.SettleUsage(summary.UsageDelta{
		AccountID:    req.AccountID,
		InputTokens:  req.InputTokens,
		OutputTokens: req.OutputTokens,
		TotalTokens:  req.TotalTokens,
		QuotaDelta:   req.QuotaDelta,
	})
	if err != nil {
		status := http.StatusBadRequest
		if errors.Is(err, ErrSummaryNotFound) {
			status = http.StatusNotFound
		}
		httpx.WriteJSON(w, status, map[string]string{"error": err.Error()})
		return
	}

	httpx.WriteJSON(w, http.StatusOK, view)
}
