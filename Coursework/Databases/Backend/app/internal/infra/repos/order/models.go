package order

import (
	"time"

	"coffee-shop-backend/app/internal/models"
)

type orderDBModel struct {
	ID               string    `db:"id"`
	CreatedAt        time.Time `db:"created_at"`
	ShiftID          string    `db:"shift_id"`
	ShiftOpenedAt    time.Time `db:"shift_opened_at"`
	EmployeeID       string    `db:"employee_id"`
	EmployeeFullName string    `db:"employee_full_name"`
	TotalPrice       float64   `db:"total_price"`
}

type createdOrderDBModel struct {
	ID        string    `db:"id"`
	CreatedAt time.Time `db:"created_at"`
}

type orderItemDBModel struct {
	ID           string  `db:"id"`
	OrderID      string  `db:"order_id"`
	MenuItemID   string  `db:"menu_item_id"`
	MenuItemName string  `db:"menu_item_name"`
	Quantity     int     `db:"quantity"`
	Price        float64 `db:"price"`
}

type orderItemCreateDBModel struct {
	OrderID    string  `db:"order_id"`
	MenuItemID string  `db:"menu_item_id"`
	Quantity   int     `db:"quantity"`
	Price      float64 `db:"price"`
}

type inventoryOperationDBModel struct {
	ID               string    `db:"id"`
	CreatedAt        time.Time `db:"created_at"`
	OperationType    string    `db:"operation_type"`
	ChangeAmount     float64   `db:"change_amount"`
	IngredientID     string    `db:"ingredient_id"`
	IngredientName   string    `db:"ingredient_name"`
	IngredientUnit   string    `db:"ingredient_unit"`
	IngredientStock  float64   `db:"ingredient_current_stock"`
	OrderID          *string   `db:"order_id"`
	ShiftID          string    `db:"shift_id"`
	EmployeeID       string    `db:"employee_id"`
	EmployeeFullName string    `db:"employee_full_name"`
}

func toServiceOrder(item orderDBModel) models.Order {
	return models.Order{
		ID:        item.ID,
		CreatedAt: item.CreatedAt,
		Shift: models.OrderShift{
			ID:       item.ShiftID,
			OpenedAt: item.ShiftOpenedAt,
		},
		Employee: models.OrderEmployee{
			ID:       item.EmployeeID,
			FullName: item.EmployeeFullName,
		},
		TotalPrice: item.TotalPrice,
	}
}

func toServiceOrderItems(items []orderItemDBModel) []models.OrderItem {
	result := make([]models.OrderItem, 0, len(items))
	for _, item := range items {
		result = append(result, models.OrderItem{
			ID:           item.ID,
			MenuItemID:   item.MenuItemID,
			MenuItemName: item.MenuItemName,
			Quantity:     item.Quantity,
			Price:        item.Price,
			TotalPrice:   item.Price * float64(item.Quantity),
		})
	}

	return result
}

func toServiceInventoryOperations(items []inventoryOperationDBModel) []models.InventoryOperation {
	result := make([]models.InventoryOperation, 0, len(items))
	for _, item := range items {
		var orderRef *models.InventoryOperationOrder
		if item.OrderID != nil {
			orderRef = &models.InventoryOperationOrder{ID: *item.OrderID}
		}

		result = append(result, models.InventoryOperation{
			ID:            item.ID,
			CreatedAt:     item.CreatedAt,
			OperationType: item.OperationType,
			ChangeAmount:  item.ChangeAmount,
			Ingredient: models.InventoryOperationIngredient{
				ID:           item.IngredientID,
				Name:         item.IngredientName,
				Unit:         item.IngredientUnit,
				CurrentStock: item.IngredientStock,
			},
			Order: orderRef,
			Shift: models.InventoryOperationShift{
				ID: item.ShiftID,
			},
			Employee: models.InventoryOperationEmployee{
				ID:       item.EmployeeID,
				FullName: item.EmployeeFullName,
			},
		})
	}

	return result
}

func toDBOrderItemCreateModels(orderID string, items []models.OrderItemCreateRecord) []orderItemCreateDBModel {
	result := make([]orderItemCreateDBModel, 0, len(items))
	for _, item := range items {
		result = append(result, orderItemCreateDBModel{
			OrderID:    orderID,
			MenuItemID: item.MenuItemID,
			Quantity:   item.Quantity,
			Price:      item.Price,
		})
	}

	return result
}
