package billing

import (
	"testing"

	billinghttp "github.com/33-token/model-api-platform/services/controlplane/internal/billing/http"
	"github.com/33-token/model-api-platform/services/controlplane/internal/billing/summary"
)

func TestBillingQueryReturnsAccountStateAndSummary(t *testing.T) {
	t.Parallel()

	view := billinghttp.BuildAccountQuotaView("active", summary.CycleSummary{
		AccountID:      "acct_1",
		QuotaTotal:     100,
		QuotaUsed:      20,
		QuotaRemaining: 80,
		Status:         "active",
	})

	if view.Status != "active" {
		t.Fatalf("expected active status, got %q", view.Status)
	}

	if view.QuotaRemaining != 80 || view.QuotaUsed != 20 || view.QuotaTotal != 100 {
		t.Fatalf("unexpected quota view: %+v", view)
	}
}
