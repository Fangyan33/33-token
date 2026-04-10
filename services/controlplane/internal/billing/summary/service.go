package summary

type UsageDelta struct {
	AccountID    string
	InputTokens  int64
	OutputTokens int64
	TotalTokens  int64
	QuotaDelta   int64
}

type CycleSummary struct {
	AccountID         string
	QuotaTotal        int64
	QuotaUsed         int64
	QuotaRemaining    int64
	InputTokensTotal  int64
	OutputTokensTotal int64
	TotalTokensTotal  int64
	Status            string
}

type SettlementResult struct {
	Summary            CycleSummary
	SubscriptionStatus string
}

type Service struct{}

func NewService() Service {
	return Service{}
}

func (Service) Settle(current CycleSummary, delta UsageDelta) SettlementResult {
	quotaDelta := delta.QuotaDelta
	if quotaDelta == 0 {
		quotaDelta = delta.TotalTokens
	}

	next := current
	next.QuotaUsed += quotaDelta
	next.QuotaRemaining -= quotaDelta
	if next.QuotaRemaining < 0 {
		next.QuotaRemaining = 0
	}

	next.InputTokensTotal += delta.InputTokens
	next.OutputTokensTotal += delta.OutputTokens
	next.TotalTokensTotal += delta.TotalTokens
	next.Status = "active"

	subscriptionStatus := "active"
	if next.QuotaRemaining <= 0 {
		next.Status = "paused"
		subscriptionStatus = "paused"
	}

	return SettlementResult{
		Summary:            next,
		SubscriptionStatus: subscriptionStatus,
	}
}
