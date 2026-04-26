package employee

import (
	"context"

	"coffee-shop-backend/app/internal/api/common"
	"coffee-shop-backend/app/internal/models"

	"github.com/gin-gonic/gin"
)

type Service interface {
	ListEmployees(ctx context.Context) ([]models.Employee, error)
}

type API struct {
	service Service
}

func NewAPI(service Service) *API {
	return &API{service: service}
}

func (a *API) BuildPath() string {
	return "/api/employee"
}

func (a *API) RegisterHandlers(group *gin.RouterGroup) {
	group.GET("", common.Wrap(a.listEmployees))
}
