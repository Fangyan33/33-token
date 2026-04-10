package db

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSummaryAndEventsSchemaSupportRealtimeQuotaChecks(t *testing.T) {
	t.Parallel()

	summaryContent, err := os.ReadFile(filepath.Join("..", "..", "..", "db", "migrations", "000004_subscription_summary.sql"))
	if err != nil {
		t.Fatalf("read summary migration: %v", err)
	}

	eventsContent, err := os.ReadFile(filepath.Join("..", "..", "..", "db", "migrations", "000005_gateway_events.sql"))
	if err != nil {
		t.Fatalf("read events migration: %v", err)
	}

	sql := string(summaryContent) + "\n" + string(eventsContent)
	for _, expected := range []string{
		"create table if not exists account_subscription_state",
		"create table if not exists account_cycle_summary",
		"quota_remaining bigint not null",
		"create table if not exists usage_event",
		"create table if not exists billing_event",
		"idempotency_key text not null unique",
	} {
		if !strings.Contains(sql, expected) {
			t.Fatalf("expected schema to contain %q", expected)
		}
	}
}
