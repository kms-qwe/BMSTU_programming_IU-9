package report

import (
	"fmt"

	"coffee-shop-backend/app/internal/api/common"
	"coffee-shop-backend/app/pkg/apimodels"

	"github.com/gin-gonic/gin"
)

// createMenuPopularityReport godoc
// @Summary Create menu popularity report
// @Description Builds an Excel menu popularity report for the selected period.
// @Tags report
// @Accept json
// @Produce application/vnd.openxmlformats-officedocument.spreadsheetml.sheet
// @Param request body apimodels.TimeRangeRequest true "Menu popularity report payload"
// @Success 200 {file} file
// @Failure 400 {object} apimodels.ErrorResponse
// @Failure 500 {object} apimodels.ErrorResponse
// @Router /api/report/menu-popularity/excel [post]
func (a *API) createMenuPopularityReport(ctx *gin.Context) error {
	var request apimodels.TimeRangeRequest
	if err := common.DecodeBody(ctx, &request); err != nil {
		return err
	}
	filter, err := toReportFilter(request)
	if err != nil {
		return err
	}
	content, err := a.service.BuildMenuPopularityReport(ctx.Request.Context(), filter)
	if err != nil {
		return err
	}
	return common.Excel(ctx, fmt.Sprintf("menu-popularity-report-%s-%s.xlsx", filter.From.Format("2006-01-02"), filter.To.Format("2006-01-02")), content)
}
