package report

import (
	"context"

	"coffee-shop-backend/app/internal/models"
)

func (s *Service) BuildEmployeeReport(ctx context.Context, filter models.ReportFilter) ([]byte, error) {
	data, err := s.repo.GetEmployeeReport(ctx, filter)
	if err != nil {
		return nil, err
	}

	pdf := newPDFDocument("Отчет по сотрудникам", true)
	names := uniqueStrings(extractEmployeeNames(data))
	addKeyValueSection(pdf, "Параметры отчета", []reportMetric{
		{Label: "Период", Value: filter.From.Format("2006-01-02 15:04 UTC") + " .. " + filter.To.Format("2006-01-02 15:04 UTC")},
		{Label: "Сотрудник", Value: employeeFilterLabel(filter, names)},
	})

	totalShifts := 0
	totalOrders := 0
	totalRevenue := 0.0
	for _, row := range data.Summaries {
		totalShifts += row.ShiftCount
		totalOrders += row.OrdersCount
		totalRevenue += row.Revenue
	}

	addKeyValueSection(pdf, "Сводные метрики", []reportMetric{
		{Label: "Сотрудников в выборке", Value: formatInt(len(data.Summaries))},
		{Label: "Всего смен", Value: formatInt(totalShifts)},
		{Label: "Всего заказов", Value: formatInt(totalOrders)},
		{Label: "Общая выручка", Value: formatMoney(totalRevenue)},
		{Label: "Средний чек по выборке", Value: formatMoney(safeDivide(totalRevenue, float64(maxInt(totalOrders, 1))))},
		{Label: "Лидер по выручке", Value: topEmployeeByRevenue(data.Summaries)},
		{Label: "Лидер по числу заказов", Value: topEmployeeByOrders(data.Summaries)},
		{Label: "Самая длинная смена", Value: longestShift(data.Shifts)},
	})

	summaryRows := make([][]string, 0, len(data.Summaries))
	for _, row := range data.Summaries {
		summaryRows = append(summaryRows, []string{
			row.EmployeeName,
			formatInt(row.ShiftCount),
			formatInt(row.OrdersCount),
			formatMoney(row.Revenue),
			formatMoney(row.AverageBill),
		})
	}
	addTable(pdf, "Сводка по сотрудникам", []string{"Сотрудник", "Смены", "Заказы", "Выручка", "Средний чек"}, []float64{85, 25, 25, 35, 35}, toStringRows(20, summaryRows))

	return writePDF(pdf)
}

func extractEmployeeNames(data *models.EmployeeReportData) []string {
	names := make([]string, 0, len(data.Summaries)+len(data.Orders)+len(data.Shifts))
	for _, row := range data.Summaries {
		names = append(names, row.EmployeeName)
	}
	for _, row := range data.Orders {
		names = append(names, row.EmployeeName)
	}
	for _, row := range data.Shifts {
		names = append(names, row.EmployeeName)
	}
	return names
}
