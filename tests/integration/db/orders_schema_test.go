package db

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestOrdersSchemaIncludesPlanSnapshotAndPaymentEvent(t *testing.T) {
	repoPath := filepath.Join("..", "..", "..", "services", "controlplane", "internal", "store", "orders_repo.go")
	if _, err := os.Stat(repoPath); err != nil {
		t.Fatalf("缺少仓库文件 %s: %v", repoPath, err)
	}

	migrationPath := filepath.Join("..", "..", "..", "db", "migrations", "000003_plans_orders.sql")
	content, err := os.ReadFile(migrationPath)
	if err != nil {
		t.Fatalf("读取迁移文件失败 %s: %v", migrationPath, err)
	}

	sql := string(content)
	for _, needle := range []string{
		"CREATE TABLE plan",
		"CREATE TABLE plan_price_snapshot",
		"CREATE TABLE \"order\"",
		"CREATE TABLE payment_event",
	} {
		if !strings.Contains(sql, needle) {
			t.Fatalf("迁移文件缺少 %q", needle)
		}
	}
}
