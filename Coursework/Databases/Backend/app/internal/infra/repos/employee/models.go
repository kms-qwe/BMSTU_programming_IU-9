package employee

import "coffee-shop-backend/app/internal/models"

type employeeDBModel struct {
	ID       string `db:"id"`
	FullName string `db:"full_name"`
	Login    string `db:"login"`
	Phone    string `db:"phone"`
}

type employeeCredentialsDBModel struct {
	ID           string `db:"id"`
	FullName     string `db:"full_name"`
	Login        string `db:"login"`
	Phone        string `db:"phone"`
	PasswordHash string `db:"password_hash"`
}

func toServiceEmployee(item employeeDBModel) models.Employee {
	return models.Employee{
		ID:       item.ID,
		FullName: item.FullName,
		Login:    item.Login,
		Phone:    item.Phone,
	}
}

func toServiceEmployees(items []employeeDBModel) []models.Employee {
	result := make([]models.Employee, 0, len(items))
	for _, item := range items {
		result = append(result, toServiceEmployee(item))
	}

	return result
}

func toServiceEmployeeCredentials(item employeeCredentialsDBModel) models.EmployeeCredentials {
	return models.EmployeeCredentials{
		Employee: models.Employee{
			ID:       item.ID,
			FullName: item.FullName,
			Login:    item.Login,
			Phone:    item.Phone,
		},
		PasswordHash: item.PasswordHash,
	}
}
