package order

import (
	"context"
	"math"

	"coffee-shop-backend/app/internal/models"
)

func (s *Service) ListOrders(ctx context.Context, filter models.OrderFilter) (*models.PaginatedResult[models.Order], error) {
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.PageSize <= 0 {
		filter.PageSize = 20
	}
	total, err := s.orderRepo.CountOrders(ctx, filter)
	if err != nil {
		return nil, err
	}
	items, err := s.orderRepo.ListOrders(ctx, filter)
	if err != nil {
		return nil, err
	}
	return &models.PaginatedResult[models.Order]{Items: items, Pagination: models.Pagination{
		Page: filter.Page, PageSize: filter.PageSize, TotalItems: total, TotalPages: int(math.Ceil(float64(total) / float64(filter.PageSize))),
	}}, nil
}
