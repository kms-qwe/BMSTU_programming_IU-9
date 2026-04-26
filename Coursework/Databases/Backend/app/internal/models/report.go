package models

import "time"

type SalesSummary struct {
	OrdersCount int
	Revenue     float64
	AverageBill float64
}

type SalesMenuItemRow struct {
	MenuItemName  string
	Quantity      int
	Revenue       float64
	SalesSharePct float64
}

type SalesOrderRow struct {
	CreatedAt    time.Time
	OrderID      string
	EmployeeName string
	TotalPrice   float64
}

type SalesReportData struct {
	Summary   SalesSummary
	MenuItems []SalesMenuItemRow
	Orders    []SalesOrderRow
}

type EmployeeSummaryRow struct {
	EmployeeName string
	ShiftCount   int
	OrdersCount  int
	Revenue      float64
	AverageBill  float64
}

type EmployeeOrderRow struct {
	CreatedAt    time.Time
	EmployeeName string
	OrderID      string
	TotalPrice   float64
}

type EmployeeShiftRow struct {
	EmployeeName string
	OpenedAt     time.Time
	ClosedAt     *time.Time
	Duration     string
}

type EmployeeReportData struct {
	Summaries []EmployeeSummaryRow
	Orders    []EmployeeOrderRow
	Shifts    []EmployeeShiftRow
}

type InventorySummaryRow struct {
	IngredientName   string
	Unit             string
	Restocked        float64
	AutoWrittenOff   float64
	ManualWrittenOff float64
	CurrentStock     float64
}

type InventoryReportData struct {
	Summaries  []InventorySummaryRow
	Operations []InventoryOperation
}

type MenuPopularityRow struct {
	Place           int
	MenuItemName    string
	SalesCount      int
	Revenue         float64
	RevenueSharePct float64
}

type MenuPopularityReportData struct {
	Rows []MenuPopularityRow
}
