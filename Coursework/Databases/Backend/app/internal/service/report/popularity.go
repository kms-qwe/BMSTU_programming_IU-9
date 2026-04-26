package report

import (
	"context"

	"coffee-shop-backend/app/internal/models"
	"github.com/xuri/excelize/v2"
)

func (s *Service) BuildMenuPopularityReport(ctx context.Context, filter models.ReportFilter) ([]byte, error) {
	data, err := s.repo.GetMenuPopularityReport(ctx, filter)
	if err != nil {
		return nil, err
	}
	file := excelize.NewFile()
	file.SetSheetName("Sheet1", "Popularity")
	rows := [][]any{{"Место", "Товар", "Количество продаж", "Выручка", "Доля"}}
	for _, row := range data.Rows {
		rows = append(rows, []any{row.Place, row.MenuItemName, row.SalesCount, row.Revenue, row.RevenueSharePct})
	}
	writeRows(file, "Popularity", rows)
	return writeWorkbook(file)
}
