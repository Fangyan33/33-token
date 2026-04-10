package summary

// UsageDelta 表示一次请求后的最小用量变更。
type UsageDelta struct {
	AccountID   string
	TotalTokens int64
}

// Settle 是周期汇总结算的最小骨架。
func Settle(delta UsageDelta) error {
	return nil
}
