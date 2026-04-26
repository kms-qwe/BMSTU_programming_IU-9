package shift

import (
	"coffee-shop-backend/app/internal/models"
	"coffee-shop-backend/app/pkg/apimodels"
)

func toShiftResponse(shift *models.Shift) *apimodels.Shift {
	if shift == nil {
		return nil
	}

	var duty *apimodels.DutyEmployee
	if shift.Duty != nil {
		duty = &apimodels.DutyEmployee{
			ID:       shift.Duty.ID,
			FullName: shift.Duty.FullName,
			Login:    shift.Duty.Login,
		}
	}

	return &apimodels.Shift{
		ID:       shift.ID,
		OpenedAt: shift.OpenedAt,
		ClosedAt: shift.ClosedAt,
		IsActive: shift.IsActive,
		Duty:     duty,
	}
}
