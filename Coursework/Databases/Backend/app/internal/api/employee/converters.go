package employee

import (
	"coffee-shop-backend/app/internal/models"
	"coffee-shop-backend/app/pkg/apimodels"
)

func toEmployeeResponse(employee models.Employee) apimodels.Employee {
	return apimodels.Employee{
		ID:       employee.ID,
		FullName: employee.FullName,
		Login:    employee.Login,
		Phone:    employee.Phone,
	}
}
