package apimodels

import "time"

type InventoryOperationIngredient struct {
	ID           string  `json:"id"`
	Name         string  `json:"name"`
	Unit         string  `json:"unit"`
	CurrentStock float64 `json:"current_stock,omitempty"`
}

type InventoryOrderReference struct {
	ID string `json:"id"`
}

type InventoryShiftReference struct {
	ID string `json:"id"`
}

type InventoryEmployeeReference struct {
	ID       string `json:"id"`
	FullName string `json:"full_name"`
}

type InventoryOperation struct {
	ID            string                       `json:"id"`
	CreatedAt     time.Time                    `json:"created_at"`
	OperationType string                       `json:"operation_type"`
	ChangeAmount  float64                      `json:"change_amount"`
	Ingredient    InventoryOperationIngredient `json:"ingredient"`
	Order         *InventoryOrderReference     `json:"order"`
	Shift         InventoryShiftReference      `json:"shift"`
	Employee      InventoryEmployeeReference   `json:"employee"`
}

type ListOperationsResponse struct {
	Items      []InventoryOperation `json:"items"`
	Pagination Pagination           `json:"pagination"`
}

type CreateOperationRequest struct {
	IngredientID  string  `json:"ingredient_id"`
	OperationType string  `json:"operation_type"`
	Amount        float64 `json:"amount"`
}

type CreateOperationResponse struct {
	Operation *InventoryOperation `json:"operation"`
}
