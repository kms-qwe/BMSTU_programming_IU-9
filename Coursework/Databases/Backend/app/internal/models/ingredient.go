package models

type IngredientUsedInRecipe struct {
	MenuItemID   string
	MenuItemName string
	Amount       float64
}

type Ingredient struct {
	ID                 string
	Name               string
	Unit               string
	CurrentStock       float64
	IsActive           bool
	UsedInRecipesCount int
	UsedInRecipes      []IngredientUsedInRecipe
}

type CreateIngredientParams struct {
	Name         string
	Unit         string
	InitialStock float64
}

type SetIngredientActivityParams struct {
	IngredientID string
	IsActive     bool
}
