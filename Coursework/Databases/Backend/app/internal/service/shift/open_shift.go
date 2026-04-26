package shift

import (
	"context"
	"strings"

	"coffee-shop-backend/app/internal/models"
	"coffee-shop-backend/app/internal/utils"
)

func (s *Service) OpenShift(ctx context.Context, login, password string) (*models.Shift, error) {
	login = strings.TrimSpace(login)
	if login == "" || password == "" {
		return nil, utils.ValidationError("login and password are required", nil)
	}
	var shift *models.Shift
	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		activeShift, err := s.shiftRepo.GetActiveShift(ctx, true)
		if err != nil {
			return err
		}
		if activeShift != nil {
			return utils.ConflictError(utils.ErrorCodeShiftAlreadyOpened, "Смена уже открыта", nil)
		}
		employee, err := s.employeeRepo.GetEmployeeByLogin(ctx, login)
		if err != nil {
			return err
		}
		if employee == nil || !passwordMatches(employee.PasswordHash, password) {
			return utils.UnauthorizedError("Неверный логин или пароль")
		}
		shift, err = s.shiftRepo.CreateShift(ctx, employee.ID)
		return err
	})
	if err != nil {
		return nil, err
	}
	return shift, nil
}
