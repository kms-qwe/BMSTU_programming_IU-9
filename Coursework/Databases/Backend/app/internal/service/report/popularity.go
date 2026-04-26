package report

import (
	"context"

	"coffee-shop-backend/app/internal/models"
)

func (s *Service) BuildMenuPopularityReport(ctx context.Context, filter models.ReportFilter) ([]byte, error) {
	data, err := s.repo.GetMenuPopularityReport(ctx, filter)
	if err != nil {
		return nil, err
	}

	pdf := newPDFDocument("Отчет по популярности меню", true)
	addKeyValueSection(pdf, "Параметры отчета", []reportMetric{
		{Label: "Период", Value: filter.From.Format("2006-01-02 15:04 UTC") + " .. " + filter.To.Format("2006-01-02 15:04 UTC")},
		{Label: "Сотрудник", Value: employeeFilterLabel(filter, nil)},
	})

	totalUnits := 0
	totalRevenue := 0.0
	for _, row := range data.Rows {
		totalUnits += row.SalesCount
		totalRevenue += row.Revenue
	}

	addKeyValueSection(pdf, "Сводные метрики", []reportMetric{
		{Label: "Позиций в выборке", Value: formatInt(len(data.Rows))},
		{Label: "Всего продано штук", Value: formatInt(totalUnits)},
		{Label: "Общая выручка", Value: formatMoney(totalRevenue)},
		{Label: "Средняя выручка за единицу", Value: formatMoney(safeDivide(totalRevenue, float64(maxInt(totalUnits, 1))))},
		{Label: "Лидер по количеству", Value: topPopularityBySales(data.Rows)},
		{Label: "Лидер по выручке", Value: topPopularityByRevenue(data.Rows)},
		{Label: "Доля топ-3 по выручке", Value: formatMoney(topThreeShare(data.Rows)) + "%"},
	})

	rows := make([][]string, 0, len(data.Rows))
	for _, row := range data.Rows {
		rows = append(rows, []string{
			formatInt(row.Place),
			row.MenuItemName,
			formatInt(row.SalesCount),
			formatMoney(row.Revenue),
			formatMoney(row.RevenueSharePct) + "%",
		})
	}
	addTable(pdf, "Рейтинг популярности", []string{"Место", "Позиция", "Продано", "Выручка", "Доля выручки"}, []float64{20, 110, 30, 35, 35}, toStringRows(25, rows))

	return writePDF(pdf)
}
