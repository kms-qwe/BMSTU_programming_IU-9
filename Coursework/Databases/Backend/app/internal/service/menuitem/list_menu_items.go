package menuitem

import (
	"context"

	"coffee-shop-backend/app/internal/models"
)

func (s *Service) ListMenuItems(ctx context.Context, includeInactive bool) ([]models.MenuItem, error) {
	return s.menuRepo.ListMenuItems(ctx, includeInactive)
}
