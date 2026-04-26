package menuitem

import (
	"context"
	"strings"

	"coffee-shop-backend/app/internal/models"
	"coffee-shop-backend/app/internal/utils"
)

func (s *Service) GetMenuItem(ctx context.Context, menuItemID string) (*models.MenuItem, error) {
	if strings.TrimSpace(menuItemID) == "" {
		return nil, utils.ValidationError("menu_item_id is required", nil)
	}
	item, err := s.menuRepo.GetMenuItem(ctx, menuItemID)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, utils.NotFoundError("Пункт меню не найден")
	}
	return item, nil
}
