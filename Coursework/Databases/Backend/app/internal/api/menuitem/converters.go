package menuitem

import (
	"coffee-shop-backend/app/internal/models"
	"coffee-shop-backend/app/pkg/apimodels"
)

func toMenuItemResponse(item models.MenuItem) apimodels.MenuItem {
	recipe := make([]apimodels.RecipeItem, 0, len(item.Recipe))
	for _, recipeItem := range item.Recipe {
		recipe = append(recipe, apimodels.RecipeItem{
			IngredientID:       recipeItem.IngredientID,
			IngredientName:     recipeItem.IngredientName,
			Unit:               recipeItem.Unit,
			Amount:             recipeItem.Amount,
			IngredientIsActive: recipeItem.IngredientIsActive,
		})
	}
	return apimodels.MenuItem{
		ID:       item.ID,
		Name:     item.Name,
		Price:    item.Price,
		IsActive: item.IsActive,
		Recipe:   recipe,
	}
}

func toMenuItemResponsePtr(item *models.MenuItem) *apimodels.MenuItem {
	if item == nil {
		return nil
	}
	result := toMenuItemResponse(*item)
	return &result
}
