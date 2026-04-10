package billing

import (
	"testing"

	"github.com/33-token/model-api-platform/services/controlplane/internal/billing/summary"
)

func TestSettleUsagePausesAccountWhenQuotaExhausted(t *testing.T) {
	t.Parallel()

	service := summary.NewService()
	result := service.Settle(summary.CycleSummary{
		AccountID:      "acct_1",
		QuotaTotal:     100,
		QuotaUsed:      90,
		QuotaRemaining: 10,
		Status:         "active",
	}, summary.UsageDelta{
		AccountID:    "acct_1",
		InputTokens:  5,
		OutputTokens: 10,
		TotalTokens:  15,
		QuotaDelta:   15,
	})

	if result.Summary.QuotaUsed != 105 {
		t.Fatalf("expected used quota to be 105, got %d", result.Summary.QuotaUsed)
	}

	if result.Summary.QuotaRemaining != 0 {
		t.Fatalf("expected remaining quota to be 0, got %d", result.Summary.QuotaRemaining)
	}

	if result.Summary.Status != "paused" || result.SubscriptionStatus != "paused" {
		t.Fatalf("expected paused statuses, got summary=%q subscription=%q", result.Summary.Status, result.SubscriptionStatus)
	}
}
