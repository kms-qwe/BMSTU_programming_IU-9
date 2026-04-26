package report

import (
	"time"

	"coffee-shop-backend/app/internal/models"
)

type salesSummaryDBModel struct {
	OrdersCount int     `db:"orders_count"`
	Revenue     float64 `db:"revenue"`
	AverageBill float64 `db:"average_bill"`
}

type salesMenuItemRowDBModel struct {
	MenuItemName  string  `db:"menu_item_name"`
	Quantity      int     `db:"quantity"`
	Revenue       float64 `db:"revenue"`
	SalesSharePct float64 `db:"sales_share_pct"`
}

type salesOrderRowDBModel struct {
	CreatedAt    time.Time `db:"created_at"`
	OrderID      string    `db:"order_id"`
	EmployeeName string    `db:"employee_name"`
	TotalPrice   float64   `db:"total_price"`
}

type employeeSummaryRowDBModel struct {
	EmployeeName string  `db:"employee_name"`
	ShiftCount   int     `db:"shift_count"`
	OrdersCount  int     `db:"orders_count"`
	Revenue      float64 `db:"revenue"`
	AverageBill  float64 `db:"average_bill"`
}

type employeeOrderRowDBModel struct {
	CreatedAt    time.Time `db:"created_at"`
	EmployeeName string    `db:"employee_name"`
	OrderID      string    `db:"order_id"`
	TotalPrice   float64   `db:"total_price"`
}

type employeeShiftRowDBModel struct {
	EmployeeName string     `db:"employee_name"`
	OpenedAt     time.Time  `db:"opened_at"`
	ClosedAt     *time.Time `db:"closed_at"`
}

type inventorySummaryRowDBModel struct {
	IngredientName   string  `db:"ingredient_name"`
	Unit             string  `db:"unit"`
	Restocked        float64 `db:"restocked"`
	AutoWrittenOff   float64 `db:"auto_written_off"`
	ManualWrittenOff float64 `db:"manual_written_off"`
	CurrentStock     float64 `db:"current_stock"`
}

type reportInventoryOperationDBModel struct {
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

type menuPopularityRowDBModel struct {
	Place           int     `db:"place"`
	MenuItemName    string  `db:"menu_item_name"`
	SalesCount      int     `db:"sales_count"`
	Revenue         float64 `db:"revenue"`
	RevenueSharePct float64 `db:"revenue_share_pct"`
}

func toServiceSalesSummary(item salesSummaryDBModel) models.SalesSummary {
	return models.SalesSummary{
		OrdersCount: item.OrdersCount,
		Revenue:     item.Revenue,
		AverageBill: item.AverageBill,
	}
}

func toServiceSalesMenuRows(items []salesMenuItemRowDBModel) []models.SalesMenuItemRow {
	result := make([]models.SalesMenuItemRow, 0, len(items))
	for _, item := range items {
		result = append(result, models.SalesMenuItemRow{
			MenuItemName:  item.MenuItemName,
			Quantity:      item.Quantity,
			Revenue:       item.Revenue,
			SalesSharePct: item.SalesSharePct,
		})
	}
	return result
}

func toServiceSalesOrderRows(items []salesOrderRowDBModel) []models.SalesOrderRow {
	result := make([]models.SalesOrderRow, 0, len(items))
	for _, item := range items {
		result = append(result, models.SalesOrderRow{
			CreatedAt:    item.CreatedAt,
			OrderID:      item.OrderID,
			EmployeeName: item.EmployeeName,
			TotalPrice:   item.TotalPrice,
		})
	}
	return result
}

func toServiceEmployeeSummaryRows(items []employeeSummaryRowDBModel) []models.EmployeeSummaryRow {
	result := make([]models.EmployeeSummaryRow, 0, len(items))
	for _, item := range items {
		result = append(result, models.EmployeeSummaryRow{
			EmployeeName: item.EmployeeName,
			ShiftCount:   item.ShiftCount,
			OrdersCount:  item.OrdersCount,
			Revenue:      item.Revenue,
			AverageBill:  item.AverageBill,
		})
	}
	return result
}

func toServiceEmployeeOrderRows(items []employeeOrderRowDBModel) []models.EmployeeOrderRow {
	result := make([]models.EmployeeOrderRow, 0, len(items))
	for _, item := range items {
		result = append(result, models.EmployeeOrderRow{
			CreatedAt:    item.CreatedAt,
			EmployeeName: item.EmployeeName,
			OrderID:      item.OrderID,
			TotalPrice:   item.TotalPrice,
		})
	}
	return result
}

func toServiceEmployeeShiftRows(items []employeeShiftRowDBModel) []models.EmployeeShiftRow {
	result := make([]models.EmployeeShiftRow, 0, len(items))
	for _, item := range items {
		row := models.EmployeeShiftRow{
			EmployeeName: item.EmployeeName,
			OpenedAt:     item.OpenedAt,
			ClosedAt:     item.ClosedAt,
		}
		if item.ClosedAt != nil {
			row.Duration = item.ClosedAt.Sub(item.OpenedAt).String()
		} else {
			row.Duration = "active"
		}
		result = append(result, row)
	}
	return result
}

func toServiceInventorySummaryRows(items []inventorySummaryRowDBModel) []models.InventorySummaryRow {
	result := make([]models.InventorySummaryRow, 0, len(items))
	for _, item := range items {
		result = append(result, models.InventorySummaryRow{
			IngredientName:   item.IngredientName,
			Unit:             item.Unit,
			Restocked:        item.Restocked,
			AutoWrittenOff:   item.AutoWrittenOff,
			ManualWrittenOff: item.ManualWrittenOff,
			CurrentStock:     item.CurrentStock,
		})
	}
	return result
}

func toServiceReportInventoryOperations(items []reportInventoryOperationDBModel) []models.InventoryOperation {
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
				CurrentStock: item.IngredientCurrentStock,
			},
			Order: orderRef,
			Shift: models.InventoryOperationShift{ID: item.ShiftID},
			Employee: models.InventoryOperationEmployee{
				ID:       item.EmployeeID,
				FullName: item.EmployeeFullName,
			},
		})
	}
	return result
}

func toServiceMenuPopularityRows(items []menuPopularityRowDBModel) []models.MenuPopularityRow {
	result := make([]models.MenuPopularityRow, 0, len(items))
	for _, item := range items {
		result = append(result, models.MenuPopularityRow{
			Place:           item.Place,
			MenuItemName:    item.MenuItemName,
			SalesCount:      item.SalesCount,
			Revenue:         item.Revenue,
			RevenueSharePct: item.RevenueSharePct,
		})
	}
	return result
}
