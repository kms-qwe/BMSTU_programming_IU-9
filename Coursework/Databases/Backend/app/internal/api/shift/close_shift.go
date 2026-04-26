package shift

import (
	"coffee-shop-backend/app/pkg/apimodels"

	"github.com/gin-gonic/gin"
)

// closeShift godoc
// @Summary Close active shift
// @Description Closes the current active shift.
// @Tags shift
// @Produce json
// @Success 200 {object} apimodels.CloseShiftResponse
// @Failure 400 {object} apimodels.ErrorResponse
// @Failure 500 {object} apimodels.ErrorResponse
// @Router /api/shift/close [post]
func (a *API) closeShift(ctx *gin.Context) error {
	shift, err := a.service.CloseShift(ctx.Request.Context())
	if err != nil {
		return err
	}

	ctx.JSON(200, apimodels.CloseShiftResponse{Shift: toShiftResponse(shift)})
	return nil
}
