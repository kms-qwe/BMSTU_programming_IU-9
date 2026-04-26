package models

type RecipeItem struct {
	IngredientID       string
	IngredientName     string
	Unit               string
	Amount             float64
	IngredientIsActive bool
}

type MenuItem struct {
	ID       string
	Name     string
	Price    float64
	IsActive bool
	Recipe   []RecipeItem
}

type CreateMenuItemRecipeInput struct {
	IngredientID string
	Amount       float64
}

type CreateMenuItemParams struct {
	Name   string
	Price  float64
	Recipe []CreateMenuItemRecipeInput
}

type UpdateMenuItemPriceParams struct {
	MenuItemID string
	Price      float64
}

type SetMenuItemActivityParams struct {
	MenuItemID string
	IsActive   bool
}
