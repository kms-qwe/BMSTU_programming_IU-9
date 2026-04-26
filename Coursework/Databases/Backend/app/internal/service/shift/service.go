package shift

import (
	"context"
	"time"

	"coffee-shop-backend/app/internal/models"
	"coffee-shop-backend/app/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

type IShiftRepo interface {
	GetActiveShift(ctx context.Context, forUpdate bool) (*models.Shift, error)
	CreateShift(ctx context.Context, employeeID string) (*models.Shift, error)
	CloseActiveShift(ctx context.Context, closedAt time.Time) (*models.Shift, error)
}

type IEmployeeRepo interface {
	GetEmployeeByLogin(ctx context.Context, login string) (*models.EmployeeCredentials, error)
}

type Service struct {
	shiftRepo    IShiftRepo
	employeeRepo IEmployeeRepo
	txManager    utils.TxManager
}

func New(shiftRepo IShiftRepo, employeeRepo IEmployeeRepo, txManager utils.TxManager) *Service {
	return &Service{shiftRepo: shiftRepo, employeeRepo: employeeRepo, txManager: txManager}
}

func passwordMatches(hash, password string) bool {
	if hash == password {
		return true
	}
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}
