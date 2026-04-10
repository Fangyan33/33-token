# 官网站点模块 Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 实现官网 MVP，覆盖首页、价格页、兼容性说明、模型列表、文档页以及注册登录和购买入口。

**Architecture:** 官网作为 React.js 公共站点实现，优先使用静态内容加少量公开 API 获取套餐和模型列表。官网不持有复杂业务状态，只负责转化入口和开发者教育。

**Tech Stack:** React.js, TypeScript, React Router, Vitest, Playwright

---

## 文件结构

- `apps/web/src/routes/public/HomePage.tsx`
- `apps/web/src/routes/public/PricingPage.tsx`
- `apps/web/src/routes/public/CompatibilityPage.tsx`
- `apps/web/src/routes/public/ModelsPage.tsx`
- `apps/web/src/routes/public/DocsPage.tsx`
- `apps/web/src/routes/public/AuthEntryPage.tsx`
- `apps/web/src/components/public/*`
- `apps/web/src/lib/api/publicSiteClient.ts`
- `apps/web/src/routes/public/__tests__/*.test.tsx`
- `tests/e2e/official-website.spec.ts`

### Task 1: 落首页与公共布局

**Files:**
- Create: `apps/web/src/routes/public/HomePage.tsx`
- Create: `apps/web/src/components/public/PublicLayout.tsx`
- Create: `apps/web/src/routes/public/__tests__/home-page.test.tsx`
- Test: `apps/web/src/routes/public/__tests__/home-page.test.tsx`

- [ ] **Step 1: 写首页失败测试**

```tsx
import { render, screen } from "@testing-library/react";
import { HomePage } from "../HomePage";

it("renders product value proposition", () => {
  render(<HomePage />);
  expect(screen.getByText("低价代码模型 API 平台")).toBeInTheDocument();
});
```

- [ ] **Step 2: 运行测试确认失败**

Run: `pnpm --filter web test home-page.test.tsx`
Expected: FAIL with `HomePage` module missing

- [ ] **Step 3: 写最小首页与布局**

```tsx
export function PublicLayout({ children }: { children: React.ReactNode }) {
  return <div>{children}</div>;
}

export function HomePage() {
  return (
    <PublicLayout>
      <h1>低价代码模型 API 平台</h1>
      <p>兼容 OpenAI 与 Anthropic 调用方式。</p>
    </PublicLayout>
  );
}
```

- [ ] **Step 4: 运行测试确认通过**

Run: `pnpm --filter web test home-page.test.tsx`
Expected: PASS

- [ ] **Step 5: 提交首页基础**

```bash
git add apps/web/src/routes/public/HomePage.tsx apps/web/src/components/public/PublicLayout.tsx apps/web/src/routes/public/__tests__/home-page.test.tsx
git commit -m "feat: add website home page"
```

### Task 2: 落价格页与模型列表页

**Files:**
- Create: `apps/web/src/routes/public/PricingPage.tsx`
- Create: `apps/web/src/routes/public/ModelsPage.tsx`
- Create: `apps/web/src/lib/api/publicSiteClient.ts`
- Create: `apps/web/src/routes/public/__tests__/pricing-models.test.tsx`
- Test: `apps/web/src/routes/public/__tests__/pricing-models.test.tsx`

- [ ] **Step 1: 写价格与模型页失败测试**

```tsx
import { render, screen } from "@testing-library/react";
import { PricingPage } from "../PricingPage";

it("shows quota suspension messaging", () => {
  render(<PricingPage plans={[]} />);
  expect(screen.getByText("超额后暂停调用")).toBeInTheDocument();
});
```

- [ ] **Step 2: 运行测试确认失败**

Run: `pnpm --filter web test pricing-models.test.tsx`
Expected: FAIL with missing component

- [ ] **Step 3: 写最小页面与公开 API 客户端**

```tsx
export function PricingPage() {
  return (
    <section>
      <h2>套餐价格</h2>
      <p>超额后暂停调用</p>
    </section>
  );
}

export function ModelsPage() {
  return <section>公开模型列表</section>;
}
```

```ts
export async function listPublicPlans() {
  return [];
}
```

- [ ] **Step 4: 运行测试确认通过**

Run: `pnpm --filter web test pricing-models.test.tsx`
Expected: PASS

- [ ] **Step 5: 提交公开信息页**

```bash
git add apps/web/src/routes/public/PricingPage.tsx apps/web/src/routes/public/ModelsPage.tsx apps/web/src/lib/api/publicSiteClient.ts apps/web/src/routes/public/__tests__/pricing-models.test.tsx
git commit -m "feat: add website pricing and models pages"
```

### Task 3: 落兼容性说明、文档页与入口页

**Files:**
- Create: `apps/web/src/routes/public/CompatibilityPage.tsx`
- Create: `apps/web/src/routes/public/DocsPage.tsx`
- Create: `apps/web/src/routes/public/AuthEntryPage.tsx`
- Create: `apps/web/src/routes/public/__tests__/docs-entry.test.tsx`
- Test: `apps/web/src/routes/public/__tests__/docs-entry.test.tsx`

- [ ] **Step 1: 写接入说明失败测试**

```tsx
import { render, screen } from "@testing-library/react";
import { DocsPage } from "../DocsPage";

it("shows OpenAI compatible quickstart", () => {
  render(<DocsPage />);
  expect(screen.getByText("OpenAI 兼容接入示例")).toBeInTheDocument();
});
```

- [ ] **Step 2: 运行测试确认失败**

Run: `pnpm --filter web test docs-entry.test.tsx`
Expected: FAIL with missing docs page

- [ ] **Step 3: 写最小接入与入口页面**

```tsx
export function CompatibilityPage() {
  return <section>支持 OpenAI 兼容 API 与 Anthropic 兼容 API</section>;
}

export function DocsPage() {
  return <section>OpenAI 兼容接入示例</section>;
}

export function AuthEntryPage() {
  return <section>注册 / 登录 / 购买入口</section>;
}
```

- [ ] **Step 4: 运行测试确认通过**

Run: `pnpm --filter web test docs-entry.test.tsx`
Expected: PASS

- [ ] **Step 5: 提交接入与入口页**

```bash
git add apps/web/src/routes/public/CompatibilityPage.tsx apps/web/src/routes/public/DocsPage.tsx apps/web/src/routes/public/AuthEntryPage.tsx apps/web/src/routes/public/__tests__/docs-entry.test.tsx
git commit -m "feat: add website docs and auth entry pages"
```

### Task 4: 加官网端到端主流程

**Files:**
- Create: `tests/e2e/official-website.spec.ts`
- Modify: `apps/web/src/main.tsx`
- Test: `tests/e2e/official-website.spec.ts`

- [ ] **Step 1: 写官网 E2E 用例**

```ts
import { test, expect } from "@playwright/test";

test("website navigation to pricing and docs", async ({ page }) => {
  await page.goto("/");
  await page.getByText("套餐价格").click();
  await expect(page.getByText("超额后暂停调用")).toBeVisible();
});
```

- [ ] **Step 2: 运行测试确认失败**

Run: `pnpm exec playwright test tests/e2e/official-website.spec.ts`
Expected: FAIL with route navigation missing

- [ ] **Step 3: 补公共路由注册**

```tsx
function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<HomePage />} />
        <Route path="/pricing" element={<PricingPage />} />
        <Route path="/docs" element={<DocsPage />} />
      </Routes>
    </BrowserRouter>
  );
}
```

- [ ] **Step 4: 运行测试确认通过**

Run: `pnpm exec playwright test tests/e2e/official-website.spec.ts`
Expected: PASS

- [ ] **Step 5: 提交官网主流程**

```bash
git add apps/web/src/main.tsx tests/e2e/official-website.spec.ts
git commit -m "test: add website e2e flow"
```
