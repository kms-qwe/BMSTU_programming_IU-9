package report

import (
	"fmt"

	"coffee-shop-backend/app/internal/api/common"
	"coffee-shop-backend/app/pkg/apimodels"

	"github.com/gin-gonic/gin"
)

// createEmployeesReport godoc
// @Summary Create employees report
// @Description Builds an Excel employees report for the selected period.
// @Tags report
// @Accept json
// @Produce application/vnd.openxmlformats-officedocument.spreadsheetml.sheet
// @Param request body apimodels.TimeRangeRequest true "Employees report payload"
// @Success 200 {file} file
// @Failure 400 {object} apimodels.ErrorResponse
// @Failure 500 {object} apimodels.ErrorResponse
// @Router /api/report/employees/excel [post]
func (a *API) createEmployeesReport(ctx *gin.Context) error {
	var request apimodels.TimeRangeRequest
	if err := common.DecodeBody(ctx, &request); err != nil {
		return err
	}
	filter, err := toReportFilter(request)
	if err != nil {
		return err
	}
	content, err := a.service.BuildEmployeeReport(ctx.Request.Context(), filter)
	if err != nil {
		return err
	}
	return common.Excel(ctx, fmt.Sprintf("employees-report-%s-%s.xlsx", filter.From.Format("2006-01-02"), filter.To.Format("2006-01-02")), content)
}
