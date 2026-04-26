package menuitem

import (
	"context"
	"strings"

	"coffee-shop-backend/app/internal/models"
	"coffee-shop-backend/app/internal/utils"
)

func (s *Service) UpdateMenuItemPrice(ctx context.Context, params models.UpdateMenuItemPriceParams) (*models.MenuItem, error) {
	if strings.TrimSpace(params.MenuItemID) == "" || params.Price < 0 {
		return nil, utils.ValidationError("invalid price payload", nil)
	}
	item, err := s.menuRepo.UpdateMenuItemPrice(ctx, params.MenuItemID, params.Price)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, utils.NotFoundError("Пункт меню не найден")
	}
	return item, nil
}
