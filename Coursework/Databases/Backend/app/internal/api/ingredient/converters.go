package ingredient

import (
	"coffee-shop-backend/app/internal/models"
	"coffee-shop-backend/app/pkg/apimodels"
)

func toIngredientResponse(item models.Ingredient) apimodels.Ingredient {
	return apimodels.Ingredient{
		ID:                 item.ID,
		Name:               item.Name,
		Unit:               item.Unit,
		CurrentStock:       item.CurrentStock,
		IsActive:           item.IsActive,
		UsedInRecipesCount: item.UsedInRecipesCount,
	}
}

func toIngredientDetailsResponse(item *models.Ingredient) *apimodels.IngredientDetails {
	if item == nil {
		return nil
	}
	used := make([]apimodels.IngredientUsage, 0, len(item.UsedInRecipes))
	for _, usage := range item.UsedInRecipes {
		used = append(used, apimodels.IngredientUsage{
			MenuItemID:   usage.MenuItemID,
			MenuItemName: usage.MenuItemName,
			Amount:       usage.Amount,
		})
	}
	return &apimodels.IngredientDetails{
		ID:                 item.ID,
		Name:               item.Name,
		Unit:               item.Unit,
		CurrentStock:       item.CurrentStock,
		IsActive:           item.IsActive,
		UsedInRecipesCount: item.UsedInRecipesCount,
		UsedInRecipes:      used,
	}
}
