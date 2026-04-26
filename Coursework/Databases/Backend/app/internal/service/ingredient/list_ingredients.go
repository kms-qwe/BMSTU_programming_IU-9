package ingredient

import (
	"context"
	"strings"

	"coffee-shop-backend/app/internal/models"
	"coffee-shop-backend/app/internal/utils"
)

func (s *Service) ListIngredients(ctx context.Context, includeInactive bool) ([]models.Ingredient, error) {
	return s.ingredientRepo.ListIngredients(ctx, includeInactive)
}
func (s *Service) GetIngredient(ctx context.Context, ingredientID string) (*models.Ingredient, error) {
	if strings.TrimSpace(ingredientID) == "" {
		return nil, utils.ValidationError("ingredient_id is required", nil)
	}
	ingredient, err := s.ingredientRepo.GetIngredient(ctx, ingredientID)
	if err != nil {
		return nil, err
	}
	if ingredient == nil {
		return nil, utils.NotFoundError("Ингредиент не найден")
	}
	usages, err := s.ingredientRepo.ListIngredientUsage(ctx, ingredientID)
	if err != nil {
		return nil, err
	}
	ingredient.UsedInRecipes = usages
	return ingredient, nil
}
