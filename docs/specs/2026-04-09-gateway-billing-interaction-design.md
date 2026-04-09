# 网关与计费核心交互时序设计

> 本文补充 [`2026-04-09-api-gateway-module-design.md`](./2026-04-09-api-gateway-module-design.md) 与 [`2026-04-09-billing-and-subscription-core-module-design.md`](./2026-04-09-billing-and-subscription-core-module-design.md) 之间的接口边界与交互时序。
> 本文不替代模块设计文档，而是细化两者在额度校验、事件记录与状态切换上的协作方式。

## 1. 设计目标

本文用于明确以下问题：

- 网关在转发前读取什么数据来判断是否允许调用
- 请求完成后 `usage_event`、`billing_event` 与周期汇总如何更新
- 额度耗尽后如何从可调用状态切换为暂停状态
- 出现写入失败、重复提交或上游异常时如何保持系统可恢复

## 2. 设计原则

本设计采用以下固定原则：

- 账户可调用性判断只读取账户当前周期汇总
- 原始事件表只负责留痕、审计、补账和重算，不参与实时放行判断
- API 网关负责请求执行与事件上报，不负责维护账务真相
- 计费与订阅核心负责汇总更新、状态切换与账务一致性
- 第一版不引入 Redis 作为额度判定前提
- 第一版不做请求级额度预占
- 第一版接受“最后一个请求可能把额度打穿一点点”的边界

## 3. 核心数据对象

### 3.1 账户周期汇总

账户周期汇总是网关实时准入时读取的主数据对象，建议至少包含：

- `account_id`
- `billing_period`
- `quota_total`
- `quota_used`
- `quota_remaining`
- `status`
- `updated_at`

其中：

- `status` 用于表达当前账户在本周期内是否可调用，例如 `active`、`paused`、`expired`
- `quota_remaining` 用于表达是否还有可消费额度

### 3.2 usage_event

`usage_event` 表示一次请求的原始事实，建议至少包含：

- `event_id`
- `account_id`
- `api_key_id`
- `protocol`
- `model_name`
- `upstream_provider`
- `input_tokens`
- `output_tokens`
- `total_tokens`
- `request_started_at`
- `request_finished_at`
- `result_status`
- `error_type`
- `latency_ms`
- `billing_period`

### 3.3 billing_event

`billing_event` 表示一次请求在账务层的处理结果，建议至少包含：

- `event_id`
- `idempotency_key`
- `account_id`
- `related_usage_event_id`
- `billing_period`
- `settlement_type`
- `quota_delta`
- `settlement_result`
- `created_at`

`usage_event` 与 `billing_event` 的区别是：

- `usage_event` 记录请求事实
- `billing_event` 记录账务动作

## 4. 实时校验链路

### 4.1 转发前校验输入

网关在转发前，至少需要以下输入：

- API Key 校验结果
- API Key 所属账户信息
- 账户当前周期汇总

### 4.2 转发前校验规则

网关按以下顺序判断是否允许调用：

1. API Key 是否存在且有效
2. 账户是否被停用
3. 是否存在有效套餐
4. 当前账单周期状态是否仍为可调用
5. `quota_remaining` 是否大于 0

若任一条件不满足，则网关必须在转发前直接阻断请求，不调用上游模型服务。

### 4.3 实时读模型要求

转发前校验只读取账户当前周期汇总，不扫描原始事件表。

原因是：

- 原始事件表适合追溯，不适合高频准入判断
- 准入判断需要稳定、简单、可预测
- 当前周期汇总更适合作为实时判定读模型

## 5. 请求完成后的记账链路

### 5.1 基本流程

请求完成后，网关将完整请求结果提交给计费与订阅核心，计费核心按以下顺序处理：

1. 写入 `usage_event`
2. 写入 `billing_event`
3. 更新账户当前周期汇总
4. 若累计后额度耗尽，则将账户状态切换为 `paused`

### 5.2 状态切换规则

若某次请求处理完成后：

- `quota_remaining` 仍大于 0，则账户继续保持可调用
- `quota_remaining` 小于等于 0，则账户切换为 `paused`

该状态更新完成后，后续请求在下一次转发前校验中会被直接阻断。

### 5.3 MVP 边界

由于第一版不做额度预占，因此存在以下边界：

- 只要请求转发前读到的 `quota_remaining > 0`，该请求就允许进入
- 某次请求可能在完成后将额度打穿一点点
- 从该请求之后的下一次请求开始，系统必须稳定阻断

这是当前 MVP 明确接受的复杂度取舍。

## 6. 内部接口边界

为避免网关和计费核心重复实现账务逻辑，建议两者之间至少有两个内部能力。

### 6.1 准入校验能力

能力名称可定义为：

- `CheckAccountQuota`

输入建议包括：

- `account_id`
- `api_key_id`

输出建议包括：

- 是否允许调用
- 不允许调用的原因
- 当前账单周期
- 当前剩余额度摘要
- 当前账户状态摘要

### 6.2 请求结算能力

能力名称可定义为：

- `RecordUsageAndSettle`

输入建议为一次完整请求结果，至少包括：

- 账户信息
- API Key 信息
- 协议信息
- 模型信息
- 请求结果状态
- token 使用量
- 请求耗时
- 错误信息
- 周期归属

该能力内部负责：

- 写入 `usage_event`
- 写入 `billing_event`
- 更新账户当前周期汇总
- 必要时切换账户到 `paused`

## 7. 异常与补偿设计

### 7.1 上游失败

若请求到达上游但失败，仍应记录 `usage_event`。

`billing_event` 是否产生扣减，取决于平台最终确定的计费规则，但规则必须固定，不能由网关临时决定。

### 7.2 原始事件写入成功、汇总更新失败

这是需要优先控制的异常。

建议：

- 保留“待汇总”或等价可重试标记
- 允许后台或异步任务重放未完成结算的记录
- 不允许因为单次汇总失败而永久丢失账务事实

### 7.3 计费核心短时不可写

若网关拿到上游结果后，计费核心暂时无法写入，系统至少要保留可重试落点。

第一版如果不引入消息队列，也应保留数据库内可重试记录，避免事件直接丢失。

### 7.4 重复提交

`billing_event` 应具备幂等键，避免重复累计。

同一请求重复上报时，应保证：

- 原始事实不丢失
- 账务累计不会重复执行

## 8. 测试与验收重点

这份交互时序设计至少应验证：

- 网关转发前只依赖账户当前周期汇总进行准入判断
- 请求完成后 `usage_event`、`billing_event` 与周期汇总按既定顺序处理
- 某次请求耗尽额度后，后续请求会被稳定阻断
- 上游失败时原始事件仍可落库
- 汇总失败时存在明确补偿路径
- 重复提交不会导致重复计费
