package report

import (
	"fmt"

	"coffee-shop-backend/app/internal/api/common"
	"coffee-shop-backend/app/pkg/apimodels"

	"github.com/gin-gonic/gin"
)

// createInventoryReport godoc
// @Summary Create inventory report
// @Description Builds an Excel inventory report for the selected period.
// @Tags report
// @Accept json
// @Produce application/vnd.openxmlformats-officedocument.spreadsheetml.sheet
// @Param request body apimodels.InventoryReportRequest true "Inventory report payload"
// @Success 200 {file} file
// @Failure 400 {object} apimodels.ErrorResponse
// @Failure 500 {object} apimodels.ErrorResponse
// @Router /api/report/inventory/excel [post]
func (a *API) createInventoryReport(ctx *gin.Context) error {
	var request apimodels.InventoryReportRequest
	if err := common.DecodeBody(ctx, &request); err != nil {
		return err
	}
	filter, err := toInventoryReportFilter(request)
	if err != nil {
		return err
	}
	content, err := a.service.BuildInventoryReport(ctx.Request.Context(), filter)
	if err != nil {
		return err
	}
	return common.Excel(ctx, fmt.Sprintf("inventory-report-%s-%s.xlsx", filter.From.Format("2006-01-02"), filter.To.Format("2006-01-02")), content)
}
