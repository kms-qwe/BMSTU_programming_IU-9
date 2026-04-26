package order

import (
	"context"

	"coffee-shop-backend/app/internal/api/common"
	"coffee-shop-backend/app/internal/models"

	"github.com/gin-gonic/gin"
)

type Service interface {
	ListOrders(ctx context.Context, filter models.OrderFilter) (*models.PaginatedResult[models.Order], error)
	GetOrder(ctx context.Context, orderID string) (*models.Order, error)
	CreateOrder(ctx context.Context, params models.CreateOrderParams) (*models.Order, error)
}

type API struct{ service Service }

func NewAPI(service Service) *API { return &API{service: service} }
func (a *API) BuildPath() string  { return "/api/order" }
func (a *API) RegisterHandlers(group *gin.RouterGroup) {
	group.GET("", common.Wrap(a.listOrders))
	group.GET("/:order_id", common.Wrap(a.getOrder))
	group.POST("", common.Wrap(a.createOrder))
}
