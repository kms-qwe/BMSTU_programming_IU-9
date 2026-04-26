package report

import (
	"context"

	"coffee-shop-backend/app/internal/models"
	"github.com/xuri/excelize/v2"
)

func (s *Service) BuildInventoryReport(ctx context.Context, filter models.InventoryReportFilter) ([]byte, error) {
	data, err := s.repo.GetInventoryReport(ctx, filter)
	if err != nil {
		return nil, err
	}
	file := excelize.NewFile()
	file.SetSheetName("Sheet1", "Inventory summary")
	rows := [][]any{{"Ингредиент", "Ед. изм.", "Пополнено", "Списано по заказам", "Списано вручную", "Текущий остаток"}}
	for _, row := range data.Summaries {
		rows = append(rows, []any{row.IngredientName, row.Unit, row.Restocked, row.AutoWrittenOff, row.ManualWrittenOff, row.CurrentStock})
	}
	writeRows(file, "Inventory summary", rows)
	file.NewSheet("Operations")
	opRows := [][]any{{"Дата", "Ингредиент", "Тип операции", "Изменение", "Сотрудник", "Заказ"}}
	for _, row := range data.Operations {
		orderID := ""
		if row.Order != nil {
			orderID = row.Order.ID
		}
		opRows = append(opRows, []any{row.CreatedAt, row.Ingredient.Name, row.OperationType, row.ChangeAmount, row.Employee.FullName, orderID})
	}
	writeRows(file, "Operations", opRows)
	return writeWorkbook(file)
}
