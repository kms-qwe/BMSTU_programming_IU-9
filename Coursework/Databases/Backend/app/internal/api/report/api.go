package report

import (
	"context"

	"coffee-shop-backend/app/internal/api/common"
	"coffee-shop-backend/app/internal/models"

	"github.com/gin-gonic/gin"
)

type Service interface {
	BuildSalesReport(ctx context.Context, filter models.ReportFilter) ([]byte, error)
	BuildEmployeeReport(ctx context.Context, filter models.ReportFilter) ([]byte, error)
	BuildInventoryReport(ctx context.Context, filter models.InventoryReportFilter) ([]byte, error)
	BuildMenuPopularityReport(ctx context.Context, filter models.ReportFilter) ([]byte, error)
}

type API struct{ service Service }

func NewAPI(service Service) *API { return &API{service: service} }
func (a *API) BuildPath() string  { return "/api/report" }
func (a *API) RegisterHandlers(group *gin.RouterGroup) {
	group.POST("/sales/pdf", common.Wrap(a.createSalesReport))
	group.POST("/employees/pdf", common.Wrap(a.createEmployeesReport))
	group.POST("/inventory/pdf", common.Wrap(a.createInventoryReport))
	group.POST("/menu-popularity/pdf", common.Wrap(a.createMenuPopularityReport))
}
