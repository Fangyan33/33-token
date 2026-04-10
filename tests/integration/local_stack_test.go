package integration

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLocalStackFilesExist(t *testing.T) {
	t.Parallel()

	requiredPaths := []string{
		filepath.Join("..", "..", "docker-compose.yaml"),
		filepath.Join("..", "..", "gateway", "envoy", "envoy.yaml"),
		filepath.Join("..", "..", "services", "controlplane", ".env.example"),
		filepath.Join("..", "..", "db", "migrations", "000001_init_extensions.sql"),
	}

	for _, path := range requiredPaths {
		if _, err := os.Stat(path); err != nil {
			t.Fatalf("expected local stack path %s to exist: %v", path, err)
		}
	}
}
