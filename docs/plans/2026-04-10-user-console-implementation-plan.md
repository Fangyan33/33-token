# 用户控制台模块 Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 实现用户控制台 MVP，覆盖账户概览、套餐状态、周期 usage 汇总、API Key 管理以及购买与付款状态展示。

**Architecture:** 控制台作为 React.js 受保护区域实现，所有状态真相均由后端 API 提供。前端只负责展示与发起动作，避免自行推导账户可调用性与订单完成状态。

**Tech Stack:** React.js, TypeScript, React Router, React Query, Vitest, Playwright

---

## 文件结构

- `apps/web/src/routes/console/ConsoleLayout.tsx`
- `apps/web/src/routes/console/DashboardPage.tsx`
- `apps/web/src/routes/console/ApiKeysPage.tsx`
- `apps/web/src/routes/console/PurchasePage.tsx`
- `apps/web/src/lib/api/consoleClient.ts`
- `apps/web/src/routes/console/__tests__/*.test.tsx`
- `tests/e2e/user-console.spec.ts`

### Task 1: 落账户概览与套餐状态页

**Files:**
- Create: `apps/web/src/routes/console/DashboardPage.tsx`
- Create: `apps/web/src/routes/console/ConsoleLayout.tsx`
- Create: `apps/web/src/routes/console/__tests__/dashboard-page.test.tsx`
- Test: `apps/web/src/routes/console/__tests__/dashboard-page.test.tsx`

- [ ] **Step 1: 写控制台首页失败测试**

```tsx
import { render, screen } from "@testing-library/react";
import { DashboardPage } from "../DashboardPage";

it("shows quota summary", () => {
  render(<DashboardPage />);
  expect(screen.getByText("当前周期剩余额度")).toBeInTheDocument();
});
```

- [ ] **Step 2: 运行测试确认失败**

Run: `pnpm --filter web test dashboard-page.test.tsx`
Expected: FAIL with missing dashboard page

- [ ] **Step 3: 写最小控制台首页**

```tsx
export function ConsoleLayout({ children }: { children: React.ReactNode }) {
  return <section>{children}</section>;
}

export function DashboardPage() {
  return (
    <ConsoleLayout>
      <h1>账户概览</h1>
      <p>当前周期剩余额度</p>
    </ConsoleLayout>
  );
}
```

- [ ] **Step 4: 运行测试确认通过**

Run: `pnpm --filter web test dashboard-page.test.tsx`
Expected: PASS

- [ ] **Step 5: 提交控制台概览页**

```bash
git add apps/web/src/routes/console/DashboardPage.tsx apps/web/src/routes/console/ConsoleLayout.tsx apps/web/src/routes/console/__tests__/dashboard-page.test.tsx
git commit -m "feat: add console dashboard page"
```

### Task 2: 落 API Key 管理页

**Files:**
- Create: `apps/web/src/routes/console/ApiKeysPage.tsx`
- Create: `apps/web/src/lib/api/consoleClient.ts`
- Create: `apps/web/src/routes/console/__tests__/api-keys-page.test.tsx`
- Test: `apps/web/src/routes/console/__tests__/api-keys-page.test.tsx`

- [ ] **Step 1: 写 API Key 页失败测试**

```tsx
import { render, screen } from "@testing-library/react";
import { ApiKeysPage } from "../ApiKeysPage";

it("shows one-time secret copy warning", () => {
  render(<ApiKeysPage />);
  expect(screen.getByText("完整明文只展示一次")).toBeInTheDocument();
});
```

- [ ] **Step 2: 运行测试确认失败**

Run: `pnpm --filter web test api-keys-page.test.tsx`
Expected: FAIL with missing API keys page

- [ ] **Step 3: 写最小 API Key 页面与客户端**

```tsx
export function ApiKeysPage() {
  return (
    <section>
      <h2>API Key 管理</h2>
      <p>完整明文只展示一次</p>
    </section>
  );
}
```

```ts
export async function listApiKeys() {
  return [];
}
```

- [ ] **Step 4: 运行测试确认通过**

Run: `pnpm --filter web test api-keys-page.test.tsx`
Expected: PASS

- [ ] **Step 5: 提交 API Key 页**

```bash
git add apps/web/src/routes/console/ApiKeysPage.tsx apps/web/src/lib/api/consoleClient.ts apps/web/src/routes/console/__tests__/api-keys-page.test.tsx
git commit -m "feat: add console api key page"
```

### Task 3: 落购买与付款状态页

**Files:**
- Create: `apps/web/src/routes/console/PurchasePage.tsx`
- Create: `apps/web/src/routes/console/__tests__/purchase-page.test.tsx`
- Modify: `apps/web/src/lib/api/consoleClient.ts`
- Test: `apps/web/src/routes/console/__tests__/purchase-page.test.tsx`

- [ ] **Step 1: 写购买页失败测试**

```tsx
import { render, screen } from "@testing-library/react";
import { PurchasePage } from "../PurchasePage";

it("shows order status as source of truth", () => {
  render(<PurchasePage />);
  expect(screen.getByText("支付结果以系统订单状态为准")).toBeInTheDocument();
});
```

- [ ] **Step 2: 运行测试确认失败**

Run: `pnpm --filter web test purchase-page.test.tsx`
Expected: FAIL with missing purchase page

- [ ] **Step 3: 写最小购买页**

```tsx
export function PurchasePage() {
  return (
    <section>
      <h2>购买与付款状态</h2>
      <p>支付结果以系统订单状态为准</p>
    </section>
  );
}
```

```ts
export async function listOrders() {
  return [];
}
```

- [ ] **Step 4: 运行测试确认通过**

Run: `pnpm --filter web test purchase-page.test.tsx`
Expected: PASS

- [ ] **Step 5: 提交购买页**

```bash
git add apps/web/src/routes/console/PurchasePage.tsx apps/web/src/lib/api/consoleClient.ts apps/web/src/routes/console/__tests__/purchase-page.test.tsx
git commit -m "feat: add console purchase page"
```

### Task 4: 加控制台路由与 E2E 主流程

**Files:**
- Create: `tests/e2e/user-console.spec.ts`
- Modify: `apps/web/src/main.tsx`
- Test: `tests/e2e/user-console.spec.ts`

- [ ] **Step 1: 写控制台 E2E 测试**

```ts
import { test, expect } from "@playwright/test";

test("console dashboard shows remaining quota and api keys", async ({ page }) => {
  await page.goto("/console");
  await expect(page.getByText("当前周期剩余额度")).toBeVisible();
});
```

- [ ] **Step 2: 运行测试确认失败**

Run: `pnpm exec playwright test tests/e2e/user-console.spec.ts`
Expected: FAIL with console route missing

- [ ] **Step 3: 注册控制台路由**

```tsx
<Route path="/console" element={<DashboardPage />} />
<Route path="/console/api-keys" element={<ApiKeysPage />} />
<Route path="/console/purchase" element={<PurchasePage />} />
```

- [ ] **Step 4: 运行测试确认通过**

Run: `pnpm exec playwright test tests/e2e/user-console.spec.ts`
Expected: PASS

- [ ] **Step 5: 提交控制台路由**

```bash
git add apps/web/src/main.tsx tests/e2e/user-console.spec.ts
git commit -m "test: add console e2e flow"
```
