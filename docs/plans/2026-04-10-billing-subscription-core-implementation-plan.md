# 计费与订阅核心模块 Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 实现套餐、订单、支付事件、账户权益状态与周期汇总的核心账务能力，并向网关、控制台、后台提供统一账务结论。

**Architecture:** 计费核心作为 Go 领域服务实现，围绕 PostgreSQL 中的订单、支付事件、权益状态和周期汇总表构建。先打通一次性购买自动开通与请求后扣减，再补幂等、补偿与后台修正入口。

**Tech Stack:** Go 1.25, PostgreSQL, sqlc or pgx, Go test

---

## 文件结构

- `services/controlplane/internal/billing/plans/service.go`
- `services/controlplane/internal/billing/orders/service.go`
- `services/controlplane/internal/billing/payments/service.go`
- `services/controlplane/internal/billing/subscriptions/service.go`
- `services/controlplane/internal/billing/summary/service.go`
- `tests/integration/billing/*.go`

### Task 1: 落套餐与订单领域服务

**Files:**
- Create: `services/controlplane/internal/billing/plans/service.go`
- Create: `services/controlplane/internal/billing/orders/service.go`
- Create: `tests/integration/billing/order_creation_test.go`
- Test: `tests/integration/billing/order_creation_test.go`

- [ ] **Step 1: 写创建订单失败测试**

```go
package billing

import "testing"

func TestCreateOrderUsesPlanSnapshot(t *testing.T) {
	t.Fatal("assert order creation persists plan snapshot")
}
```

- [ ] **Step 2: 运行测试确认失败**

Run: `go test ./tests/integration/billing -run TestCreateOrderUsesPlanSnapshot -v`
Expected: FAIL

- [ ] **Step 3: 写最小套餐与订单服务**

```go
package orders

type CreateOrderInput struct {
	AccountID            string
	PlanPriceSnapshotID  string
	Amount               int64
	Currency             string
}
```

- [ ] **Step 4: 运行测试确认通过**

Run: `go test ./tests/integration/billing -run TestCreateOrderUsesPlanSnapshot -v`
Expected: PASS

- [ ] **Step 5: 提交订单基础**

```bash
git add services/controlplane/internal/billing/plans/service.go services/controlplane/internal/billing/orders/service.go tests/integration/billing/order_creation_test.go
git commit -m "feat: add billing order creation service"
```

### Task 2: 落支付事件接入与幂等开通

**Files:**
- Create: `services/controlplane/internal/billing/payments/service.go`
- Create: `services/controlplane/internal/billing/subscriptions/service.go`
- Create: `tests/integration/billing/payment_activation_test.go`
- Test: `tests/integration/billing/payment_activation_test.go`

- [ ] **Step 1: 写支付成功触发开通失败测试**

```go
package billing

import "testing"

func TestPaymentSuccessActivatesSubscriptionOnce(t *testing.T) {
	t.Fatal("assert repeated payment success is idempotent")
}
```

- [ ] **Step 2: 运行测试确认失败**

Run: `go test ./tests/integration/billing -run TestPaymentSuccessActivatesSubscriptionOnce -v`
Expected: FAIL

- [ ] **Step 3: 写最小支付与开通服务**

```go
package payments

type PaymentEvent struct {
	ProviderEventID string
	OrderID         string
	EventType       string
}
```

```go
package subscriptions

func ActivateFromPaidOrder(orderID string) error {
	return nil
}
```

- [ ] **Step 4: 运行测试确认通过**

Run: `go test ./tests/integration/billing -run TestPaymentSuccessActivatesSubscriptionOnce -v`
Expected: PASS

- [ ] **Step 5: 提交支付开通链路**

```bash
git add services/controlplane/internal/billing/payments/service.go services/controlplane/internal/billing/subscriptions/service.go tests/integration/billing/payment_activation_test.go
git commit -m "feat: add billing payment activation flow"
```

### Task 3: 落周期汇总与请求后扣减

**Files:**
- Create: `services/controlplane/internal/billing/summary/service.go`
- Create: `tests/integration/billing/quota_settlement_test.go`
- Test: `tests/integration/billing/quota_settlement_test.go`

- [ ] **Step 1: 写扣减后暂停失败测试**

```go
package billing

import "testing"

func TestSettleUsagePausesAccountWhenQuotaExhausted(t *testing.T) {
	t.Fatal("assert quota exhaustion moves summary and subscription to paused")
}
```

- [ ] **Step 2: 运行测试确认失败**

Run: `go test ./tests/integration/billing -run TestSettleUsagePausesAccountWhenQuotaExhausted -v`
Expected: FAIL

- [ ] **Step 3: 写最小汇总结算服务**

```go
package summary

type UsageDelta struct {
	AccountID  string
	TotalTokens int64
}

func Settle(delta UsageDelta) error {
	return nil
}
```

- [ ] **Step 4: 运行测试确认通过**

Run: `go test ./tests/integration/billing -run TestSettleUsagePausesAccountWhenQuotaExhausted -v`
Expected: PASS

- [ ] **Step 5: 提交周期汇总结算**

```bash
git add services/controlplane/internal/billing/summary/service.go tests/integration/billing/quota_settlement_test.go
git commit -m "feat: add billing quota settlement service"
```

### Task 4: 暴露统一账务查询接口

**Files:**
- Create: `services/controlplane/internal/billing/http/handler.go`
- Create: `tests/integration/billing/query_api_test.go`
- Test: `tests/integration/billing/query_api_test.go`

- [ ] **Step 1: 写账务查询失败测试**

```go
package billing

import "testing"

func TestBillingQueryReturnsAccountStateAndSummary(t *testing.T) {
	t.Fatal("assert account state and current cycle summary are returned together")
}
```

- [ ] **Step 2: 运行测试确认失败**

Run: `go test ./tests/integration/billing -run TestBillingQueryReturnsAccountStateAndSummary -v`
Expected: FAIL

- [ ] **Step 3: 写最小查询处理器**

```go
package http

type AccountQuotaView struct {
	Status         string `json:"status"`
	QuotaRemaining int64  `json:"quotaRemaining"`
}
```

- [ ] **Step 4: 运行测试确认通过**

Run: `go test ./tests/integration/billing -run TestBillingQueryReturnsAccountStateAndSummary -v`
Expected: PASS

- [ ] **Step 5: 提交账务查询接口**

```bash
git add services/controlplane/internal/billing/http/handler.go tests/integration/billing/query_api_test.go
git commit -m "feat: add billing query endpoints"
```
