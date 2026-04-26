package order

import (
	"coffee-shop-backend/app/internal/models"
	"coffee-shop-backend/app/pkg/apimodels"
)

func toOrderResponse(item models.Order) apimodels.Order {
	orderItems := make([]apimodels.OrderItem, 0, len(item.Items))
	for _, orderItem := range item.Items {
		orderItems = append(orderItems, apimodels.OrderItem{
			ID: orderItem.ID, MenuItemID: orderItem.MenuItemID, MenuItemName: orderItem.MenuItemName,
			Quantity: orderItem.Quantity, Price: orderItem.Price, TotalPrice: orderItem.TotalPrice,
		})
	}
	ops := make([]apimodels.OrderInventoryOperation, 0, len(item.InventoryOperations))
	for _, op := range item.InventoryOperations {
		ops = append(ops, apimodels.OrderInventoryOperation{
			ID: op.ID, IngredientID: op.Ingredient.ID, IngredientName: op.Ingredient.Name,
			ChangeAmount: op.ChangeAmount, Unit: op.Ingredient.Unit, OperationType: op.OperationType, CreatedAt: op.CreatedAt,
		})
	}
	return apimodels.Order{
		ID: item.ID, CreatedAt: item.CreatedAt,
		Shift:    apimodels.OrderShift{ID: item.Shift.ID, OpenedAt: item.Shift.OpenedAt},
		Employee: apimodels.OrderEmployee{ID: item.Employee.ID, FullName: item.Employee.FullName},
		Items:    orderItems, InventoryOperations: ops, TotalPrice: item.TotalPrice,
	}
}

func toOrderResponsePtr(item *models.Order) *apimodels.Order {
	if item == nil {
		return nil
	}
	result := toOrderResponse(*item)
	return &result
}
