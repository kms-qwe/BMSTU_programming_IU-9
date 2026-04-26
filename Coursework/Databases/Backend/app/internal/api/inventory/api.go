package inventory

import (
	"context"

	"coffee-shop-backend/app/internal/api/common"
	"coffee-shop-backend/app/internal/models"

	"github.com/gin-gonic/gin"
)

type Service interface {
	ListInventoryOperations(ctx context.Context, filter models.InventoryOperationFilter) (*models.PaginatedResult[models.InventoryOperation], error)
	CreateManualInventoryOperation(ctx context.Context, params models.CreateInventoryOperationParams) (*models.InventoryOperation, error)
}

type API struct{ service Service }

func NewAPI(service Service) *API { return &API{service: service} }
func (a *API) BuildPath() string  { return "/api/inventory" }
func (a *API) RegisterHandlers(group *gin.RouterGroup) {
	group.GET("/operation", common.Wrap(a.listInventoryOperations))
	group.POST("/operation", common.Wrap(a.createInventoryOperation))
}
