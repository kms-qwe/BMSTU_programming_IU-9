package menuitem

import (
	"context"
	"strings"

	"coffee-shop-backend/app/internal/models"
	"coffee-shop-backend/app/internal/utils"
)

func (s *Service) CreateMenuItem(ctx context.Context, params models.CreateMenuItemParams) (*models.MenuItem, error) {
	params.Name = strings.TrimSpace(params.Name)
	if params.Name == "" || params.Price < 0 || len(params.Recipe) == 0 {
		return nil, utils.ValidationError("invalid menu item payload", nil)
	}
	ids := make([]string, 0, len(params.Recipe))
	seen := map[string]struct{}{}
	for _, recipeItem := range params.Recipe {
		if strings.TrimSpace(recipeItem.IngredientID) == "" || recipeItem.Amount <= 0 {
			return nil, utils.ValidationError("recipe items must contain ingredient_id and positive amount", nil)
		}
		if _, ok := seen[recipeItem.IngredientID]; ok {
			return nil, utils.ValidationError("duplicate ingredients in recipe are not allowed", nil)
		}
		seen[recipeItem.IngredientID] = struct{}{}
		ids = append(ids, recipeItem.IngredientID)
	}
	ingredients, err := s.ingredientRepo.GetIngredientsByIDs(ctx, ids, false)
	if err != nil {
		return nil, err
	}
	if len(ingredients) != len(ids) {
		return nil, utils.NotFoundError("Не все ингредиенты найдены")
	}
	for _, ingredient := range ingredients {
		if !ingredient.IsActive {
			return nil, utils.BusinessError(utils.ErrorCodeInactiveIngredient, "Нельзя использовать неактивный ингредиент в рецепте", nil)
		}
	}
	return s.menuRepo.CreateMenuItem(ctx, params)
}
