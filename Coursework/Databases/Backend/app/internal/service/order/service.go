package order

import (
	"context"

	"coffee-shop-backend/app/internal/models"
	"coffee-shop-backend/app/internal/utils"
)

type IShiftRepo interface {
	GetActiveShift(ctx context.Context, forUpdate bool) (*models.Shift, error)
}
type IMenuItemRepo interface {
	ListMenuItemsByIDs(ctx context.Context, menuItemIDs []string) ([]models.MenuItem, error)
}
type IIngredientRepo interface {
	GetIngredientsByIDs(ctx context.Context, ingredientIDs []string, forUpdate bool) ([]models.Ingredient, error)
	AdjustIngredientStock(ctx context.Context, ingredientID string, delta float64) error
}
type IOrderRepo interface {
	CreateOrder(ctx context.Context, shiftID string) (*models.Order, error)
	AddOrderItems(ctx context.Context, orderID string, items []models.OrderItemCreateRecord) error
	ListOrders(ctx context.Context, filter models.OrderFilter) ([]models.Order, error)
	CountOrders(ctx context.Context, filter models.OrderFilter) (int, error)
	GetOrder(ctx context.Context, orderID string) (*models.Order, error)
}
type IInventoryRepo interface {
	CreateInventoryOperation(ctx context.Context, params models.InventoryOperationCreateRecord) (*models.InventoryOperation, error)
}

type Service struct {
	shiftRepo      IShiftRepo
	menuRepo       IMenuItemRepo
	ingredientRepo IIngredientRepo
	orderRepo      IOrderRepo
	inventoryRepo  IInventoryRepo
	txManager      utils.TxManager
}

func New(shiftRepo IShiftRepo, menuRepo IMenuItemRepo, ingredientRepo IIngredientRepo, orderRepo IOrderRepo, inventoryRepo IInventoryRepo, txManager utils.TxManager) *Service {
	return &Service{shiftRepo: shiftRepo, menuRepo: menuRepo, ingredientRepo: ingredientRepo, orderRepo: orderRepo, inventoryRepo: inventoryRepo, txManager: txManager}
}
