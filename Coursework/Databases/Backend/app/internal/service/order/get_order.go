package order

import (
	"context"
	"strings"

	"coffee-shop-backend/app/internal/models"
	"coffee-shop-backend/app/internal/utils"
)

func (s *Service) GetOrder(ctx context.Context, orderID string) (*models.Order, error) {
	if strings.TrimSpace(orderID) == "" {
		return nil, utils.ValidationError("order_id is required", nil)
	}
	item, err := s.orderRepo.GetOrder(ctx, orderID)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, utils.NotFoundError("Заказ не найден")
	}
	return item, nil
}
