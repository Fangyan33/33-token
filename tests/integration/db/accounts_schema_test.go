package db

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestAccountsSchemaIncludesAccountUserIdentityAndAPIKey(t *testing.T) {
	t.Parallel()

	content, err := os.ReadFile(filepath.Join("..", "..", "..", "db", "migrations", "000002_accounts.sql"))
	if err != nil {
		t.Fatalf("read migration: %v", err)
	}

	sql := string(content)
	for _, expected := range []string{
		"create table if not exists account",
		"create table if not exists user_identity",
		"create table if not exists api_key",
		"key_hash text not null",
	} {
		if !strings.Contains(sql, expected) {
			t.Fatalf("expected migration to contain %q", expected)
		}
	}
}
