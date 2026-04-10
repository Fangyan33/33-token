package orders

import (
	"errors"
	"time"

	"github.com/33-token/model-api-platform/services/controlplane/internal/store"
)

var ErrSnapshotMismatch = errors.New("plan price snapshot does not match input")

type CreateOrderInput struct {
	AccountID           string
	PlanPriceSnapshotID string
	Amount              int64
	Currency            string
	PaymentProvider     string
	OrderType           string
}

type CreatedOrder struct {
	Order    store.Order
	Snapshot store.PlanPriceSnapshot
}

type Service struct{}

func NewService() Service {
	return Service{}
}

func (Service) CreateOrder(input CreateOrderInput, snapshot store.PlanPriceSnapshot) (CreatedOrder, error) {
	if snapshot.ID != input.PlanPriceSnapshotID {
		return CreatedOrder{}, ErrSnapshotMismatch
	}

	now := time.Now().UTC()

	return CreatedOrder{
		Order: store.Order{
			ID:                  "order-bootstrap",
			AccountID:           input.AccountID,
			PlanPriceSnapshotID: input.PlanPriceSnapshotID,
			OrderType:           defaultString(input.OrderType, "one_time_purchase"),
			Status:              "pending_payment",
			PaymentProvider:     defaultString(input.PaymentProvider, "paypal"),
			Amount:              input.Amount,
			Currency:            input.Currency,
			CreatedAt:           now,
		},
		Snapshot: snapshot,
	}, nil
}

func defaultString(value string, fallback string) string {
	if value == "" {
		return fallback
	}
	return value
}
