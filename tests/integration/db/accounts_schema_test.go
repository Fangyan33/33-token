package db

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestAccountsSchemaFilesAndTablesExist(t *testing.T) {
	repoPath := filepath.Join("..", "..", "..", "services", "controlplane", "internal", "store", "accounts_repo.go")
	if _, err := os.Stat(repoPath); err != nil {
		t.Fatalf("缺少仓库文件 %s: %v", repoPath, err)
	}

	migrationPath := filepath.Join("..", "..", "..", "db", "migrations", "000002_accounts.sql")
	content, err := os.ReadFile(migrationPath)
	if err != nil {
		t.Fatalf("读取迁移文件失败 %s: %v", migrationPath, err)
	}

	sql := string(content)
	for _, needle := range []string{
		"CREATE TABLE account",
		"CREATE TABLE user_identity",
		"CREATE TABLE api_key",
	} {
		if !strings.Contains(sql, needle) {
			t.Fatalf("迁移文件缺少 %q", needle)
		}
	}
}
