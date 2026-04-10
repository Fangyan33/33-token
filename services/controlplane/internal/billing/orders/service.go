package orders

// CreateOrderInput 表示创建订单所需的最小输入模型。
type CreateOrderInput struct {
	AccountID           string
	PlanPriceSnapshotID string
	Amount              int64
	Currency            string
}

// OrdersService 是订单领域的最小服务占位。
// 当前任务只需要文件骨架，不接真实业务逻辑。
type OrdersService struct{}
