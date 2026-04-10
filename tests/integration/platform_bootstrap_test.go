package integration

import (
	"os"
	"path/filepath"
	"testing"
)

func TestPlatformBootstrapChecklist(t *testing.T) {
	files := []string{
		"package.json",
		"pnpm-workspace.yaml",
		"apps/web/package.json",
		"apps/web/src/main.tsx",
		"services/controlplane/go.mod",
		"services/controlplane/cmd/api/main.go",
		"gateway/envoy/docker-compose.yaml",
	}

	for _, file := range files {
		path := filepath.Join("..", "..", file)
		if _, err := os.Stat(path); err != nil {
			t.Fatalf("缺少骨架文件 %s: %v", file, err)
		}
	}
}
