package report

import (
	"bytes"
	"context"
	_ "embed"
	"fmt"
	"math"
	"sort"
	"strings"
	"time"

	"coffee-shop-backend/app/internal/models"
	"coffee-shop-backend/app/internal/utils"

	"github.com/go-pdf/fpdf"
)

type IReportRepo interface {
	GetSalesReport(ctx context.Context, filter models.ReportFilter) (*models.SalesReportData, error)
	GetEmployeeReport(ctx context.Context, filter models.ReportFilter) (*models.EmployeeReportData, error)
	GetInventoryReport(ctx context.Context, filter models.InventoryReportFilter) (*models.InventoryReportData, error)
	GetMenuPopularityReport(ctx context.Context, filter models.ReportFilter) (*models.MenuPopularityReportData, error)
}

type Service struct{ repo IReportRepo }

func New(repo IReportRepo) *Service { return &Service{repo: repo} }

//go:embed assets/ArialUnicode.ttf
var reportFontRegular []byte

const reportFontFamily = "ArialUnicode"

type reportMetric struct {
	Label string
	Value string
}

func newPDFDocument(title string, landscape bool) *fpdf.Fpdf {
	orientation := "P"
	if landscape {
		orientation = "L"
	}

	pdf := fpdf.New(orientation, "mm", "A4", "")
	pdf.SetMargins(10, 10, 10)
	pdf.SetAutoPageBreak(true, 10)
	pdf.AliasNbPages("")
	pdf.SetTitle(title, false)
	pdf.SetCreator("coffee-shop-backend", false)
	pdf.SetAuthor("coffee-shop-backend", false)
	registerReportFonts(pdf)
	pdf.SetFont(reportFontFamily, "B", 16)
	pdf.AddPage()
	pdf.CellFormat(0, 10, title, "", 1, "L", false, 0, "")
	pdf.SetFont(reportFontFamily, "", 9)
	pdf.SetTextColor(90, 90, 90)
	pdf.CellFormat(0, 6, "Сформирован: "+time.Now().UTC().Format(time.RFC3339), "", 1, "L", false, 0, "")
	pdf.SetTextColor(0, 0, 0)
	pdf.Ln(2)
	return pdf
}

func registerReportFonts(pdf *fpdf.Fpdf) {
	pdf.AddUTF8FontFromBytes(reportFontFamily, "", reportFontRegular)
	pdf.AddUTF8FontFromBytes(reportFontFamily, "B", reportFontRegular)
}

func writePDF(pdf *fpdf.Fpdf) ([]byte, error) {
	var buffer bytes.Buffer
	if err := pdf.Output(&buffer); err != nil {
		return nil, utils.InternalError("failed to build pdf report", err)
	}
	return bytes.Clone(buffer.Bytes()), nil
}

func addSectionTitle(pdf *fpdf.Fpdf, title string) {
	pdf.SetFont(reportFontFamily, "B", 12)
	pdf.CellFormat(0, 8, title, "", 1, "L", false, 0, "")
	pdf.SetFont(reportFontFamily, "", 9)
}

func addKeyValueSection(pdf *fpdf.Fpdf, title string, rows []reportMetric) {
	addSectionTitle(pdf, title)
	for _, row := range rows {
		pdf.SetFont(reportFontFamily, "B", 9)
		pdf.CellFormat(55, 6, row.Label, "0", 0, "L", false, 0, "")
		pdf.SetFont(reportFontFamily, "", 9)
		pdf.MultiCell(0, 6, row.Value, "0", "L", false)
	}
	pdf.Ln(1)
}

func addTable(pdf *fpdf.Fpdf, title string, headers []string, widths []float64, rows [][]string) {
	addSectionTitle(pdf, title)
	renderHeader := func() {
		pdf.SetFont(reportFontFamily, "B", 8)
		for index, header := range headers {
			pdf.CellFormat(widths[index], 7, header, "1", 0, "L", false, 0, "")
		}
		pdf.Ln(-1)
		pdf.SetFont(reportFontFamily, "", 8)
	}

	renderHeader()
	for _, row := range rows {
		if pdf.GetY() > 185 {
			pdf.AddPage()
			renderHeader()
		}
		for index, value := range row {
			pdf.CellFormat(widths[index], 6, truncateText(value, widths[index]), "1", 0, "L", false, 0, "")
		}
		pdf.Ln(-1)
	}
	pdf.Ln(2)
}

func truncateText(text string, width float64) string {
	limit := int(math.Max(8, width/2.4))
	runes := []rune(text)
	if len(runes) <= limit {
		return text
	}
	if limit <= 3 {
		return string(runes[:limit])
	}
	return string(runes[:limit-3]) + "..."
}

func formatMoney(value float64) string {
	return fmt.Sprintf("%.2f", value)
}

func formatInt(value int) string {
	return fmt.Sprintf("%d", value)
}

func formatQuantity(value float64) string {
	return fmt.Sprintf("%.3f", value)
}

func formatDateTime(value time.Time) string {
	return value.UTC().Format("2006-01-02 15:04")
}

func formatOptionalDateTime(value *time.Time) string {
	if value == nil {
		return "active"
	}
	return formatDateTime(*value)
}

func employeeFilterLabel(filter models.ReportFilter, names []string) string {
	if filter.EmployeeID == nil {
		return "Все сотрудники"
	}
	if len(names) > 0 {
		return names[0] + " (конкретный сотрудник)"
	}
	return "Конкретный сотрудник, ID: " + *filter.EmployeeID
}

func inventoryEmployeeFilterLabel(filter models.InventoryReportFilter, names []string) string {
	if filter.EmployeeID == nil {
		return "Все сотрудники"
	}
	if len(names) > 0 {
		return names[0] + " (конкретный сотрудник)"
	}
	return "Конкретный сотрудник, ID: " + *filter.EmployeeID
}

func ingredientFilterLabel(filter models.InventoryReportFilter, names []string) string {
	if filter.IngredientID == nil {
		return "Все ингредиенты"
	}
	if len(names) > 0 {
		return names[0] + " (конкретный ингредиент)"
	}
	return "Конкретный ингредиент, ID: " + *filter.IngredientID
}

func operationTypeFilterLabel(filter models.InventoryReportFilter) string {
	if filter.OperationType == nil {
		return "Все типы операций"
	}
	return *filter.OperationType
}

func uniqueStrings(values []string) []string {
	seen := make(map[string]struct{}, len(values))
	result := make([]string, 0, len(values))
	for _, value := range values {
		if value == "" {
			continue
		}
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		result = append(result, value)
	}
	sort.Strings(result)
	return result
}

func topRevenueEmployee(rows []models.SalesOrderRow) (string, float64) {
	totals := make(map[string]float64)
	for _, row := range rows {
		totals[row.EmployeeName] += row.TotalPrice
	}
	var bestName string
	var bestValue float64
	for name, value := range totals {
		if value > bestValue || bestName == "" {
			bestName = name
			bestValue = value
		}
	}
	if bestName == "" {
		return "No data", 0
	}
	return bestName, bestValue
}

func peakHour(rows []models.SalesOrderRow) string {
	counts := map[int]int{}
	for _, row := range rows {
		counts[row.CreatedAt.UTC().Hour()]++
	}
	bestHour := -1
	bestCount := -1
	for hour, count := range counts {
		if count > bestCount || (count == bestCount && hour < bestHour) {
			bestHour = hour
			bestCount = count
		}
	}
	if bestHour < 0 {
		return "Нет данных"
	}
	return fmt.Sprintf("%02d:00-%02d:59 UTC (%d заказов)", bestHour, bestHour, bestCount)
}

func orderBounds(rows []models.SalesOrderRow) (float64, float64) {
	if len(rows) == 0 {
		return 0, 0
	}
	minValue := rows[0].TotalPrice
	maxValue := rows[0].TotalPrice
	for _, row := range rows[1:] {
		minValue = math.Min(minValue, row.TotalPrice)
		maxValue = math.Max(maxValue, row.TotalPrice)
	}
	return minValue, maxValue
}

func longestShift(rows []models.EmployeeShiftRow) string {
	var bestName string
	var bestDuration time.Duration
	for _, row := range rows {
		if row.ClosedAt == nil {
			continue
		}
		duration := row.ClosedAt.Sub(row.OpenedAt)
		if duration > bestDuration || bestName == "" {
			bestName = row.EmployeeName
			bestDuration = duration
		}
	}
	if bestName == "" {
		return "Нет завершенных смен"
	}
	return fmt.Sprintf("%s (%s)", bestName, bestDuration.Round(time.Minute))
}

func topEmployeeByRevenue(rows []models.EmployeeSummaryRow) string {
	var best models.EmployeeSummaryRow
	found := false
	for _, row := range rows {
		if !found || row.Revenue > best.Revenue {
			best = row
			found = true
		}
	}
	if !found {
		return "Нет данных"
	}
	return fmt.Sprintf("%s (выручка %s)", best.EmployeeName, formatMoney(best.Revenue))
}

func topEmployeeByOrders(rows []models.EmployeeSummaryRow) string {
	var best models.EmployeeSummaryRow
	found := false
	for _, row := range rows {
		if !found || row.OrdersCount > best.OrdersCount {
			best = row
			found = true
		}
	}
	if !found {
		return "Нет данных"
	}
	return fmt.Sprintf("%s (%d заказов)", best.EmployeeName, best.OrdersCount)
}

func countLowStockItems(rows []models.InventorySummaryRow) int {
	count := 0
	for _, row := range rows {
		if row.CurrentStock <= 5 {
			count++
		}
	}
	return count
}

func topInventoryMetric(rows []models.InventorySummaryRow, selector func(models.InventorySummaryRow) float64) string {
	var best models.InventorySummaryRow
	var bestValue float64
	found := false
	for _, row := range rows {
		value := selector(row)
		if !found || value > bestValue {
			best = row
			bestValue = value
			found = true
		}
	}
	if !found {
		return "Нет данных"
	}
	return fmt.Sprintf("%s (%s %s)", best.IngredientName, formatQuantity(bestValue), best.Unit)
}

func topPopularityBySales(rows []models.MenuPopularityRow) string {
	if len(rows) == 0 {
		return "Нет данных"
	}
	best := rows[0]
	for _, row := range rows[1:] {
		if row.SalesCount > best.SalesCount {
			best = row
		}
	}
	return fmt.Sprintf("%s (%d шт.)", best.MenuItemName, best.SalesCount)
}

func topPopularityByRevenue(rows []models.MenuPopularityRow) string {
	if len(rows) == 0 {
		return "Нет данных"
	}
	best := rows[0]
	for _, row := range rows[1:] {
		if row.Revenue > best.Revenue {
			best = row
		}
	}
	return fmt.Sprintf("%s (%s)", best.MenuItemName, formatMoney(best.Revenue))
}

func topThreeShare(rows []models.MenuPopularityRow) float64 {
	limit := 3
	if len(rows) < limit {
		limit = len(rows)
	}
	total := 0.0
	for i := 0; i < limit; i++ {
		total += rows[i].RevenueSharePct
	}
	return total
}

func toStringRows(limit int, rows [][]string) [][]string {
	if len(rows) <= limit {
		return rows
	}
	return rows[:limit]
}

func joinNames(values []string) string {
	if len(values) == 0 {
		return "Нет данных"
	}
	return strings.Join(values, ", ")
}

func safeDivide(numerator, denominator float64) float64 {
	if denominator == 0 {
		return 0
	}
	return numerator / denominator
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
