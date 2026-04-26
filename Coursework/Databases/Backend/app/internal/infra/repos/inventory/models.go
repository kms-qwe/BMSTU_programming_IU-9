package inventory

import (
	"time"

	"coffee-shop-backend/app/internal/models"
)

type inventoryOperationDBModel struct {
	ID                     string    `db:"id"`
	CreatedAt              time.Time `db:"created_at"`
	OperationType          string    `db:"operation_type"`
	ChangeAmount           float64   `db:"change_amount"`
	IngredientID           string    `db:"ingredient_id"`
	IngredientName         string    `db:"ingredient_name"`
	IngredientUnit         string    `db:"ingredient_unit"`
	IngredientCurrentStock float64   `db:"ingredient_current_stock"`
	OrderID                *string   `db:"order_id"`
	ShiftID                string    `db:"shift_id"`
	EmployeeID             string    `db:"employee_id"`
	EmployeeFullName       string    `db:"employee_full_name"`
}

type createdInventoryOperationDBModel struct {
	ID string `db:"id"`
}

func toServiceInventoryOperation(item inventoryOperationDBModel) models.InventoryOperation {
	var orderRef *models.InventoryOperationOrder
	if item.OrderID != nil {
		orderRef = &models.InventoryOperationOrder{ID: *item.OrderID}
	}

	return models.InventoryOperation{
		ID:            item.ID,
		CreatedAt:     item.CreatedAt,
		OperationType: item.OperationType,
		ChangeAmount:  item.ChangeAmount,
		Ingredient: models.InventoryOperationIngredient{
			ID:           item.IngredientID,
			Name:         item.IngredientName,
			Unit:         item.IngredientUnit,
			CurrentStock: item.IngredientCurrentStock,
		},
		Order: orderRef,
		Shift: models.InventoryOperationShift{
			ID: item.ShiftID,
		},
		Employee: models.InventoryOperationEmployee{
			ID:       item.EmployeeID,
			FullName: item.EmployeeFullName,
		},
	}
}

func toServiceInventoryOperations(items []inventoryOperationDBModel) []models.InventoryOperation {
	result := make([]models.InventoryOperation, 0, len(items))
	for _, item := range items {
		result = append(result, toServiceInventoryOperation(item))
	}

	return result
}
