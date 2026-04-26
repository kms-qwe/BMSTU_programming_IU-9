package report

import (
	"bytes"
	"context"

	"coffee-shop-backend/app/internal/models"
	"coffee-shop-backend/app/internal/utils"
	"github.com/xuri/excelize/v2"
)

type IReportRepo interface {
	GetSalesReport(ctx context.Context, filter models.ReportFilter) (*models.SalesReportData, error)
	GetEmployeeReport(ctx context.Context, filter models.ReportFilter) (*models.EmployeeReportData, error)
	GetInventoryReport(ctx context.Context, filter models.InventoryReportFilter) (*models.InventoryReportData, error)
	GetMenuPopularityReport(ctx context.Context, filter models.ReportFilter) (*models.MenuPopularityReportData, error)
}

type Service struct{ repo IReportRepo }

func New(repo IReportRepo) *Service { return &Service{repo: repo} }

func writeRows(file *excelize.File, sheet string, rows [][]any) {
	for rowIdx, row := range rows {
		cell, _ := excelize.CoordinatesToCellName(1, rowIdx+1)
		_ = file.SetSheetRow(sheet, cell, &row)
	}
	style, _ := file.NewStyle(&excelize.Style{Font: &excelize.Font{Bold: true}})
	_ = file.SetCellStyle(sheet, "A1", "F1", style)
}
func writeWorkbook(file *excelize.File) ([]byte, error) {
	buffer, err := file.WriteToBuffer()
	if err != nil {
		return nil, utils.InternalError("failed to build excel report", err)
	}
	return bytes.Clone(buffer.Bytes()), nil
}
