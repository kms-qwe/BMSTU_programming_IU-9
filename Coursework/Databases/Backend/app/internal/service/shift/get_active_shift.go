package shift

import (
	"context"

	"coffee-shop-backend/app/internal/models"
)

func (s *Service) GetActiveShift(ctx context.Context) (*models.Shift, error) {
	return s.shiftRepo.GetActiveShift(ctx, false)
}
