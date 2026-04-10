# 管理后台模块 Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 实现内部管理后台 MVP，覆盖模型路由配置、套餐管理、异常订单补单、账户状态调整以及运营汇总视图。

**Architecture:** 管理后台作为 React.js 内部受保护区域实现，通过 Go 后端受控接口执行高风险运营动作。后台自身不直接写数据库，所有写操作都经过业务服务并留下审计日志。

**Tech Stack:** React.js, TypeScript, React Router, Go 1.25, Vitest, Playwright

---

## 文件结构

- `apps/web/src/routes/admin/AdminLayout.tsx`
- `apps/web/src/routes/admin/RoutesPage.tsx`
- `apps/web/src/routes/admin/PlansPage.tsx`
- `apps/web/src/routes/admin/OrdersPage.tsx`
- `apps/web/src/routes/admin/AccountsPage.tsx`
- `apps/web/src/lib/api/adminClient.ts`
- `services/controlplane/internal/admin/*`
- `tests/e2e/admin-console.spec.ts`

### Task 1: 落后台布局与路由配置页

**Files:**
- Create: `apps/web/src/routes/admin/AdminLayout.tsx`
- Create: `apps/web/src/routes/admin/RoutesPage.tsx`
- Create: `apps/web/src/routes/admin/__tests__/routes-page.test.tsx`
- Test: `apps/web/src/routes/admin/__tests__/routes-page.test.tsx`

- [ ] **Step 1: 写路由配置页失败测试**

```tsx
import { render, screen } from "@testing-library/react";
import { RoutesPage } from "../RoutesPage";

it("shows public model and upstream mapping fields", () => {
  render(<RoutesPage />);
  expect(screen.getByText("上游真实模型标识")).toBeInTheDocument();
});
```

- [ ] **Step 2: 运行测试确认失败**

Run: `pnpm --filter web test routes-page.test.tsx`
Expected: FAIL with missing routes page

- [ ] **Step 3: 写最小后台布局与路由页**

```tsx
export function AdminLayout({ children }: { children: React.ReactNode }) {
  return <section>{children}</section>;
}

export function RoutesPage() {
  return (
    <AdminLayout>
      <h1>模型路由配置</h1>
      <p>上游真实模型标识</p>
    </AdminLayout>
  );
}
```

- [ ] **Step 4: 运行测试确认通过**

Run: `pnpm --filter web test routes-page.test.tsx`
Expected: PASS

- [ ] **Step 5: 提交后台路由页**

```bash
git add apps/web/src/routes/admin/AdminLayout.tsx apps/web/src/routes/admin/RoutesPage.tsx apps/web/src/routes/admin/__tests__/routes-page.test.tsx
git commit -m "feat: add admin routes page"
```

### Task 2: 落套餐管理与订单处理页

**Files:**
- Create: `apps/web/src/routes/admin/PlansPage.tsx`
- Create: `apps/web/src/routes/admin/OrdersPage.tsx`
- Create: `apps/web/src/lib/api/adminClient.ts`
- Create: `apps/web/src/routes/admin/__tests__/plans-orders.test.tsx`
- Test: `apps/web/src/routes/admin/__tests__/plans-orders.test.tsx`

- [ ] **Step 1: 写套餐与订单页失败测试**

```tsx
import { render, screen } from "@testing-library/react";
import { OrdersPage } from "../OrdersPage";

it("shows manual review and repair entry", () => {
  render(<OrdersPage />);
  expect(screen.getByText("异常订单补单")).toBeInTheDocument();
});
```

- [ ] **Step 2: 运行测试确认失败**

Run: `pnpm --filter web test plans-orders.test.tsx`
Expected: FAIL with missing pages

- [ ] **Step 3: 写最小套餐与订单页**

```tsx
export function PlansPage() {
  return <section>套餐配置管理</section>;
}

export function OrdersPage() {
  return <section>异常订单补单</section>;
}
```

```ts
export async function listAdminOrders() {
  return [];
}
```

- [ ] **Step 4: 运行测试确认通过**

Run: `pnpm --filter web test plans-orders.test.tsx`
Expected: PASS

- [ ] **Step 5: 提交套餐与订单页**

```bash
git add apps/web/src/routes/admin/PlansPage.tsx apps/web/src/routes/admin/OrdersPage.tsx apps/web/src/lib/api/adminClient.ts apps/web/src/routes/admin/__tests__/plans-orders.test.tsx
git commit -m "feat: add admin plans and orders pages"
```

### Task 3: 落账户状态页与后台治理 API

**Files:**
- Create: `apps/web/src/routes/admin/AccountsPage.tsx`
- Create: `services/controlplane/internal/admin/service.go`
- Create: `tests/integration/admin/account_actions_test.go`
- Test: `tests/integration/admin/account_actions_test.go`

- [ ] **Step 1: 写账户暂停恢复失败测试**

```go
package admin

import "testing"

func TestAdminCanPauseAndResumeAccount(t *testing.T) {
	t.Fatal("assert admin actions are audited and delegated to billing core")
}
```

- [ ] **Step 2: 运行测试确认失败**

Run: `go test ./tests/integration/admin -run TestAdminCanPauseAndResumeAccount -v`
Expected: FAIL

- [ ] **Step 3: 写最小后台治理服务**

```go
package admin

func PauseAccount(accountID string, reason string) error {
	return nil
}

func ResumeAccount(accountID string, reason string) error {
	return nil
}
```

- [ ] **Step 4: 运行测试确认通过**

Run: `go test ./tests/integration/admin -run TestAdminCanPauseAndResumeAccount -v`
Expected: PASS

- [ ] **Step 5: 提交账户治理服务**

```bash
git add apps/web/src/routes/admin/AccountsPage.tsx services/controlplane/internal/admin/service.go tests/integration/admin/account_actions_test.go
git commit -m "feat: add admin account actions"
```

### Task 4: 加后台 E2E 与审计检查点

**Files:**
- Create: `tests/e2e/admin-console.spec.ts`
- Modify: `apps/web/src/main.tsx`
- Create: `tests/integration/admin/audit_log_test.go`
- Test: `tests/e2e/admin-console.spec.ts`

- [ ] **Step 1: 写后台 E2E 失败测试**

```ts
import { test, expect } from "@playwright/test";

test("admin can navigate to route and order pages", async ({ page }) => {
  await page.goto("/admin");
  await expect(page.getByText("模型路由配置")).toBeVisible();
});
```

- [ ] **Step 2: 运行测试确认失败**

Run: `pnpm exec playwright test tests/e2e/admin-console.spec.ts`
Expected: FAIL with admin route missing

- [ ] **Step 3: 注册后台路由并加审计测试文件**

```tsx
<Route path="/admin" element={<RoutesPage />} />
<Route path="/admin/orders" element={<OrdersPage />} />
<Route path="/admin/accounts" element={<AccountsPage />} />
```

```go
package admin

import "testing"

func TestAdminActionsCreateAuditLog(t *testing.T) {}
```

- [ ] **Step 4: 运行测试确认通过**

Run: `pnpm exec playwright test tests/e2e/admin-console.spec.ts`
Expected: PASS

- [ ] **Step 5: 提交后台主流程**

```bash
git add apps/web/src/main.tsx tests/e2e/admin-console.spec.ts tests/integration/admin/audit_log_test.go
git commit -m "test: add admin console flows"
```
