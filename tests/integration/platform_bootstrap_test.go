package integration

import (
	"os"
	"path/filepath"
	"testing"
)

func TestPlatformBootstrapChecklist(t *testing.T) {
	t.Parallel()

	requiredPaths := []string{
		filepath.Join("..", "..", "package.json"),
		filepath.Join("..", "..", "pnpm-workspace.yaml"),
		filepath.Join("..", "..", "apps", "web", "src", "main.tsx"),
		filepath.Join("..", "..", "services", "controlplane", "go.mod"),
		filepath.Join("..", "..", "services", "controlplane", "cmd", "api", "main.go"),
		filepath.Join("..", "..", "gateway", "envoy", "docker-compose.yaml"),
	}

	for _, path := range requiredPaths {
		if _, err := os.Stat(path); err != nil {
			t.Fatalf("expected path %s to exist: %v", path, err)
		}
	}
}
