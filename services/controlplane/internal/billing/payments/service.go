package payments

import "time"

type PaymentEvent struct {
	ProviderEventID string
	OrderID         string
	EventType       string
	EventStatus     string
	OccurredAt      time.Time
}

type Service struct{}

func NewService() Service {
	return Service{}
}

func (Service) IsPaid(event PaymentEvent) bool {
	return event.EventType == "payment_succeeded" && event.EventStatus == "paid"
}
