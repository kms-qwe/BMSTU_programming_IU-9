package employee

import (
	"context"

	"coffee-shop-backend/app/internal/models"
)

type IEmployeeRepo interface {
	ListEmployees(ctx context.Context) ([]models.Employee, error)
}

type Service struct {
	repo IEmployeeRepo
}

func New(repo IEmployeeRepo) *Service {
	return &Service{repo: repo}
}
