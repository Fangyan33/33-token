package gateway

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestEnvoyConfigIncludesOpenAIAndAnthropicRoutes(t *testing.T) {
	envoyRoot := filepath.Join("..", "..", "..", "gateway", "envoy")

	envoyConfig, err := os.ReadFile(filepath.Join(envoyRoot, "envoy.yaml"))
	if err != nil {
		t.Fatalf("读取 envoy 主配置失败: %v", err)
	}

	openAIRoute, err := os.ReadFile(filepath.Join(envoyRoot, "routes", "openai.yaml"))
	if err != nil {
		t.Fatalf("读取 OpenAI 路由配置失败: %v", err)
	}

	anthropicRoute, err := os.ReadFile(filepath.Join(envoyRoot, "routes", "anthropic.yaml"))
	if err != nil {
		t.Fatalf("读取 Anthropic 路由配置失败: %v", err)
	}

	if !strings.Contains(string(envoyConfig), "openai-compatible") {
		t.Fatal("envoy 主配置未包含 openai-compatible 引用")
	}

	if !strings.Contains(string(envoyConfig), "anthropic-compatible") {
		t.Fatal("envoy 主配置未包含 anthropic-compatible 引用")
	}

	if !strings.Contains(string(envoyConfig), "gateway/envoy/routes/openai.yaml") {
		t.Fatal("envoy 主配置未引用 routes/openai.yaml")
	}

	if !strings.Contains(string(envoyConfig), "gateway/envoy/routes/anthropic.yaml") {
		t.Fatal("envoy 主配置未引用 routes/anthropic.yaml")
	}

	if !strings.Contains(string(envoyConfig), "/v1") {
		t.Fatal("envoy 主配置未包含 /v1 入口提示")
	}

	if !strings.Contains(string(envoyConfig), "/anthropic") {
		t.Fatal("envoy 主配置未包含 /anthropic 入口提示")
	}

	if !strings.Contains(string(openAIRoute), "path_prefix: /v1") {
		t.Fatal("OpenAI 路由未使用 /v1 前缀")
	}

	if !strings.Contains(string(anthropicRoute), "path_prefix: /anthropic") {
		t.Fatal("Anthropic 路由未使用 /anthropic 前缀")
	}
}
