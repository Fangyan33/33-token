package db

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestOrdersSchemaIncludesPlanSnapshotAndPaymentEvent(t *testing.T) {
	t.Parallel()

	content, err := os.ReadFile(filepath.Join("..", "..", "..", "db", "migrations", "000003_plans_orders.sql"))
	if err != nil {
		t.Fatalf("read migration: %v", err)
	}

	sql := string(content)
	for _, expected := range []string{
		"create table if not exists plan",
		"create table if not exists plan_price_snapshot",
		"create table if not exists \"order\"",
		"create table if not exists payment_event",
		"plan_price_snapshot_id uuid not null references plan_price_snapshot(id)",
	} {
		if !strings.Contains(sql, expected) {
			t.Fatalf("expected migration to contain %q", expected)
		}
	}
}
