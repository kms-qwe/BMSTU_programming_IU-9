package shift

import (
	"strings"

	"coffee-shop-backend/app/internal/api/common"
	"coffee-shop-backend/app/internal/utils"
	"coffee-shop-backend/app/pkg/apimodels"

	"github.com/gin-gonic/gin"
)

// openShift godoc
// @Summary Open shift
// @Description Opens a new shift for an employee by login and password.
// @Tags shift
// @Accept json
// @Produce json
// @Param request body apimodels.OpenShiftRequest true "Open shift payload"
// @Success 200 {object} apimodels.OpenShiftResponse
// @Failure 400 {object} apimodels.ErrorResponse
// @Failure 401 {object} apimodels.ErrorResponse
// @Failure 409 {object} apimodels.ErrorResponse
// @Failure 500 {object} apimodels.ErrorResponse
// @Router /api/shift/open [post]
func (a *API) openShift(ctx *gin.Context) error {
	var request apimodels.OpenShiftRequest
	if err := common.DecodeBody(ctx, &request); err != nil {
		return err
	}
	if strings.TrimSpace(request.Login) == "" || strings.TrimSpace(request.Password) == "" {
		return utils.ValidationError("login and password are required", nil)
	}

	shift, err := a.service.OpenShift(ctx.Request.Context(), request.Login, request.Password)
	if err != nil {
		return err
	}

	ctx.JSON(200, apimodels.OpenShiftResponse{Shift: toShiftResponse(shift)})
	return nil
}
