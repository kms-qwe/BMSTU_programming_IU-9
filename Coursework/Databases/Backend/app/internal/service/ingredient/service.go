package ingredient

import (
	"context"

	"coffee-shop-backend/app/internal/models"
	"coffee-shop-backend/app/internal/utils"
)

type IIngredientRepo interface {
	ListIngredients(ctx context.Context, includeInactive bool) ([]models.Ingredient, error)
	GetIngredient(ctx context.Context, ingredientID string) (*models.Ingredient, error)
	GetIngredientsByIDs(ctx context.Context, ingredientIDs []string, forUpdate bool) ([]models.Ingredient, error)
	ListIngredientUsage(ctx context.Context, ingredientID string) ([]models.IngredientUsedInRecipe, error)
	CreateIngredient(ctx context.Context, params models.CreateIngredientParams) (*models.Ingredient, error)
	SetIngredientActive(ctx context.Context, ingredientID string, isActive bool) (*models.Ingredient, error)
}

type IShiftRepo interface {
	GetActiveShift(ctx context.Context, forUpdate bool) (*models.Shift, error)
}

type IInventoryRepo interface {
	CreateInventoryOperation(ctx context.Context, params models.InventoryOperationCreateRecord) (*models.InventoryOperation, error)
}

type Service struct {
	ingredientRepo IIngredientRepo
	shiftRepo      IShiftRepo
	inventoryRepo  IInventoryRepo
	txManager      utils.TxManager
}

func New(ingredientRepo IIngredientRepo, shiftRepo IShiftRepo, inventoryRepo IInventoryRepo, txManager utils.TxManager) *Service {
	return &Service{ingredientRepo: ingredientRepo, shiftRepo: shiftRepo, inventoryRepo: inventoryRepo, txManager: txManager}
}
