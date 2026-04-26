package models

import "time"

type OrderShift struct {
	ID       string
	OpenedAt time.Time
}

type OrderEmployee struct {
	ID       string
	FullName string
}

type OrderItem struct {
	ID           string
	MenuItemID   string
	MenuItemName string
	Quantity     int
	Price        float64
	TotalPrice   float64
}

type Order struct {
	ID                  string
	CreatedAt           time.Time
	Shift               OrderShift
	Employee            OrderEmployee
	Items               []OrderItem
	InventoryOperations []InventoryOperation
	TotalPrice          float64
}

type CreateOrderItemInput struct {
	MenuItemID string
	Quantity   int
}

type CreateOrderParams struct {
	Items []CreateOrderItemInput
}

type OrderItemCreateRecord struct {
	MenuItemID string
	Quantity   int
	Price      float64
}
