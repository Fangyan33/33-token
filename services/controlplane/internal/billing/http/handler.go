package http

// AccountQuotaView 表示统一账务查询接口的最小响应视图。
type AccountQuotaView struct {
	Status         string `json:"status"`
	QuotaRemaining int64  `json:"quotaRemaining"`
}
