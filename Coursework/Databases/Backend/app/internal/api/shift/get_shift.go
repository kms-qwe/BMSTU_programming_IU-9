package shift

import (
	"coffee-shop-backend/app/pkg/apimodels"

	"github.com/gin-gonic/gin"
)

// getShift godoc
// @Summary Get active shift
// @Description Returns the current active shift or null when there is no active shift.
// @Tags shift
// @Produce json
// @Success 200 {object} apimodels.ActiveShiftResponse
// @Failure 500 {object} apimodels.ErrorResponse
// @Router /api/shift/active [get]
func (a *API) getShift(ctx *gin.Context) error {
	shift, err := a.service.GetActiveShift(ctx.Request.Context())
	if err != nil {
		return err
	}

	ctx.JSON(200, apimodels.ActiveShiftResponse{
		IsActive: shift != nil,
		Shift:    toShiftResponse(shift),
	})
	return nil
}
