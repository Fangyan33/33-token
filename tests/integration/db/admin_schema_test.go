package db

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestAdminSchemaIncludesModelRouteCredentialRefAndAuditLog(t *testing.T) {
	repoPath := filepath.Join("..", "..", "..", "services", "controlplane", "internal", "store", "admin_repo.go")
	if _, err := os.Stat(repoPath); err != nil {
		t.Fatalf("缺少仓库文件 %s: %v", repoPath, err)
	}

	migrationPath := filepath.Join("..", "..", "..", "db", "migrations", "000006_admin_config.sql")
	content, err := os.ReadFile(migrationPath)
	if err != nil {
		t.Fatalf("读取迁移文件失败 %s: %v", migrationPath, err)
	}

	sql := string(content)
	for _, needle := range []string{
		"CREATE TABLE upstream_credential_ref",
		"CREATE TABLE model_route",
		"CREATE TABLE admin_audit_log",
	} {
		if !strings.Contains(sql, needle) {
			t.Fatalf("迁移文件缺少 %q", needle)
		}
	}
}
