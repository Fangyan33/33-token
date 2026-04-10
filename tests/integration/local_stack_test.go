package integration

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLocalStackFilesExist(t *testing.T) {
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
}
