package inventory

import (
	"context"
	"math"

	"coffee-shop-backend/app/internal/models"
)

func (s *Service) ListInventoryOperations(ctx context.Context, filter models.InventoryOperationFilter) (*models.PaginatedResult[models.InventoryOperation], error) {
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.PageSize <= 0 {
		filter.PageSize = 20
	}
	total, err := s.inventoryRepo.CountInventoryOperations(ctx, filter)
	if err != nil {
		return nil, err
	}
	items, err := s.inventoryRepo.ListInventoryOperations(ctx, filter)
	if err != nil {
		return nil, err
	}
	return &models.PaginatedResult[models.InventoryOperation]{Items: items, Pagination: models.Pagination{Page: filter.Page, PageSize: filter.PageSize, TotalItems: total, TotalPages: int(math.Ceil(float64(total) / float64(filter.PageSize)))}}, nil
}
