package inventory

import (
	"context"

	"coffee-shop-backend/app/internal/models"
	"coffee-shop-backend/app/internal/utils"
)

type IInventoryRepo interface {
	ListInventoryOperations(ctx context.Context, filter models.InventoryOperationFilter) ([]models.InventoryOperation, error)
	CountInventoryOperations(ctx context.Context, filter models.InventoryOperationFilter) (int, error)
	CreateInventoryOperation(ctx context.Context, params models.InventoryOperationCreateRecord) (*models.InventoryOperation, error)
}
type IShiftRepo interface {
	GetActiveShift(ctx context.Context, forUpdate bool) (*models.Shift, error)
}
type IIngredientRepo interface {
	GetIngredientsByIDs(ctx context.Context, ingredientIDs []string, forUpdate bool) ([]models.Ingredient, error)
	AdjustIngredientStock(ctx context.Context, ingredientID string, delta float64) error
}

type Service struct {
	inventoryRepo  IInventoryRepo
	shiftRepo      IShiftRepo
	ingredientRepo IIngredientRepo
	txManager      utils.TxManager
}

func New(inventoryRepo IInventoryRepo, shiftRepo IShiftRepo, ingredientRepo IIngredientRepo, txManager utils.TxManager) *Service {
	return &Service{inventoryRepo: inventoryRepo, shiftRepo: shiftRepo, ingredientRepo: ingredientRepo, txManager: txManager}
}
