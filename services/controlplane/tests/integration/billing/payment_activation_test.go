package billing

import (
	"testing"

	"github.com/33-token/model-api-platform/services/controlplane/internal/billing/subscriptions"
)

func TestPaymentSuccessActivatesSubscriptionOnce(t *testing.T) {
	t.Parallel()

	service := subscriptions.NewService()

	first, err := service.Activate("order_1")
	if err != nil {
		t.Fatalf("first activation failed: %v", err)
	}

	second, err := service.Activate("order_1")
	if err != nil {
		t.Fatalf("second activation failed: %v", err)
	}

	if !first.Activated {
		t.Fatal("expected first activation to activate subscription")
	}

	if second.Activated {
		t.Fatal("expected repeated activation to be idempotent")
	}

	if second.Status != "active" {
		t.Fatalf("expected repeated activation to keep active status, got %q", second.Status)
	}
}
