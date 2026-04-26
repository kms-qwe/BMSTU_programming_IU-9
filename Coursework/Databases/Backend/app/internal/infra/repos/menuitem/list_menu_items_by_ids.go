package menuitem

import (
	"context"

	"coffee-shop-backend/app/internal/models"
)

func (r *Repo) ListMenuItemsByIDs(ctx context.Context, menuItemIDs []string) ([]models.MenuItem, error) {
	return r.loadMenuItems(ctx, true, menuItemIDs)
}
