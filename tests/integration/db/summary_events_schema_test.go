package db

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSummaryAndEventsSchemaSupportRealtimeQuotaChecks(t *testing.T) {
	repoPath := filepath.Join("..", "..", "..", "services", "controlplane", "internal", "store", "summary_repo.go")
	if _, err := os.Stat(repoPath); err != nil {
		t.Fatalf("缺少仓库文件 %s: %v", repoPath, err)
	}

	migrationPaths := []string{
		filepath.Join("..", "..", "..", "db", "migrations", "000004_subscription_summary.sql"),
		filepath.Join("..", "..", "..", "db", "migrations", "000005_gateway_events.sql"),
	}

	combinedSQL := ""
	for _, migrationPath := range migrationPaths {
		content, err := os.ReadFile(migrationPath)
		if err != nil {
			t.Fatalf("读取迁移文件失败 %s: %v", migrationPath, err)
		}
		combinedSQL += string(content)
	}

	for _, needle := range []string{
		"CREATE TABLE account_subscription_state",
		"CREATE TABLE account_cycle_summary",
		"CREATE TABLE usage_event",
		"CREATE TABLE billing_event",
	} {
		if !strings.Contains(combinedSQL, needle) {
			t.Fatalf("迁移文件缺少 %q", needle)
		}
	}
}
