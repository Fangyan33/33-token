package db

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestAdminSchemaIncludesModelRouteCredentialRefAndAuditLog(t *testing.T) {
	t.Parallel()

	content, err := os.ReadFile(filepath.Join("..", "..", "..", "db", "migrations", "000006_admin_config.sql"))
	if err != nil {
		t.Fatalf("read migration: %v", err)
	}

	sql := string(content)
	for _, expected := range []string{
		"create table if not exists upstream_credential_ref",
		"create table if not exists model_route",
		"upstream_credential_ref_id uuid not null references upstream_credential_ref(id)",
		"create table if not exists admin_audit_log",
	} {
		if !strings.Contains(sql, expected) {
			t.Fatalf("expected migration to contain %q", expected)
		}
	}
}
