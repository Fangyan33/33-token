# API 代理网关模块 Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 实现支持 OpenAI 兼容与 Anthropic 兼容入口的 API 代理网关 MVP，并接入账户校验、路由解析、上游转发与事件记录。

**Architecture:** Envoy AI Gateway 负责对外协议入口与基础流量编排，Go 侧提供准入校验、路由配置读取与结算上报能力。先打通统一执行链路，再补标准化错误映射和基础风控。

**Tech Stack:** Envoy AI Gateway, Go 1.25, PostgreSQL, Docker Compose, Go test

---

## 文件结构

- `gateway/envoy/envoy.yaml`
- `gateway/envoy/routes/openai.yaml`
- `gateway/envoy/routes/anthropic.yaml`
- `services/controlplane/internal/gateway/access/service.go`
- `services/controlplane/internal/gateway/routes/service.go`
- `services/controlplane/internal/gateway/upstream/client.go`
- `services/controlplane/internal/gateway/http/handler.go`
- `tests/integration/gateway/*.go`

### Task 1: 建立统一准入校验接口

**Files:**
- Create: `services/controlplane/internal/gateway/access/service.go`
- Create: `tests/integration/gateway/access_check_test.go`
- Test: `tests/integration/gateway/access_check_test.go`

- [ ] **Step 1: 写准入校验失败测试**

```go
package gateway

import "testing"

func TestCheckAccessRejectsZeroQuota(t *testing.T) {
	t.Fatal("implement access service and replace this placeholder assertion")
}
```

- [ ] **Step 2: 运行测试确认失败**

Run: `go test ./tests/integration/gateway -run TestCheckAccessRejectsZeroQuota -v`
Expected: FAIL

- [ ] **Step 3: 写最小准入校验服务**

```go
package access

type Summary struct {
	QuotaRemaining int64
	Status         string
}

func Check(summary Summary) error {
	if summary.Status != "active" || summary.QuotaRemaining <= 0 {
		return ErrQuotaBlocked
	}
	return nil
}
```

- [ ] **Step 4: 运行测试确认通过**

Run: `go test ./tests/integration/gateway -run TestCheckAccessRejectsZeroQuota -v`
Expected: PASS

- [ ] **Step 5: 提交准入校验**

```bash
git add services/controlplane/internal/gateway/access/service.go tests/integration/gateway/access_check_test.go
git commit -m "feat: add gateway access check service"
```

### Task 2: 落 OpenAI 与 Anthropic 路由入口

**Files:**
- Create: `gateway/envoy/routes/openai.yaml`
- Create: `gateway/envoy/routes/anthropic.yaml`
- Modify: `gateway/envoy/envoy.yaml`
- Test: `tests/integration/gateway/envoy_config_test.go`

- [ ] **Step 1: 写 Envoy 配置失败测试**

```go
package gateway

import "testing"

func TestEnvoyConfigIncludesOpenAIAndAnthropicRoutes(t *testing.T) {
	t.Fatal("assert envoy route config contains both protocol routes")
}
```

- [ ] **Step 2: 运行测试确认失败**

Run: `go test ./tests/integration/gateway -run TestEnvoyConfigIncludesOpenAIAndAnthropicRoutes -v`
Expected: FAIL

- [ ] **Step 3: 写最小路由配置**

```yaml
name: openai-compatible
match:
  path_prefix: /v1
```

```yaml
name: anthropic-compatible
match:
  path_prefix: /anthropic
```

- [ ] **Step 4: 运行测试确认通过**

Run: `go test ./tests/integration/gateway -run TestEnvoyConfigIncludesOpenAIAndAnthropicRoutes -v`
Expected: PASS

- [ ] **Step 5: 提交双协议入口**

```bash
git add gateway/envoy/envoy.yaml gateway/envoy/routes/openai.yaml gateway/envoy/routes/anthropic.yaml tests/integration/gateway/envoy_config_test.go
git commit -m "feat: add gateway protocol routes"
```

### Task 3: 落上游转发与标准化错误

**Files:**
- Create: `services/controlplane/internal/gateway/upstream/client.go`
- Create: `services/controlplane/internal/gateway/http/handler.go`
- Create: `tests/integration/gateway/upstream_proxy_test.go`
- Test: `tests/integration/gateway/upstream_proxy_test.go`

- [ ] **Step 1: 写转发失败测试**

```go
package gateway

import "testing"

func TestProxyMapsUpstreamErrorToStandardResponse(t *testing.T) {
	t.Fatal("assert upstream errors are normalized")
}
```

- [ ] **Step 2: 运行测试确认失败**

Run: `go test ./tests/integration/gateway -run TestProxyMapsUpstreamErrorToStandardResponse -v`
Expected: FAIL

- [ ] **Step 3: 写最小转发客户端与错误映射**

```go
package upstream

type StandardError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func MapError(code string) StandardError {
	return StandardError{Code: code, Message: "upstream request failed"}
}
```

- [ ] **Step 4: 运行测试确认通过**

Run: `go test ./tests/integration/gateway -run TestProxyMapsUpstreamErrorToStandardResponse -v`
Expected: PASS

- [ ] **Step 5: 提交转发与错误映射**

```bash
git add services/controlplane/internal/gateway/upstream/client.go services/controlplane/internal/gateway/http/handler.go tests/integration/gateway/upstream_proxy_test.go
git commit -m "feat: add gateway upstream proxy handling"
```

### Task 4: 加事件记录与风控阻断检查

**Files:**
- Create: `services/controlplane/internal/gateway/events/service.go`
- Create: `tests/integration/gateway/events_settlement_test.go`
- Test: `tests/integration/gateway/events_settlement_test.go`

- [ ] **Step 1: 写 usage 与 billing 事件失败测试**

```go
package gateway

import "testing"

func TestGatewayRecordsUsageAndBillingEvents(t *testing.T) {
	t.Fatal("assert usage and billing events are created together")
}
```

- [ ] **Step 2: 运行测试确认失败**

Run: `go test ./tests/integration/gateway -run TestGatewayRecordsUsageAndBillingEvents -v`
Expected: FAIL

- [ ] **Step 3: 写最小事件记录服务**

```go
package events

type UsageEvent struct {
	RequestID string
	AccountID string
}

type BillingEvent struct {
	IdempotencyKey string
	AccountID      string
}
```

- [ ] **Step 4: 运行测试确认通过**

Run: `go test ./tests/integration/gateway -run TestGatewayRecordsUsageAndBillingEvents -v`
Expected: PASS

- [ ] **Step 5: 提交事件记录基础**

```bash
git add services/controlplane/internal/gateway/events/service.go tests/integration/gateway/events_settlement_test.go
git commit -m "feat: add gateway event recording model"
```
