package models

import "time"

type InventoryOperationIngredient struct {
	ID           string
	Name         string
	Unit         string
	CurrentStock float64
}

type InventoryOperationOrder struct {
	ID string
}

type InventoryOperationShift struct {
	ID string
}

type InventoryOperationEmployee struct {
	ID       string
	FullName string
}

type InventoryOperation struct {
	ID            string
	CreatedAt     time.Time
	OperationType string
	ChangeAmount  float64
	Ingredient    InventoryOperationIngredient
	Order         *InventoryOperationOrder
	Shift         InventoryOperationShift
	Employee      InventoryOperationEmployee
}

type CreateInventoryOperationParams struct {
	IngredientID  string
	OperationType string
	Amount        float64
}

type InventoryOperationCreateRecord struct {
	ShiftID       string
	IngredientID  string
	OrderID       *string
	ChangeAmount  float64
	OperationType string
}
