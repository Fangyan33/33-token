# 网关与计费核心交互时序 Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 实现网关与计费核心之间的准入校验、请求结算、状态切换与异常补偿链路。

**Architecture:** Go 控制面提供两个清晰内部能力：准入校验能力和请求结算能力。网关在转发前只读周期汇总，转发后统一上报 usage 与 billing 事件，由计费核心更新汇总并在额度耗尽后切换为暂停状态。

**Tech Stack:** Go 1.25, PostgreSQL, Envoy AI Gateway, Go test

---

## 文件结构

- `services/controlplane/internal/gateway/access/service.go`
- `services/controlplane/internal/gateway/events/service.go`
- `services/controlplane/internal/billing/summary/service.go`
- `services/controlplane/internal/billing/http/internal_api.go`
- `tests/integration/gateway_billing/*.go`

### Task 1: 固定准入校验响应契约

**Files:**
- Create: `services/controlplane/internal/billing/http/internal_api.go`
- Create: `tests/integration/gateway_billing/check_quota_contract_test.go`
- Test: `tests/integration/gateway_billing/check_quota_contract_test.go`

- [ ] **Step 1: 写准入契约失败测试**

```go
package gatewaybilling

import "testing"

func TestCheckAccountQuotaContract(t *testing.T) {
	t.Fatal("assert response includes allow, reason, summary status, and remaining quota")
}
```

- [ ] **Step 2: 运行测试确认失败**

Run: `go test ./tests/integration/gateway_billing -run TestCheckAccountQuotaContract -v`
Expected: FAIL

- [ ] **Step 3: 写最小内部准入 API 契约**

```go
package http

type CheckAccountQuotaResponse struct {
	Allow          bool   `json:"allow"`
	Reason         string `json:"reason"`
	Status         string `json:"status"`
	QuotaRemaining int64  `json:"quotaRemaining"`
}
```

- [ ] **Step 4: 运行测试确认通过**

Run: `go test ./tests/integration/gateway_billing -run TestCheckAccountQuotaContract -v`
Expected: PASS

- [ ] **Step 5: 提交准入契约**

```bash
git add services/controlplane/internal/billing/http/internal_api.go tests/integration/gateway_billing/check_quota_contract_test.go
git commit -m "feat: define gateway billing access contract"
```

### Task 2: 固定请求结算输入契约

**Files:**
- Modify: `services/controlplane/internal/billing/http/internal_api.go`
- Create: `tests/integration/gateway_billing/settlement_contract_test.go`
- Test: `tests/integration/gateway_billing/settlement_contract_test.go`

- [ ] **Step 1: 写结算契约失败测试**

```go
package gatewaybilling

import "testing"

func TestRecordUsageAndSettleContract(t *testing.T) {
	t.Fatal("assert settlement contract includes request, tokens, status, and idempotency key")
}
```

- [ ] **Step 2: 运行测试确认失败**

Run: `go test ./tests/integration/gateway_billing -run TestRecordUsageAndSettleContract -v`
Expected: FAIL

- [ ] **Step 3: 写最小结算输入模型**

```go
type RecordUsageAndSettleRequest struct {
	RequestID      string `json:"requestId"`
	IdempotencyKey string `json:"idempotencyKey"`
	AccountID      string `json:"accountId"`
	APIKeyID       string `json:"apiKeyId"`
	ResultStatus   string `json:"resultStatus"`
	TotalTokens    int64  `json:"totalTokens"`
}
```

- [ ] **Step 4: 运行测试确认通过**

Run: `go test ./tests/integration/gateway_billing -run TestRecordUsageAndSettleContract -v`
Expected: PASS

- [ ] **Step 5: 提交结算契约**

```bash
git add services/controlplane/internal/billing/http/internal_api.go tests/integration/gateway_billing/settlement_contract_test.go
git commit -m "feat: define gateway billing settlement contract"
```

### Task 3: 落请求后结算与状态切换

**Files:**
- Modify: `services/controlplane/internal/gateway/events/service.go`
- Modify: `services/controlplane/internal/billing/summary/service.go`
- Create: `tests/integration/gateway_billing/quota_pause_flow_test.go`
- Test: `tests/integration/gateway_billing/quota_pause_flow_test.go`

- [ ] **Step 1: 写额度耗尽暂停失败测试**

```go
package gatewaybilling

import "testing"

func TestSettlementPausesFutureRequestsAfterQuotaExhausted(t *testing.T) {
	t.Fatal("assert one request may exhaust quota and next request is blocked")
}
```

- [ ] **Step 2: 运行测试确认失败**

Run: `go test ./tests/integration/gateway_billing -run TestSettlementPausesFutureRequestsAfterQuotaExhausted -v`
Expected: FAIL

- [ ] **Step 3: 写最小状态切换实现**

```go
package summary

func ApplySettlement(accountID string, totalTokens int64) error {
	return nil
}
```

- [ ] **Step 4: 运行测试确认通过**

Run: `go test ./tests/integration/gateway_billing -run TestSettlementPausesFutureRequestsAfterQuotaExhausted -v`
Expected: PASS

- [ ] **Step 5: 提交结算状态切换**

```bash
git add services/controlplane/internal/gateway/events/service.go services/controlplane/internal/billing/summary/service.go tests/integration/gateway_billing/quota_pause_flow_test.go
git commit -m "feat: add gateway billing settlement flow"
```

### Task 4: 落异常补偿与重复提交保护

**Files:**
- Create: `services/controlplane/internal/billing/compensation/service.go`
- Create: `tests/integration/gateway_billing/compensation_test.go`
- Test: `tests/integration/gateway_billing/compensation_test.go`

- [ ] **Step 1: 写补偿失败测试**

```go
package gatewaybilling

import "testing"

func TestCompensationReplaysPendingSettlementSafely(t *testing.T) {
	t.Fatal("assert repeated settlement uses idempotency key")
}
```

- [ ] **Step 2: 运行测试确认失败**

Run: `go test ./tests/integration/gateway_billing -run TestCompensationReplaysPendingSettlementSafely -v`
Expected: FAIL

- [ ] **Step 3: 写最小补偿服务**

```go
package compensation

func ReplayPendingSettlement(idempotencyKey string) error {
	return nil
}
```

- [ ] **Step 4: 运行测试确认通过**

Run: `go test ./tests/integration/gateway_billing -run TestCompensationReplaysPendingSettlementSafely -v`
Expected: PASS

- [ ] **Step 5: 提交补偿层**

```bash
git add services/controlplane/internal/billing/compensation/service.go tests/integration/gateway_billing/compensation_test.go
git commit -m "feat: add gateway billing compensation flow"
```
