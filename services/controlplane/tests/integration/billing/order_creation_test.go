package billing

import (
	"testing"

	"github.com/33-token/model-api-platform/services/controlplane/internal/billing/orders"
	"github.com/33-token/model-api-platform/services/controlplane/internal/store"
)

func TestCreateOrderUsesPlanSnapshot(t *testing.T) {
	t.Parallel()

	service := orders.NewService()
	snapshot := store.PlanPriceSnapshot{
		ID:          "snapshot_basic",
		PlanCode:    "basic",
		PlanName:    "Basic",
		PriceAmount: 1900,
		Currency:    "USD",
		QuotaTotal:  1000000,
	}

	created, err := service.CreateOrder(orders.CreateOrderInput{
		AccountID:           "acct_1",
		PlanPriceSnapshotID: "snapshot_basic",
		Amount:              1900,
		Currency:            "USD",
	}, snapshot)
	if err != nil {
		t.Fatalf("create order: %v", err)
	}

	if created.Order.PlanPriceSnapshotID != snapshot.ID {
		t.Fatalf("expected order to persist snapshot id %q, got %q", snapshot.ID, created.Order.PlanPriceSnapshotID)
	}

	if created.Snapshot.PlanCode != "basic" || created.Snapshot.PriceAmount != 1900 {
		t.Fatalf("expected snapshot data to be retained, got %+v", created.Snapshot)
	}

	if created.Order.Status != "pending_payment" {
		t.Fatalf("expected pending_payment status, got %q", created.Order.Status)
	}
}
