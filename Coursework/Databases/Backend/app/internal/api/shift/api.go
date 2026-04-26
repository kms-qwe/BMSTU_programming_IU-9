package shift

import (
	"context"

	"coffee-shop-backend/app/internal/api/common"
	"coffee-shop-backend/app/internal/models"

	"github.com/gin-gonic/gin"
)

type Service interface {
	GetActiveShift(ctx context.Context) (*models.Shift, error)
	OpenShift(ctx context.Context, login, password string) (*models.Shift, error)
	CloseShift(ctx context.Context) (*models.Shift, error)
}

type API struct {
	service Service
}

func NewAPI(service Service) *API {
	return &API{service: service}
}

func (a *API) BuildPath() string {
	return "/api/shift"
}

func (a *API) RegisterHandlers(group *gin.RouterGroup) {
	group.GET("/active", common.Wrap(a.getShift))
	group.POST("/open", common.Wrap(a.openShift))
	group.POST("/close", common.Wrap(a.closeShift))
}
