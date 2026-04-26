package menuitem

import (
	"context"
	"strings"

	"coffee-shop-backend/app/internal/models"
	"coffee-shop-backend/app/internal/utils"
)

func (s *Service) SetMenuItemActivity(ctx context.Context, params models.SetMenuItemActivityParams) (*models.MenuItem, error) {
	if strings.TrimSpace(params.MenuItemID) == "" {
		return nil, utils.ValidationError("menu_item_id is required", nil)
	}
	item, err := s.menuRepo.SetMenuItemActive(ctx, params.MenuItemID, params.IsActive)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, utils.NotFoundError("Пункт меню не найден")
	}
	return item, nil
}
