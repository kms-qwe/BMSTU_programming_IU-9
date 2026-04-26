package menuitem

import (
	"context"

	"coffee-shop-backend/app/internal/models"
)

type IMenuItemRepo interface {
	ListMenuItems(ctx context.Context, includeInactive bool) ([]models.MenuItem, error)
	GetMenuItem(ctx context.Context, menuItemID string) (*models.MenuItem, error)
	ListMenuItemsByIDs(ctx context.Context, menuItemIDs []string) ([]models.MenuItem, error)
	CreateMenuItem(ctx context.Context, params models.CreateMenuItemParams) (*models.MenuItem, error)
	UpdateMenuItemPrice(ctx context.Context, menuItemID string, price float64) (*models.MenuItem, error)
	SetMenuItemActive(ctx context.Context, menuItemID string, isActive bool) (*models.MenuItem, error)
}

type IIngredientRepo interface {
	GetIngredientsByIDs(ctx context.Context, ingredientIDs []string, forUpdate bool) ([]models.Ingredient, error)
}

type Service struct {
	menuRepo       IMenuItemRepo
	ingredientRepo IIngredientRepo
}

func New(menuRepo IMenuItemRepo, ingredientRepo IIngredientRepo) *Service {
	return &Service{menuRepo: menuRepo, ingredientRepo: ingredientRepo}
}
