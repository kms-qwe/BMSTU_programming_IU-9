package report

import (
	"fmt"

	"coffee-shop-backend/app/internal/api/common"
	"coffee-shop-backend/app/pkg/apimodels"

	"github.com/gin-gonic/gin"
)

// createSalesReport godoc
// @Summary Create sales report
// @Description Builds an Excel sales report for the selected period.
// @Tags report
// @Accept json
// @Produce application/vnd.openxmlformats-officedocument.spreadsheetml.sheet
// @Param request body apimodels.TimeRangeRequest true "Sales report payload"
// @Success 200 {file} file
// @Failure 400 {object} apimodels.ErrorResponse
// @Failure 500 {object} apimodels.ErrorResponse
// @Router /api/report/sales/excel [post]
func (a *API) createSalesReport(ctx *gin.Context) error {
	var request apimodels.TimeRangeRequest
	if err := common.DecodeBody(ctx, &request); err != nil {
		return err
	}
	filter, err := toReportFilter(request)
	if err != nil {
		return err
	}
	content, err := a.service.BuildSalesReport(ctx.Request.Context(), filter)
	if err != nil {
		return err
	}
	return common.Excel(ctx, fmt.Sprintf("sales-report-%s-%s.xlsx", filter.From.Format("2006-01-02"), filter.To.Format("2006-01-02")), content)
}
