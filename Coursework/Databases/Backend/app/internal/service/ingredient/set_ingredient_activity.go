package ingredient

import (
	"context"
	"strings"

	"coffee-shop-backend/app/internal/models"
	"coffee-shop-backend/app/internal/utils"
)

func (s *Service) SetIngredientActivity(ctx context.Context, params models.SetIngredientActivityParams) (*models.Ingredient, error) {
	if strings.TrimSpace(params.IngredientID) == "" {
		return nil, utils.ValidationError("ingredient_id is required", nil)
	}
	if !params.IsActive {
		usages, err := s.ingredientRepo.ListIngredientUsage(ctx, params.IngredientID)
		if err != nil {
			return nil, err
		}
		if len(usages) > 0 {
			return nil, utils.BusinessError(utils.ErrorCodeIngredientUsedRecipe, "Ингредиент нельзя сделать неактивным, потому что он используется в рецептах", map[string]any{"used_in_recipes": usages})
		}
	}
	ingredient, err := s.ingredientRepo.SetIngredientActive(ctx, params.IngredientID, params.IsActive)
	if err != nil {
		return nil, err
	}
	if ingredient == nil {
		return nil, utils.NotFoundError("Ингредиент не найден")
	}
	return ingredient, nil
}
