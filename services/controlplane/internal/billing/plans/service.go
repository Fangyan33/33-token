package plans

import "github.com/33-token/model-api-platform/services/controlplane/internal/store"

type Service struct{}

func NewService() Service {
	return Service{}
}

func (Service) ListSellable(plans []store.Plan) []store.Plan {
	sellable := make([]store.Plan, 0, len(plans))
	for _, plan := range plans {
		if plan.Status == "active" {
			sellable = append(sellable, plan)
		}
	}
	return sellable
}
