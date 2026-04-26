package report

import (
	"context"

	"coffee-shop-backend/app/internal/models"
	"github.com/xuri/excelize/v2"
)

func (s *Service) BuildSalesReport(ctx context.Context, filter models.ReportFilter) ([]byte, error) {
	data, err := s.repo.GetSalesReport(ctx, filter)
	if err != nil {
		return nil, err
	}
	file := excelize.NewFile()
	file.SetSheetName("Sheet1", "Summary")
	writeRows(file, "Summary", [][]any{{"Показатель", "Значение"}, {"Период", filter.From.Format("02.01.2006") + "-" + filter.To.Format("02.01.2006")}, {"Количество заказов", data.Summary.OrdersCount}, {"Выручка", data.Summary.Revenue}, {"Средний чек", data.Summary.AverageBill}})
	file.NewSheet("Menu items")
	rows := [][]any{{"Товар", "Количество", "Выручка", "Доля от продаж"}}
	for _, row := range data.MenuItems {
		rows = append(rows, []any{row.MenuItemName, row.Quantity, row.Revenue, row.SalesSharePct})
	}
	writeRows(file, "Menu items", rows)
	file.NewSheet("Orders")
	orderRows := [][]any{{"Дата", "Заказ", "Сотрудник", "Сумма"}}
	for _, row := range data.Orders {
		orderRows = append(orderRows, []any{row.CreatedAt, row.OrderID, row.EmployeeName, row.TotalPrice})
	}
	writeRows(file, "Orders", orderRows)
	return writeWorkbook(file)
}
