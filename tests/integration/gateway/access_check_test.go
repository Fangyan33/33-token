package gateway

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestCheckAccessRejectsZeroQuota(t *testing.T) {
	controlplaneRoot := filepath.Join("..", "..", "..", "services", "controlplane")
	tempDir, err := os.MkdirTemp(controlplaneRoot, "gateway-access-check-*")
	if err != nil {
		t.Fatalf("创建临时测试目录失败: %v", err)
	}
	t.Cleanup(func() {
		_ = os.RemoveAll(tempDir)
	})

	testSource := []byte(`package accesscheck

import (
	"errors"
	"testing"

	"platform.local/controlplane/internal/gateway/access"
)

func TestAccessCheckRejectsZeroQuota(t *testing.T) {
	err := access.Check(access.Summary{Status: "active", QuotaRemaining: 0})
	if !errors.Is(err, access.ErrQuotaBlocked) {
		t.Fatalf("expected quota blocked error, got %v", err)
	}
}

func TestAccessCheckRejectsInactiveStatus(t *testing.T) {
	err := access.Check(access.Summary{Status: "suspended", QuotaRemaining: 10})
	if !errors.Is(err, access.ErrQuotaBlocked) {
		t.Fatalf("expected quota blocked error, got %v", err)
	}
}
`)

	testFile := filepath.Join(tempDir, "access_check_test.go")
	if err := os.WriteFile(testFile, testSource, 0o600); err != nil {
		t.Fatalf("写入临时测试文件失败: %v", err)
	}

	cmd := exec.Command("go", "test", "-run", "TestAccessCheckRejects(ZeroQuota|InactiveStatus)", "-v", ".")
	cmd.Dir = tempDir
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("临时行为测试执行失败: %v\n输出:\n%s", err, string(output))
	}

	if !strings.Contains(string(output), "PASS") {
		t.Fatalf("临时行为测试未通过，输出:\n%s", string(output))
	}
}
