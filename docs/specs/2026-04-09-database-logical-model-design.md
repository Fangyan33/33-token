# 数据库逻辑模型设计

> 本文统一整理模型服务 API 平台 MVP 阶段的核心数据库逻辑模型。
> 本文补充 [`2026-03-27-model-api-platform-design.md`](./2026-03-27-model-api-platform-design.md) 及各模块细化文档中的数据对象定义，重点描述逻辑实体、职责边界与关联关系，不展开具体 SQL DDL。

## 1. 设计目标

本文用于统一以下内容：

- 各模块共享的数据对象边界
- 核心表之间的逻辑关系
- 实时判定读模型与原始事实数据的分层方式
- MVP 阶段应保留和应避免的数据库复杂度

本文不直接定义：

- 具体数据库引擎参数
- 具体索引语句
- 具体分库分表策略
- ORM 映射细节

## 2. 总体建模原则

数据库逻辑模型遵循以下原则：

- 账户是计费与权限归属的根对象
- 请求事实、账务动作、实时汇总三类数据必须分开
- 套餐定义、订单、支付事件必须分开
- 运营配置、敏感凭证引用、后台审计必须分开
- MVP 阶段不引入 Redis 作为额度判定前提
- MVP 阶段不做复杂额度预占，接受“最后一个请求可能把额度打穿一点点”

## 3. 逻辑域划分

建议将数据库逻辑模型拆分为 5 个逻辑域：

1. 账户与认证域
2. 套餐与订单域
3. 账户权益与周期汇总域
4. 网关计量与事件域
5. 运营配置与后台审计域

## 4. 账户与认证域

### 4.1 account

`account` 是业务主体和计费主体。

建议至少包含：

- `id`
- `status`
- `display_name`
- `default_contact_email`
- `created_at`
- `updated_at`

说明：

- `status` 只表达账户级长期状态，例如 `active`、`disabled`
- 套餐过期、额度耗尽等周期性状态不放在该对象中

### 4.2 user_identity

`user_identity` 是登录主体，用于表示谁能登录官网、控制台或后台。

建议至少包含：

- `id`
- `account_id`
- `login_email`
- `auth_provider`
- `status`
- `last_login_at`
- `created_at`

说明：

- MVP 即使是一账户一用户，也建议保留该对象
- 后续新增 OAuth、多身份或后台独立登录体系时无需重构账户模型

### 4.3 api_key

`api_key` 是调用凭证主体，必须直接归属到账户。

建议至少包含：

- `id`
- `account_id`
- `key_prefix`
- `key_hash`
- `status`
- `created_at`
- `disabled_at`
- `last_used_at`

说明：

- 只存 `key_hash`，不存明文
- 控制台创建时只返回一次明文
- `status` 至少支持 `active`、`disabled`

### 4.4 关系

- `user_identity.account_id -> account.id`
- `api_key.account_id -> account.id`

## 5. 套餐与订单域

### 5.1 plan

`plan` 表示当前对外出售的套餐定义。

建议至少包含：

- `id`
- `code`
- `name`
- `status`
- `billing_period_type`
- `quota_total`
- `rate_limit_policy`
- `display_priority`
- `created_at`
- `updated_at`

### 5.2 plan_price_snapshot

`plan_price_snapshot` 表示用户下单时冻结的成交快照。

建议至少包含：

- `id`
- `plan_id`
- `plan_code`
- `plan_name`
- `price_amount`
- `currency`
- `billing_period_type`
- `quota_total`
- `rate_limit_policy_snapshot`
- `created_at`

说明：

- 订单应关联套餐快照，而不是直接依赖当前套餐定义
- 套餐改价或改额度后，历史订单仍保留成交依据

### 5.3 order

`order` 表示一次购买过程。

建议至少包含：

- `id`
- `account_id`
- `plan_price_snapshot_id`
- `order_type`
- `status`
- `payment_provider`
- `payment_provider_order_id`
- `amount`
- `currency`
- `created_at`
- `paid_at`
- `completed_at`

`status` 至少建议支持：

- `pending_payment`
- `paid_pending_activation`
- `completed`
- `failed`
- `cancelled`
- `manual_review_required`

### 5.4 payment_event

`payment_event` 表示支付平台返回的外部事实。

建议至少包含：

- `id`
- `order_id`
- `payment_provider`
- `provider_event_id`
- `event_type`
- `event_status`
- `raw_reference`
- `event_occurred_at`
- `received_at`

说明：

- 一个订单可对应多条支付事件
- 支付事件不直接等于订单已完成

### 5.5 关系

- `plan_price_snapshot.plan_id -> plan.id`
- `order.account_id -> account.id`
- `order.plan_price_snapshot_id -> plan_price_snapshot.id`
- `payment_event.order_id -> order.id`

## 6. 账户权益与周期汇总域

### 6.1 account_subscription_state

`account_subscription_state` 表示账户当前长期权益状态。

建议至少包含：

- `id`
- `account_id`
- `current_order_id`
- `current_plan_price_snapshot_id`
- `status`
- `billing_period_start`
- `billing_period_end`
- `activated_at`
- `expired_at`
- `paused_at`
- `updated_at`

`status` 至少建议支持：

- `inactive`
- `active`
- `paused`
- `expired`
- `disabled`

说明：

- `disabled` 偏人工停用或封禁
- `paused` 偏额度耗尽或运营动作导致的当前不可调用
- `expired` 偏周期自然结束

### 6.2 account_cycle_summary

`account_cycle_summary` 表示账户在当前周期内的实时汇总读模型。

建议至少包含：

- `id`
- `account_id`
- `billing_period_start`
- `billing_period_end`
- `quota_total`
- `quota_used`
- `quota_remaining`
- `input_tokens_total`
- `output_tokens_total`
- `total_tokens_total`
- `status`
- `updated_at`

`status` 至少建议支持：

- `active`
- `paused`
- `closed`

说明：

- 该对象用于网关转发前的实时准入判断
- 该对象不承载订单和支付来源等长期状态信息

### 6.3 关系

- `account_subscription_state.account_id -> account.id`
- `account_subscription_state.current_order_id -> order.id`
- `account_subscription_state.current_plan_price_snapshot_id -> plan_price_snapshot.id`
- `account_cycle_summary.account_id -> account.id`

### 6.4 职责边界

- `account_subscription_state` 回答账户当前处于什么权益状态
- `account_cycle_summary` 回答当前周期还剩多少额度、是否可继续调用

## 7. 网关计量与事件域

### 7.1 usage_event

`usage_event` 表示单次请求的原始事实。

建议至少包含：

- `id`
- `request_id`
- `account_id`
- `api_key_id`
- `protocol`
- `model_name`
- `upstream_provider`
- `upstream_model_id`
- `request_started_at`
- `request_finished_at`
- `latency_ms`
- `result_status`
- `error_type`
- `input_tokens`
- `output_tokens`
- `total_tokens`
- `billing_period_start`
- `billing_period_end`
- `created_at`

### 7.2 billing_event

`billing_event` 表示单次请求在账务层的处理结果。

建议至少包含：

- `id`
- `idempotency_key`
- `account_id`
- `usage_event_id`
- `billing_period_start`
- `billing_period_end`
- `settlement_type`
- `quota_delta`
- `before_quota_remaining`
- `after_quota_remaining`
- `result_status`
- `failure_reason`
- `created_at`
- `settled_at`

说明：

- `usage_event` 记录请求事实
- `billing_event` 记录账务动作
- 两者不能合并

### 7.3 建议状态

`usage_event.result_status` 至少可区分：

- `success`
- `upstream_error`
- `blocked`
- `internal_error`

`billing_event.result_status` 至少可区分：

- `applied`
- `skipped`
- `failed`

### 7.4 关系

- `usage_event.account_id -> account.id`
- `usage_event.api_key_id -> api_key.id`
- `billing_event.account_id -> account.id`
- `billing_event.usage_event_id -> usage_event.id`

## 8. 运营配置与后台审计域

### 8.1 model_route

`model_route` 表示平台对外模型名与上游模型之间的路由映射。

建议至少包含：

- `id`
- `protocol`
- `public_model_name`
- `upstream_provider`
- `upstream_model_id`
- `upstream_credential_ref_id`
- `status`
- `request_adapter_type`
- `response_adapter_type`
- `priority`
- `created_at`
- `updated_at`

### 8.2 upstream_credential_ref

`upstream_credential_ref` 表示上游凭证引用，而不是凭证明文本体。

建议至少包含：

- `id`
- `provider`
- `credential_key`
- `display_name`
- `status`
- `created_at`
- `updated_at`

说明：

- 数据库逻辑模型中只保存引用关系
- 真正密钥应交给环境变量或受控密钥管理系统

### 8.3 admin_audit_log

`admin_audit_log` 表示后台高风险写操作的审计日志。

建议至少包含：

- `id`
- `operator_id`
- `operator_type`
- `action_type`
- `target_type`
- `target_id`
- `before_snapshot`
- `after_snapshot`
- `reason`
- `created_at`

### 8.4 关系

- `model_route.upstream_credential_ref_id -> upstream_credential_ref.id`
- `admin_audit_log` 通过 `target_type + target_id` 指向被操作对象

## 9. 关键读写路径

### 9.1 网关实时准入

网关转发前主要读取：

- `api_key`
- `account`
- `account_subscription_state`
- `account_cycle_summary`

判断逻辑聚焦于：

- API Key 是否有效
- 账户是否被停用
- 当前权益是否有效
- 当前周期状态是否可调用
- `quota_remaining > 0`

### 9.2 请求结算

请求完成后主要写入：

- `usage_event`
- `billing_event`
- `account_cycle_summary`

必要时同步更新：

- `account_subscription_state`

### 9.3 购买开通

购买与开通流程主要涉及：

- 创建 `order`
- 接收 `payment_event`
- 更新 `account_subscription_state`
- 初始化或刷新 `account_cycle_summary`

## 10. MVP 阶段明确不做的复杂度

MVP 阶段不建议提前引入以下复杂度：

- 多租户组织模型
- 多成员协作权限
- API Key 独立额度
- 多权益包叠加
- 跨周期额度结转
- 复杂财务分录系统
- Kafka 等事件总线作为前提
- 分库分表
- 请求正文与响应正文全文存储

## 11. 与其他设计文档的关系

本文重点补充以下文档中的数据模型层：

- [`2026-04-09-api-gateway-module-design.md`](./2026-04-09-api-gateway-module-design.md)
- [`2026-04-09-billing-and-subscription-core-module-design.md`](./2026-04-09-billing-and-subscription-core-module-design.md)
- [`2026-04-09-user-console-module-design.md`](./2026-04-09-user-console-module-design.md)
- [`2026-04-09-admin-console-module-design.md`](./2026-04-09-admin-console-module-design.md)
- [`2026-04-09-gateway-billing-interaction-design.md`](./2026-04-09-gateway-billing-interaction-design.md)

本文可作为后续数据库表设计、服务接口设计和实现计划拆解的统一依据。
