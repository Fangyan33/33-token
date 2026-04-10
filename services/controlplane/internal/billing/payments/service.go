package payments

// PaymentEvent 表示支付网关回调的最小事件模型。
type PaymentEvent struct {
	ProviderEventID string
	OrderID         string
	EventType       string
}
