package report

import (
	"strings"
	"time"

	"coffee-shop-backend/app/internal/models"
	"coffee-shop-backend/app/internal/utils"
	"coffee-shop-backend/app/pkg/apimodels"
)

func toReportFilter(request apimodels.TimeRangeRequest) (models.ReportFilter, error) {
	from, to, err := parseRange(request.From, request.To)
	if err != nil {
		return models.ReportFilter{}, err
	}
	return models.ReportFilter{From: from, To: to, EmployeeID: normalizeOptional(request.EmployeeID)}, nil
}

func toInventoryReportFilter(request apimodels.InventoryReportRequest) (models.InventoryReportFilter, error) {
	from, to, err := parseRange(request.From, request.To)
	if err != nil {
		return models.InventoryReportFilter{}, err
	}
	return models.InventoryReportFilter{From: from, To: to, EmployeeID: normalizeOptional(request.EmployeeID), IngredientID: normalizeOptional(request.IngredientID), OperationType: normalizeOptional(request.OperationType)}, nil
}

func parseRange(fromRaw, toRaw string) (time.Time, time.Time, error) {
	from, err := time.Parse(time.RFC3339, strings.TrimSpace(fromRaw))
	if err != nil {
		return time.Time{}, time.Time{}, utils.ValidationError("invalid from date", nil)
	}
	to, err := time.Parse(time.RFC3339, strings.TrimSpace(toRaw))
	if err != nil {
		return time.Time{}, time.Time{}, utils.ValidationError("invalid to date", nil)
	}
	if to.Before(from) {
		return time.Time{}, time.Time{}, utils.ValidationError("to must be greater than or equal to from", nil)
	}
	return from, to, nil
}

func normalizeOptional(value *string) *string {
	if value == nil {
		return nil
	}
	trimmed := strings.TrimSpace(*value)
	if trimmed == "" {
		return nil
	}
	return &trimmed
}
