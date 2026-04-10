package integration

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLocalStackChecklist(t *testing.T) {
	files := []string{
		"docker-compose.yaml",
		"gateway/envoy/envoy.yaml",
		"services/controlplane/.env.example",
		"db/migrations/000001_init_extensions.sql",
	}

	for _, file := range files {
		path := filepath.Join("..", "..", file)
		if _, err := os.Stat(path); err != nil {
			t.Fatalf("缺少本地栈文件 %s: %v", file, err)
		}
	}

	compose := mustRead(t, filepath.Join("..", "..", "docker-compose.yaml"))
	mustContain(t, compose, "postgres:")
	mustContain(t, compose, "api:")
	mustContain(t, compose, "envoy:")
	mustContain(t, compose, "postgres:17-alpine")
	mustContain(t, compose, "go run ./cmd/api")

	envoy := mustRead(t, filepath.Join("..", "..", "gateway/envoy/envoy.yaml"))
	mustContain(t, envoy, "name: envoy")
	mustContain(t, envoy, "direct_response")

	migration := mustRead(t, filepath.Join("..", "..", "db/migrations/000001_init_extensions.sql"))
	mustContain(t, migration, `CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`)
}

func mustRead(t *testing.T, path string) string {
	t.Helper()

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("读取文件失败 %s: %v", path, err)
	}

	return string(data)
}

func mustContain(t *testing.T, content, needle string) {
	t.Helper()

	if !strings.Contains(content, needle) {
		t.Fatalf("内容缺少 %q", needle)
	}
}
