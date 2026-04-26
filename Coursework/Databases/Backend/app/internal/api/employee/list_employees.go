package employee

import (
	"coffee-shop-backend/app/pkg/apimodels"

	"github.com/gin-gonic/gin"
)

// listEmployees godoc
// @Summary List employees
// @Description Returns employees for filters and operational screens.
// @Tags employee
// @Produce json
// @Success 200 {object} apimodels.ListEmployeesResponse
// @Failure 500 {object} apimodels.ErrorResponse
// @Router /api/employee [get]
func (a *API) listEmployees(ctx *gin.Context) error {
	employees, err := a.service.ListEmployees(ctx.Request.Context())
	if err != nil {
		return err
	}

	items := make([]apimodels.Employee, 0, len(employees))
	for _, employee := range employees {
		items = append(items, toEmployeeResponse(employee))
	}

	ctx.JSON(200, apimodels.ListEmployeesResponse{Employees: items})
	return nil
}
