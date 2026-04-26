package employee

import (
	"context"

	"coffee-shop-backend/app/internal/models"
)

func (s *Service) ListEmployees(ctx context.Context) ([]models.Employee, error) {
	return s.repo.ListEmployees(ctx)
}
