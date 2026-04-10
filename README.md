# 模型服务 API 平台

本仓库用于承载模型服务 API 平台 MVP 的整体工程骨架。

当前阶段目标：

- 使用 React 搭建官网、用户控制台、管理后台的统一前端应用
- 使用 Go 1.25 搭建 controlplane API 服务骨架
- 使用 Envoy AI Gateway 预留 OpenAI / Anthropic 兼容入口配置
- 使用 PostgreSQL 17 作为核心数据存储

## 仓库结构

- `apps/web`：React 前端应用，承载官网、控制台、后台
- `services/controlplane`：Go 控制面服务
- `gateway/envoy`：Envoy 网关配置
- `db/migrations`：数据库迁移脚本
- `tests/integration`：集成级骨架检查
- `tests/e2e`：端到端 smoke 测试
- `docs/specs`：设计文档
- `docs/plans`：实施计划文档

## 常用命令

- `pnpm install`
- `pnpm dev:web`
- `pnpm test:web`
- `pnpm test:e2e`
- `go test ./tests/integration/...`
- `cd services/controlplane && go test ./...`
- `docker compose up --build`

## 当前状态

当前只落地工程骨架和占位入口，未实现注册登录、支付、API Key、网关转发与 usage 结算等业务能力。
