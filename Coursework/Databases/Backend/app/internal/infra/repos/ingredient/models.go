package ingredient

import "coffee-shop-backend/app/internal/models"

type ingredientDBModel struct {
	ID                 string  `db:"id"`
	Name               string  `db:"name"`
	Unit               string  `db:"unit"`
	CurrentStock       float64 `db:"current_stock"`
	IsActive           bool    `db:"is_active"`
	UsedInRecipesCount int     `db:"used_in_recipes_count"`
}

type ingredientUsageDBModel struct {
	MenuItemID   string  `db:"menu_item_id"`
	MenuItemName string  `db:"menu_item_name"`
	Amount       float64 `db:"amount"`
}

type createdIngredientDBModel struct {
	ID string `db:"id"`
}

func toServiceIngredient(item ingredientDBModel) models.Ingredient {
	return models.Ingredient{
		ID:                 item.ID,
		Name:               item.Name,
		Unit:               item.Unit,
		CurrentStock:       item.CurrentStock,
		IsActive:           item.IsActive,
		UsedInRecipesCount: item.UsedInRecipesCount,
	}
}

func toServiceIngredients(items []ingredientDBModel) []models.Ingredient {
	result := make([]models.Ingredient, 0, len(items))
	for _, item := range items {
		result = append(result, toServiceIngredient(item))
	}

	return result
}

func toServiceIngredientUsage(items []ingredientUsageDBModel) []models.IngredientUsedInRecipe {
	result := make([]models.IngredientUsedInRecipe, 0, len(items))
	for _, item := range items {
		result = append(result, models.IngredientUsedInRecipe{
			MenuItemID:   item.MenuItemID,
			MenuItemName: item.MenuItemName,
			Amount:       item.Amount,
		})
	}

	return result
}
