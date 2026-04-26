package menuitem

import (
	"context"

	"coffee-shop-backend/app/internal/models"
)

func (r *Repo) GetMenuItem(ctx context.Context, menuItemID string) (*models.MenuItem, error) {
	items, err := r.loadMenuItems(ctx, true, []string{menuItemID})
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, nil
	}

	return &items[0], nil
}
