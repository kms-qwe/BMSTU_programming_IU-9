package report

import (
	"context"

	"coffee-shop-backend/app/internal/models"
)

func (s *Service) BuildSalesReport(ctx context.Context, filter models.ReportFilter) ([]byte, error) {
	data, err := s.repo.GetSalesReport(ctx, filter)
	if err != nil {
		return nil, err
	}

	pdf := newPDFDocument("Отчет по продажам", true)
	employeeNames := uniqueStrings(extractSalesEmployeeNames(data.Orders))
	addKeyValueSection(pdf, "Параметры отчета", []reportMetric{
		{Label: "Период", Value: filter.From.Format("2006-01-02 15:04 UTC") + " .. " + filter.To.Format("2006-01-02 15:04 UTC")},
		{Label: "Сотрудник", Value: employeeFilterLabel(filter, employeeNames)},
	})

	minOrder, maxOrder := orderBounds(data.Orders)
	bestEmployeeName, bestEmployeeRevenue := topRevenueEmployee(data.Orders)
	topMenuByRevenue := "No data"
	topMenuByQuantity := "No data"
	if len(data.MenuItems) > 0 {
		topMenuByRevenue = data.MenuItems[0].MenuItemName + " (" + formatMoney(data.MenuItems[0].Revenue) + ")"
		bestByQuantity := data.MenuItems[0]
		for _, item := range data.MenuItems[1:] {
			if item.Quantity > bestByQuantity.Quantity {
				bestByQuantity = item
			}
		}
		topMenuByQuantity = bestByQuantity.MenuItemName + " (" + formatInt(bestByQuantity.Quantity) + " шт.)"
	} else {
		topMenuByRevenue = "Нет данных"
		topMenuByQuantity = "Нет данных"
	}

	addKeyValueSection(pdf, "Сводные метрики", []reportMetric{
		{Label: "Количество заказов", Value: formatInt(data.Summary.OrdersCount)},
		{Label: "Выручка", Value: formatMoney(data.Summary.Revenue)},
		{Label: "Средний чек", Value: formatMoney(data.Summary.AverageBill)},
		{Label: "Уникальных позиций продано", Value: formatInt(len(data.MenuItems))},
		{Label: "Самый большой заказ", Value: formatMoney(maxOrder)},
		{Label: "Самый маленький заказ", Value: formatMoney(minOrder)},
		{Label: "Пиковый час продаж", Value: peakHour(data.Orders)},
		{Label: "Лучший сотрудник по выручке", Value: bestEmployeeName + " (" + formatMoney(bestEmployeeRevenue) + ")"},
		{Label: "Лидер по выручке среди позиций", Value: topMenuByRevenue},
		{Label: "Лидер по количеству продаж", Value: topMenuByQuantity},
	})

	menuRows := make([][]string, 0, len(data.MenuItems))
	for _, row := range data.MenuItems {
		menuRows = append(menuRows, []string{
			row.MenuItemName,
			formatInt(row.Quantity),
			formatMoney(row.Revenue),
			formatMoney(row.SalesSharePct) + "%",
		})
	}
	addTable(pdf, "Разбивка по позициям меню", []string{"Позиция", "Шт.", "Выручка", "Доля выручки"}, []float64{95, 35, 45, 40}, toStringRows(20, menuRows))

	return writePDF(pdf)
}

func extractSalesEmployeeNames(rows []models.SalesOrderRow) []string {
	names := make([]string, 0, len(rows))
	for _, row := range rows {
		names = append(names, row.EmployeeName)
	}
	return names
}
