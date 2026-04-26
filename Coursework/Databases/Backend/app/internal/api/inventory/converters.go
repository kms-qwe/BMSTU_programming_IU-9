package inventory

import (
	"coffee-shop-backend/app/internal/models"
	"coffee-shop-backend/app/pkg/apimodels"
)

func toInventoryOperationResponse(item models.InventoryOperation) apimodels.InventoryOperation {
	var order *apimodels.InventoryOrderReference
	if item.Order != nil {
		order = &apimodels.InventoryOrderReference{ID: item.Order.ID}
	}
	return apimodels.InventoryOperation{
		ID: item.ID, CreatedAt: item.CreatedAt, OperationType: item.OperationType, ChangeAmount: item.ChangeAmount,
		Ingredient: apimodels.InventoryOperationIngredient{ID: item.Ingredient.ID, Name: item.Ingredient.Name, Unit: item.Ingredient.Unit, CurrentStock: item.Ingredient.CurrentStock},
		Order:      order,
		Shift:      apimodels.InventoryShiftReference{ID: item.Shift.ID},
		Employee:   apimodels.InventoryEmployeeReference{ID: item.Employee.ID, FullName: item.Employee.FullName},
	}
}

func toInventoryOperationResponsePtr(item *models.InventoryOperation) *apimodels.InventoryOperation {
	if item == nil {
		return nil
	}
	result := toInventoryOperationResponse(*item)
	return &result
}
