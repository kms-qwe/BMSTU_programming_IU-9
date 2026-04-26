package shift

import (
	"context"
	"time"

	"coffee-shop-backend/app/internal/models"
	"coffee-shop-backend/app/internal/utils"
)

func (s *Service) CloseShift(ctx context.Context) (*models.Shift, error) {
	var shift *models.Shift
	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		activeShift, err := s.shiftRepo.GetActiveShift(ctx, true)
		if err != nil {
			return err
		}
		if activeShift == nil {
			return utils.BusinessError(utils.ErrorCodeShiftNotActive, "Нет активной смены", nil)
		}
		shift, err = s.shiftRepo.CloseActiveShift(ctx, time.Now().UTC())
		return err
	})
	if err != nil {
		return nil, err
	}
	return shift, nil
}
