package menuitem

import (
	"context"

	"coffee-shop-backend/app/internal/models"
)

func (r *Repo) ListMenuItems(ctx context.Context, includeInactive bool) ([]models.MenuItem, error) {
	return r.loadMenuItems(ctx, includeInactive, nil)
}
