package report

import (
	"context"

	"coffee-shop-backend/app/internal/models"
	"github.com/xuri/excelize/v2"
)

func (s *Service) BuildEmployeeReport(ctx context.Context, filter models.ReportFilter) ([]byte, error) {
	data, err := s.repo.GetEmployeeReport(ctx, filter)
	if err != nil {
		return nil, err
	}
	file := excelize.NewFile()
	file.SetSheetName("Sheet1", "Employees summary")
	rows := [][]any{{"Сотрудник", "Количество смен", "Количество заказов", "Выручка", "Средний чек"}}
	for _, row := range data.Summaries {
		rows = append(rows, []any{row.EmployeeName, row.ShiftCount, row.OrdersCount, row.Revenue, row.AverageBill})
	}
	writeRows(file, "Employees summary", rows)
	file.NewSheet("Employee orders")
	orderRows := [][]any{{"Дата", "Сотрудник", "Заказ", "Сумма"}}
	for _, row := range data.Orders {
		orderRows = append(orderRows, []any{row.CreatedAt, row.EmployeeName, row.OrderID, row.TotalPrice})
	}
	writeRows(file, "Employee orders", orderRows)
	file.NewSheet("Shifts")
	shiftRows := [][]any{{"Сотрудник", "Открытие смены", "Закрытие смены", "Длительность"}}
	for _, row := range data.Shifts {
		closedAt := any("")
		if row.ClosedAt != nil {
			closedAt = *row.ClosedAt
		}
		shiftRows = append(shiftRows, []any{row.EmployeeName, row.OpenedAt, closedAt, row.Duration})
	}
	writeRows(file, "Shifts", shiftRows)
	return writeWorkbook(file)
}
