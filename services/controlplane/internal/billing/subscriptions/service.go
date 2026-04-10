package subscriptions

import "errors"

var ErrOrderIDRequired = errors.New("order id is required")

type ActivationResult struct {
	OrderID   string
	Activated bool
	Status    string
}

type Service struct {
	activated map[string]struct{}
}

func NewService() *Service {
	return &Service{
		activated: map[string]struct{}{},
	}
}

var defaultService = NewService()

func ActivateFromPaidOrder(orderID string) error {
	_, err := defaultService.Activate(orderID)
	return err
}

func (s *Service) Activate(orderID string) (ActivationResult, error) {
	if orderID == "" {
		return ActivationResult{}, ErrOrderIDRequired
	}

	if _, ok := s.activated[orderID]; ok {
		return ActivationResult{
			OrderID:   orderID,
			Activated: false,
			Status:    "active",
		}, nil
	}

	s.activated[orderID] = struct{}{}

	return ActivationResult{
		OrderID:   orderID,
		Activated: true,
		Status:    "active",
	}, nil
}
