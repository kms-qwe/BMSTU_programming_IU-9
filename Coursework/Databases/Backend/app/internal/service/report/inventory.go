package report

import (
	"context"

	"coffee-shop-backend/app/internal/models"
)

func (s *Service) BuildInventoryReport(ctx context.Context, filter models.InventoryReportFilter) ([]byte, error) {
	data, err := s.repo.GetInventoryReport(ctx, filter)
	if err != nil {
		return nil, err
	}

	pdf := newPDFDocument("Отчет по складу", true)
	employeeNames := uniqueStrings(extractInventoryEmployeeNames(data.Operations))
	ingredientNames := uniqueStrings(extractInventoryIngredientNames(data))
	addKeyValueSection(pdf, "Параметры отчета", []reportMetric{
		{Label: "Период", Value: filter.From.Format("2006-01-02 15:04 UTC") + " .. " + filter.To.Format("2006-01-02 15:04 UTC")},
		{Label: "Сотрудник", Value: inventoryEmployeeFilterLabel(filter, employeeNames)},
		{Label: "Ингредиент", Value: ingredientFilterLabel(filter, ingredientNames)},
		{Label: "Тип операции", Value: operationTypeFilterLabel(filter)},
	})

	totalRestocked := 0.0
	totalAuto := 0.0
	totalManual := 0.0
	for _, row := range data.Summaries {
		totalRestocked += row.Restocked
		totalAuto += row.AutoWrittenOff
		totalManual += row.ManualWrittenOff
	}

	addKeyValueSection(pdf, "Сводные метрики", []reportMetric{
		{Label: "Ингредиентов в выборке", Value: formatInt(len(data.Summaries))},
		{Label: "Всего пополнено", Value: formatQuantity(totalRestocked)},
		{Label: "Автосписание", Value: formatQuantity(totalAuto)},
		{Label: "Ручное списание", Value: formatQuantity(totalManual)},
		{Label: "Чистое движение", Value: formatQuantity(totalRestocked - totalAuto - totalManual)},
		{Label: "Операций в выборке", Value: formatInt(len(data.Operations))},
		{Label: "Низкий остаток (<= 5)", Value: formatInt(countLowStockItems(data.Summaries))},
		{Label: "Лидер автосписания", Value: topInventoryMetric(data.Summaries, func(row models.InventorySummaryRow) float64 { return row.AutoWrittenOff })},
		{Label: "Лидер ручного списания", Value: topInventoryMetric(data.Summaries, func(row models.InventorySummaryRow) float64 { return row.ManualWrittenOff })},
		{Label: "Лидер пополнения", Value: topInventoryMetric(data.Summaries, func(row models.InventorySummaryRow) float64 { return row.Restocked })},
	})

	summaryRows := make([][]string, 0, len(data.Summaries))
	for _, row := range data.Summaries {
		summaryRows = append(summaryRows, []string{
			row.IngredientName,
			row.Unit,
			formatQuantity(row.Restocked),
			formatQuantity(row.AutoWrittenOff),
			formatQuantity(row.ManualWrittenOff),
			formatQuantity(row.CurrentStock),
		})
	}
	addTable(pdf, "Сводка по ингредиентам", []string{"Ингредиент", "Ед.", "Пополнено", "Автосписание", "Ручное списание", "Текущий остаток"}, []float64{70, 20, 35, 35, 35, 35}, toStringRows(20, summaryRows))

	return writePDF(pdf)
}

func extractInventoryEmployeeNames(rows []models.InventoryOperation) []string {
	names := make([]string, 0, len(rows))
	for _, row := range rows {
		names = append(names, row.Employee.FullName)
	}
	return names
}

func extractInventoryIngredientNames(data *models.InventoryReportData) []string {
	names := make([]string, 0, len(data.Summaries)+len(data.Operations))
	for _, row := range data.Summaries {
		names = append(names, row.IngredientName)
	}
	for _, row := range data.Operations {
		names = append(names, row.Ingredient.Name)
	}
	return names
}
