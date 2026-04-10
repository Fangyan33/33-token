package store

import (
	"encoding/json"
	"time"
)

type Plan struct {
	ID                string
	Code              string
	Name              string
	Status            string
	BillingPeriodType string
	QuotaTotal        int64
	RateLimitPolicy   json.RawMessage
	DisplayPriority   int
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

type PlanPriceSnapshot struct {
	ID                      string
	PlanID                  string
	PlanCode                string
	PlanName                string
	PriceAmount             int64
	Currency                string
	BillingPeriodType       string
	QuotaTotal              int64
	RateLimitPolicySnapshot json.RawMessage
	CreatedAt               time.Time
}

type Order struct {
	ID                     string
	AccountID              string
	PlanPriceSnapshotID    string
	OrderType              string
	Status                 string
	PaymentProvider        string
	PaymentProviderOrderID string
	Amount                 int64
	Currency               string
	CreatedAt              time.Time
	PaidAt                 *time.Time
	CompletedAt            *time.Time
}

type PaymentEvent struct {
	ID              string
	OrderID         string
	PaymentProvider string
	ProviderEventID string
	EventType       string
	EventStatus     string
	RawReference    string
	EventOccurredAt time.Time
	ReceivedAt      time.Time
}

type OrdersRepository interface {
	CreateOrder(order Order) (Order, error)
	GetPlanSnapshot(snapshotID string) (PlanPriceSnapshot, error)
	AppendPaymentEvent(event PaymentEvent) (PaymentEvent, error)
}
