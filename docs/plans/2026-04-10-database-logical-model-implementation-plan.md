# 数据库逻辑模型 Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 按数据库逻辑模型设计落 PostgreSQL schema、迁移脚本、基础仓储和关键查询路径。

**Architecture:** 按账户与认证、套餐与订单、权益与周期汇总、网关事件、运营配置与审计五个逻辑域拆分 migration 与 repository。优先实现支持实时准入和请求结算的核心表，再补运营治理与查询辅助结构。

**Tech Stack:** PostgreSQL, Go 1.25, dbmate or goose, pgx, Go test

---

## 文件结构

- `db/migrations/000002_accounts.sql`
- `db/migrations/000003_plans_orders.sql`
- `db/migrations/000004_subscription_summary.sql`
- `db/migrations/000005_gateway_events.sql`
- `db/migrations/000006_admin_config.sql`
- `services/controlplane/internal/store/*`
- `tests/integration/db/*.go`

### Task 1: 落账户与认证域 schema

**Files:**
- Create: `db/migrations/000002_accounts.sql`
- Create: `services/controlplane/internal/store/accounts_repo.go`
- Create: `tests/integration/db/accounts_schema_test.go`
- Test: `tests/integration/db/accounts_schema_test.go`

- [ ] **Step 1: 写账户域失败测试**

```go
package db

import "testing"

func TestAccountsSchemaIncludesAccountUserIdentityAndAPIKey(t *testing.T) {
	t.Fatal("assert accounts migration creates account, user_identity, and api_key tables")
}
```

- [ ] **Step 2: 运行测试确认失败**

Run: `go test ./tests/integration/db -run TestAccountsSchemaIncludesAccountUserIdentityAndAPIKey -v`
Expected: FAIL

- [ ] **Step 3: 写最小账户域 migration**

```sql
create table account (
  id uuid primary key,
  status text not null
);

create table user_identity (
  id uuid primary key,
  account_id uuid not null references account(id)
);

create table api_key (
  id uuid primary key,
  account_id uuid not null references account(id),
  key_hash text not null
);
```

- [ ] **Step 4: 运行测试确认通过**

Run: `go test ./tests/integration/db -run TestAccountsSchemaIncludesAccountUserIdentityAndAPIKey -v`
Expected: PASS

- [ ] **Step 5: 提交账户域 schema**

```bash
git add db/migrations/000002_accounts.sql services/controlplane/internal/store/accounts_repo.go tests/integration/db/accounts_schema_test.go
git commit -m "feat: add accounts schema"
```

### Task 2: 落套餐与订单域 schema

**Files:**
- Create: `db/migrations/000003_plans_orders.sql`
- Create: `services/controlplane/internal/store/orders_repo.go`
- Create: `tests/integration/db/orders_schema_test.go`
- Test: `tests/integration/db/orders_schema_test.go`

- [ ] **Step 1: 写订单域失败测试**

```go
package db

import "testing"

func TestOrdersSchemaIncludesPlanSnapshotAndPaymentEvent(t *testing.T) {
	t.Fatal("assert plan, plan_price_snapshot, order, and payment_event tables exist")
}
```

- [ ] **Step 2: 运行测试确认失败**

Run: `go test ./tests/integration/db -run TestOrdersSchemaIncludesPlanSnapshotAndPaymentEvent -v`
Expected: FAIL

- [ ] **Step 3: 写最小订单域 migration**

```sql
create table plan (
  id uuid primary key,
  code text not null unique
);

create table plan_price_snapshot (
  id uuid primary key,
  plan_id uuid not null references plan(id)
);

create table "order" (
  id uuid primary key,
  account_id uuid not null references account(id),
  plan_price_snapshot_id uuid not null references plan_price_snapshot(id)
);

create table payment_event (
  id uuid primary key,
  order_id uuid not null references "order"(id)
);
```

- [ ] **Step 4: 运行测试确认通过**

Run: `go test ./tests/integration/db -run TestOrdersSchemaIncludesPlanSnapshotAndPaymentEvent -v`
Expected: PASS

- [ ] **Step 5: 提交订单域 schema**

```bash
git add db/migrations/000003_plans_orders.sql services/controlplane/internal/store/orders_repo.go tests/integration/db/orders_schema_test.go
git commit -m "feat: add plans and orders schema"
```

### Task 3: 落权益、汇总与事件域 schema

**Files:**
- Create: `db/migrations/000004_subscription_summary.sql`
- Create: `db/migrations/000005_gateway_events.sql`
- Create: `services/controlplane/internal/store/summary_repo.go`
- Create: `tests/integration/db/summary_events_schema_test.go`
- Test: `tests/integration/db/summary_events_schema_test.go`

- [ ] **Step 1: 写权益与事件域失败测试**

```go
package db

import "testing"

func TestSummaryAndEventsSchemaSupportRealtimeQuotaChecks(t *testing.T) {
	t.Fatal("assert subscription, cycle summary, usage_event, and billing_event tables exist")
}
```

- [ ] **Step 2: 运行测试确认失败**

Run: `go test ./tests/integration/db -run TestSummaryAndEventsSchemaSupportRealtimeQuotaChecks -v`
Expected: FAIL

- [ ] **Step 3: 写最小权益与事件域 migration**

```sql
create table account_subscription_state (
  id uuid primary key,
  account_id uuid not null references account(id),
  status text not null
);

create table account_cycle_summary (
  id uuid primary key,
  account_id uuid not null references account(id),
  quota_remaining bigint not null
);

create table usage_event (
  id uuid primary key,
  account_id uuid not null references account(id)
);

create table billing_event (
  id uuid primary key,
  account_id uuid not null references account(id),
  usage_event_id uuid not null references usage_event(id)
);
```

- [ ] **Step 4: 运行测试确认通过**

Run: `go test ./tests/integration/db -run TestSummaryAndEventsSchemaSupportRealtimeQuotaChecks -v`
Expected: PASS

- [ ] **Step 5: 提交权益与事件域 schema**

```bash
git add db/migrations/000004_subscription_summary.sql db/migrations/000005_gateway_events.sql services/controlplane/internal/store/summary_repo.go tests/integration/db/summary_events_schema_test.go
git commit -m "feat: add subscription and gateway event schema"
```

### Task 4: 落运营配置与后台审计域 schema

**Files:**
- Create: `db/migrations/000006_admin_config.sql`
- Create: `services/controlplane/internal/store/admin_repo.go`
- Create: `tests/integration/db/admin_schema_test.go`
- Test: `tests/integration/db/admin_schema_test.go`

- [ ] **Step 1: 写后台配置域失败测试**

```go
package db

import "testing"

func TestAdminSchemaIncludesModelRouteCredentialRefAndAuditLog(t *testing.T) {
	t.Fatal("assert model_route, upstream_credential_ref, and admin_audit_log exist")
}
```

- [ ] **Step 2: 运行测试确认失败**

Run: `go test ./tests/integration/db -run TestAdminSchemaIncludesModelRouteCredentialRefAndAuditLog -v`
Expected: FAIL

- [ ] **Step 3: 写最小后台配置 migration**

```sql
create table upstream_credential_ref (
  id uuid primary key,
  credential_key text not null
);

create table model_route (
  id uuid primary key,
  upstream_credential_ref_id uuid not null references upstream_credential_ref(id)
);

create table admin_audit_log (
  id uuid primary key,
  operator_id uuid not null,
  action_type text not null
);
```

- [ ] **Step 4: 运行测试确认通过**

Run: `go test ./tests/integration/db -run TestAdminSchemaIncludesModelRouteCredentialRefAndAuditLog -v`
Expected: PASS

- [ ] **Step 5: 提交后台配置域 schema**

```bash
git add db/migrations/000006_admin_config.sql services/controlplane/internal/store/admin_repo.go tests/integration/db/admin_schema_test.go
git commit -m "feat: add admin config schema"
```
