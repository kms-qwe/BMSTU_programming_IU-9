package menuitem

import (
	"coffee-shop-backend/app/internal/models"
)

type menuItemRowDBModel struct {
	ID                 string   `db:"id"`
	Name               string   `db:"name"`
	Price              float64  `db:"price"`
	IsActive           bool     `db:"is_active"`
	IngredientID       *string  `db:"ingredient_id"`
	IngredientName     *string  `db:"ingredient_name"`
	Unit               *string  `db:"unit"`
	Amount             *float64 `db:"amount"`
	IngredientIsActive bool     `db:"ingredient_is_active"`
}

type createdMenuItemDBModel struct {
	ID string `db:"id"`
}

func toServiceMenuItems(rows []menuItemRowDBModel) []models.MenuItem {
	menuMap := make(map[string]*models.MenuItem)
	order := make([]string, 0)

	for _, row := range rows {
		item, exists := menuMap[row.ID]
		if !exists {
			item = &models.MenuItem{
				ID:       row.ID,
				Name:     row.Name,
				Price:    row.Price,
				IsActive: row.IsActive,
				Recipe:   []models.RecipeItem{},
			}
			menuMap[row.ID] = item
			order = append(order, row.ID)
		}

		if row.IngredientID != nil {
			recipeItem := models.RecipeItem{
				IngredientID:       *row.IngredientID,
				IngredientIsActive: row.IngredientIsActive,
			}
			if row.IngredientName != nil {
				recipeItem.IngredientName = *row.IngredientName
			}
			if row.Unit != nil {
				recipeItem.Unit = *row.Unit
			}
			if row.Amount != nil {
				recipeItem.Amount = *row.Amount
			}
			item.Recipe = append(item.Recipe, recipeItem)
		}
	}

	result := make([]models.MenuItem, 0, len(order))
	for _, id := range order {
		result = append(result, *menuMap[id])
	}

	return result
}
