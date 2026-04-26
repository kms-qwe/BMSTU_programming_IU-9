package ingredient

import (
	"context"

	"coffee-shop-backend/app/internal/api/common"
	"coffee-shop-backend/app/internal/models"

	"github.com/gin-gonic/gin"
)

type Service interface {
	ListIngredients(ctx context.Context, includeInactive bool) ([]models.Ingredient, error)
	GetIngredient(ctx context.Context, ingredientID string) (*models.Ingredient, error)
	CreateIngredient(ctx context.Context, params models.CreateIngredientParams) (*models.Ingredient, error)
	SetIngredientActivity(ctx context.Context, params models.SetIngredientActivityParams) (*models.Ingredient, error)
}

type API struct {
	service Service
}

func NewAPI(service Service) *API {
	return &API{service: service}
}

func (a *API) BuildPath() string {
	return "/api/ingredient"
}

func (a *API) RegisterHandlers(group *gin.RouterGroup) {
	group.GET("", common.Wrap(a.listIngredients))
	group.GET("/:ingredient_id", common.Wrap(a.getIngredient))
	group.POST("", common.Wrap(a.createIngredient))
	group.PATCH("/:ingredient_id/activity", common.Wrap(a.setIngredientActivity))
}
