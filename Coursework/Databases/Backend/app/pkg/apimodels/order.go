package apimodels

import "time"

type OrderShift struct {
	ID       string    `json:"id"`
	OpenedAt time.Time `json:"opened_at"`
}

type OrderEmployee struct {
	ID       string `json:"id"`
	FullName string `json:"full_name"`
}

type OrderItem struct {
	ID           string  `json:"id"`
	MenuItemID   string  `json:"menu_item_id"`
	MenuItemName string  `json:"menu_item_name"`
	Quantity     int     `json:"quantity"`
	Price        float64 `json:"price"`
	TotalPrice   float64 `json:"total_price"`
}

type OrderInventoryOperation struct {
	ID             string    `json:"id"`
	IngredientID   string    `json:"ingredient_id"`
	IngredientName string    `json:"ingredient_name"`
	ChangeAmount   float64   `json:"change_amount"`
	Unit           string    `json:"unit"`
	OperationType  string    `json:"operation_type"`
	CreatedAt      time.Time `json:"created_at"`
}

type Order struct {
	ID                  string                    `json:"id"`
	CreatedAt           time.Time                 `json:"created_at"`
	Shift               OrderShift                `json:"shift"`
	Employee            OrderEmployee             `json:"employee"`
	Items               []OrderItem               `json:"items"`
	InventoryOperations []OrderInventoryOperation `json:"inventory_operations,omitempty"`
	TotalPrice          float64                   `json:"total_price"`
}

type ListOrdersResponse struct {
	Items      []Order    `json:"items"`
	Pagination Pagination `json:"pagination"`
}

type GetOrderResponse struct {
	Order *Order `json:"order"`
}

type CreateOrderItemRequest struct {
	MenuItemID string `json:"menu_item_id"`
	Quantity   int    `json:"quantity"`
}

type CreateOrderRequest struct {
	Items []CreateOrderItemRequest `json:"items"`
}

type CreateOrderResponse struct {
	Order *Order `json:"order"`
}
