package apimodels

type RecipeItem struct {
	IngredientID       string  `json:"ingredient_id"`
	IngredientName     string  `json:"ingredient_name"`
	Unit               string  `json:"unit"`
	Amount             float64 `json:"amount"`
	IngredientIsActive bool    `json:"ingredient_is_active"`
}

type MenuItem struct {
	ID       string       `json:"id"`
	Name     string       `json:"name"`
	Price    float64      `json:"price"`
	IsActive bool         `json:"is_active"`
	Recipe   []RecipeItem `json:"recipe"`
}

type ListMenuItemsResponse struct {
	MenuItems []MenuItem `json:"menu_items"`
}

type GetMenuItemResponse struct {
	MenuItem *MenuItem `json:"menu_item"`
}

type CreateMenuItemRecipeItem struct {
	IngredientID string  `json:"ingredient_id"`
	Amount       float64 `json:"amount"`
}

type CreateMenuItemRequest struct {
	Name   string                     `json:"name"`
	Price  float64                    `json:"price"`
	Recipe []CreateMenuItemRecipeItem `json:"recipe"`
}

type CreateMenuItemResponse struct {
	MenuItem *MenuItem `json:"menu_item"`
}

type UpdateMenuItemPriceRequest struct {
	Price float64 `json:"price"`
}

type UpdateMenuItemPriceResponse struct {
	MenuItem *MenuItem `json:"menu_item"`
}

type SetMenuItemActivityRequest struct {
	IsActive bool `json:"is_active"`
}

type SetMenuItemActivityResponse struct {
	MenuItem *MenuItem `json:"menu_item"`
}
