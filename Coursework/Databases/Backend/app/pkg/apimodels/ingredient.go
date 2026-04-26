package apimodels

type Ingredient struct {
	ID                 string  `json:"id"`
	Name               string  `json:"name"`
	Unit               string  `json:"unit"`
	CurrentStock       float64 `json:"current_stock"`
	IsActive           bool    `json:"is_active"`
	UsedInRecipesCount int     `json:"used_in_recipes_count"`
}

type IngredientUsage struct {
	MenuItemID   string  `json:"menu_item_id"`
	MenuItemName string  `json:"menu_item_name"`
	Amount       float64 `json:"amount"`
}

type IngredientDetails struct {
	ID                 string            `json:"id"`
	Name               string            `json:"name"`
	Unit               string            `json:"unit"`
	CurrentStock       float64           `json:"current_stock"`
	IsActive           bool              `json:"is_active"`
	UsedInRecipesCount int               `json:"used_in_recipes_count"`
	UsedInRecipes      []IngredientUsage `json:"used_in_recipes"`
}

type ListIngredientsResponse struct {
	Ingredients []Ingredient `json:"ingredients"`
}

type GetIngredientResponse struct {
	Ingredient *IngredientDetails `json:"ingredient"`
}

type CreateIngredientRequest struct {
	Name         string   `json:"name"`
	Unit         string   `json:"unit"`
	InitialStock *float64 `json:"initial_stock"`
}

type CreateIngredientResponse struct {
	Ingredient Ingredient `json:"ingredient"`
}

type SetIngredientActivityRequest struct {
	IsActive bool `json:"is_active"`
}

type SetIngredientActivityResponse struct {
	Ingredient Ingredient `json:"ingredient"`
}
