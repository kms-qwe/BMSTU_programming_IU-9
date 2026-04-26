package menuitem

import (
	"context"

	"coffee-shop-backend/app/internal/api/common"
	"coffee-shop-backend/app/internal/models"

	"github.com/gin-gonic/gin"
)

type Service interface {
	ListMenuItems(ctx context.Context, includeInactive bool) ([]models.MenuItem, error)
	GetMenuItem(ctx context.Context, menuItemID string) (*models.MenuItem, error)
	CreateMenuItem(ctx context.Context, params models.CreateMenuItemParams) (*models.MenuItem, error)
	UpdateMenuItemPrice(ctx context.Context, params models.UpdateMenuItemPriceParams) (*models.MenuItem, error)
	SetMenuItemActivity(ctx context.Context, params models.SetMenuItemActivityParams) (*models.MenuItem, error)
}

type API struct {
	service Service
}

func NewAPI(service Service) *API {
	return &API{service: service}
}

func (a *API) BuildPath() string {
	return "/api/menu-item"
}

func (a *API) RegisterHandlers(group *gin.RouterGroup) {
	group.GET("", common.Wrap(a.listMenuItems))
	group.GET("/:menu_item_id", common.Wrap(a.getMenuItem))
	group.POST("", common.Wrap(a.createMenuItem))
	group.PATCH("/:menu_item_id/price", common.Wrap(a.updateMenuItemPrice))
	group.PATCH("/:menu_item_id/activity", common.Wrap(a.setMenuItemActivity))
}
