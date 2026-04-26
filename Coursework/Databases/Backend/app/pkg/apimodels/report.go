package apimodels

type TimeRangeRequest struct {
	From       string  `json:"from"`
	To         string  `json:"to"`
	EmployeeID *string `json:"employee_id"`
}

type InventoryReportRequest struct {
	From          string  `json:"from"`
	To            string  `json:"to"`
	EmployeeID    *string `json:"employee_id"`
	IngredientID  *string `json:"ingredient_id"`
	OperationType *string `json:"operation_type"`
}
