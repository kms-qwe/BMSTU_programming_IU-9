package models

import "time"

type Pagination struct {
	Page       int
	PageSize   int
	TotalItems int
	TotalPages int
}

type PaginatedResult[T any] struct {
	Items      []T
	Pagination Pagination
}

type OrderFilter struct {
	Page       int
	PageSize   int
	EmployeeID string
	From       *time.Time
	To         *time.Time
}

type InventoryOperationFilter struct {
	Page          int
	PageSize      int
	EmployeeID    string
	IngredientID  string
	OperationType string
	From          *time.Time
	To            *time.Time
}

type ReportFilter struct {
	From       time.Time
	To         time.Time
	EmployeeID *string
}

type InventoryReportFilter struct {
	From          time.Time
	To            time.Time
	EmployeeID    *string
	IngredientID  *string
	OperationType *string
}
