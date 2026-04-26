package common

import (
	"errors"
	"fmt"
	"net/http"

	"coffee-shop-backend/app/internal/utils"
	"coffee-shop-backend/app/pkg/apimodels"

	"github.com/gin-gonic/gin"
)

type HandlerFunc func(ctx *gin.Context) error

func Wrap(handler HandlerFunc) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if err := handler(ctx); err != nil {
			WriteError(ctx, err)
		}
	}
}

func WriteError(ctx *gin.Context, err error) {
	var appErr *utils.AppError
	if !errors.As(err, &appErr) {
		utils.Logger().Error("unexpected handler error", "error", err, "path", ctx.Request.URL.Path)
		appErr = utils.InternalError("internal server error", err)
	}

	ctx.AbortWithStatusJSON(appErr.Status, apimodels.ErrorResponse{
		Error: apimodels.ErrorBody{
			Code:    string(appErr.Code),
			Message: appErr.Message,
			Details: appErr.Details,
		},
	})
}

func DecodeBody(ctx *gin.Context, target any) error {
	if err := ctx.ShouldBindJSON(target); err != nil {
		return utils.ValidationError("invalid request body", map[string]any{"reason": err.Error()})
	}
	return nil
}

func Excel(ctx *gin.Context, fileName string, content []byte) error {
	ctx.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%q", fileName))
	ctx.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", content)
	return nil
}
