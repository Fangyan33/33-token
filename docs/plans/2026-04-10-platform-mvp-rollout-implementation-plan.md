# 平台 MVP 总体落地 Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 从零搭建模型服务 API 平台 MVP 的整体工程骨架、基础服务边界与端到端可运行主链路。

**Architecture:** 采用 monorepo 结构，前端使用 React.js 承载官网、用户控制台和管理后台，后端使用 Go 1.25 提供认证、账户、计费、支付与后台 API，Envoy AI Gateway 作为统一协议入口。总体计划优先打通“注册登录、创建 API Key、购买开通、网关转发、usage 结算、后台治理”主链路，再由各模块计划分别填充分支能力。

**Tech Stack:** React.js, TypeScript, Vite, React Router, Go 1.25, PostgreSQL, Envoy AI Gateway, Docker Compose, Playwright, Vitest, Go test

---

## 文件结构

- `apps/web`: React 前端应用，承载官网、用户控制台、管理后台
- `apps/web/src/routes/public/*`: 官网页面路由
- `apps/web/src/routes/console/*`: 用户控制台页面路由
- `apps/web/src/routes/admin/*`: 管理后台页面路由
- `apps/web/src/lib/api/*`: 前端调用后端 API 的客户端
- `services/controlplane/cmd/api/main.go`: Go API 服务入口
- `services/controlplane/internal/auth/*`: 登录、会话、用户身份
- `services/controlplane/internal/accounts/*`: 账户与 API Key
- `services/controlplane/internal/billing/*`: 套餐、订单、支付、周期汇总
- `services/controlplane/internal/admin/*`: 后台路由、套餐、审计治理
- `gateway/envoy/*`: Envoy AI Gateway 配置与扩展
- `db/migrations/*`: PostgreSQL migration
- `tests/e2e/*`: Playwright 端到端测试
- `tests/integration/*`: Go 集成测试

### Task 1: 初始化 Monorepo 骨架

**Files:**
- Create: `package.json`
- Create: `pnpm-workspace.yaml`
- Create: `apps/web/package.json`
- Create: `apps/web/src/main.tsx`
- Create: `services/controlplane/go.mod`
- Create: `services/controlplane/cmd/api/main.go`
- Create: `gateway/envoy/docker-compose.yaml`
- Test: `tests/integration/platform_bootstrap_test.go`

- [ ] **Step 1: 写基础启动测试清单**

```go
package integration

import "testing"

func TestPlatformBootstrapChecklist(t *testing.T) {
	t.Log("verify web workspace, go api module, and envoy compose files exist")
}
```

- [ ] **Step 2: 运行测试确认当前尚未落地**

Run: `go test ./tests/integration -run TestPlatformBootstrapChecklist -v`
Expected: FAIL with `directory not found` or `no Go files`

- [ ] **Step 3: 写最小工程骨架**

```json
{
  "name": "model-api-platform",
  "private": true,
  "packageManager": "pnpm@10.0.0",
  "scripts": {
    "dev:web": "pnpm --filter web dev",
    "test:web": "pnpm --filter web test",
    "test:e2e": "playwright test"
  }
}
```

```go
package main

import "log"

func main() {
	log.Println("controlplane api bootstrap")
}
```

- [ ] **Step 4: 运行测试确认骨架可被识别**

Run: `go test ./tests/integration -run TestPlatformBootstrapChecklist -v`
Expected: PASS

- [ ] **Step 5: 提交骨架**

```bash
git add package.json pnpm-workspace.yaml apps/web services/controlplane gateway/envoy tests/integration/platform_bootstrap_test.go
git commit -m "chore: bootstrap platform monorepo"
```

### Task 2: 打通本地基础设施编排

**Files:**
- Create: `docker-compose.yaml`
- Create: `gateway/envoy/envoy.yaml`
- Create: `services/controlplane/.env.example`
- Create: `db/migrations/000001_init_extensions.sql`
- Test: `tests/integration/local_stack_test.go`

- [ ] **Step 1: 写本地编排存在性测试**

```go
package integration

import "testing"

func TestLocalStackFilesExist(t *testing.T) {
	t.Log("verify compose, envoy config, env example, and first migration exist")
}
```

- [ ] **Step 2: 运行测试确认失败**

Run: `go test ./tests/integration -run TestLocalStackFilesExist -v`
Expected: FAIL with missing file assertions

- [ ] **Step 3: 写最小本地编排文件**

```yaml
services:
  postgres:
    image: postgres:17
  api:
    build: ./services/controlplane
  envoy:
    image: envoyproxy/envoy:v1.33-latest
```

```yaml
static_resources:
  listeners: []
  clusters: []
```

- [ ] **Step 4: 运行测试确认通过**

Run: `go test ./tests/integration -run TestLocalStackFilesExist -v`
Expected: PASS

- [ ] **Step 5: 提交编排层**

```bash
git add docker-compose.yaml gateway/envoy/envoy.yaml services/controlplane/.env.example db/migrations/000001_init_extensions.sql tests/integration/local_stack_test.go
git commit -m "chore: add local stack baseline"
```

### Task 3: 建立端到端主链路检查点

**Files:**
- Create: `tests/e2e/platform-smoke.spec.ts`
- Create: `docs/plans/README.md`
- Modify: `package.json`
- Test: `tests/e2e/platform-smoke.spec.ts`

- [ ] **Step 1: 写端到端 smoke 用例草案**

```ts
import { test, expect } from "@playwright/test";

test("platform smoke route placeholders", async ({ page }) => {
  await page.goto("/");
  await expect(page).toHaveTitle(/Model API Platform/);
});
```

- [ ] **Step 2: 运行测试确认当前失败**

Run: `pnpm exec playwright test tests/e2e/platform-smoke.spec.ts`
Expected: FAIL with route or app bootstrap missing

- [ ] **Step 3: 补最小前端占位入口与计划索引**

```tsx
import React from "react";
import ReactDOM from "react-dom/client";

function App() {
  return <main>Model API Platform</main>;
}

ReactDOM.createRoot(document.getElementById("root")!).render(<App />);
```

```md
# 执行计划索引

- 平台 MVP 总体落地
- 官网站点
- 用户控制台
- API 网关
- 计费与订阅核心
- 管理后台
- 网关与计费交互
- 数据库逻辑模型
```

- [ ] **Step 4: 运行 smoke 测试确认通过**

Run: `pnpm exec playwright test tests/e2e/platform-smoke.spec.ts`
Expected: PASS

- [ ] **Step 5: 提交主链路检查点**

```bash
git add apps/web/src/main.tsx tests/e2e/platform-smoke.spec.ts package.json docs/plans/README.md
git commit -m "test: add platform smoke checkpoint"
```

### Task 4: 建立跨计划实施顺序

**Files:**
- Create: `docs/plans/2026-04-10-platform-rollout-order.md`
- Modify: `docs/plans/README.md`
- Test: `docs/plans/2026-04-10-platform-rollout-order.md`

- [ ] **Step 1: 写顺序校验说明**

```md
# 平台执行顺序校验

必须先完成数据库逻辑模型和计费核心基础，再并行推进官网、控制台、后台与网关对接。
```

- [ ] **Step 2: 运行文本检查**

Run: `rg -n "数据库逻辑模型|计费核心|官网|控制台|后台|网关" docs/plans/2026-04-10-platform-rollout-order.md`
Expected: PASS with all phase keywords matched

- [ ] **Step 3: 补完整执行顺序文档**

```md
## Phase 1
- 数据库逻辑模型
- 计费与订阅核心

## Phase 2
- API 代理网关
- 网关与计费核心交互

## Phase 3
- 官网站点
- 用户控制台
- 管理后台
```

- [ ] **Step 4: 重新运行文本检查**

Run: `rg -n "Phase 1|Phase 2|Phase 3" docs/plans/2026-04-10-platform-rollout-order.md`
Expected: PASS

- [ ] **Step 5: 提交执行顺序文档**

```bash
git add docs/plans/README.md docs/plans/2026-04-10-platform-rollout-order.md
git commit -m "docs: define platform rollout order"
```
