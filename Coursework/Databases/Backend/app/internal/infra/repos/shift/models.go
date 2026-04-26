package shift

import (
	"time"

	"coffee-shop-backend/app/internal/models"
)

type shiftDBModel struct {
	ID           string     `db:"id"`
	OpenedAt     time.Time  `db:"opened_at"`
	ClosedAt     *time.Time `db:"closed_at"`
	IsActive     bool       `db:"is_active"`
	DutyID       string     `db:"duty_id"`
	DutyFullName string     `db:"duty_full_name"`
	DutyLogin    string     `db:"duty_login"`
}

type createdShiftDBModel struct {
	ID string `db:"id"`
}

func toServiceShift(item shiftDBModel) models.Shift {
	return models.Shift{
		ID:       item.ID,
		OpenedAt: item.OpenedAt,
		ClosedAt: item.ClosedAt,
		IsActive: item.IsActive,
		Duty: &models.ShiftEmployee{
			ID:       item.DutyID,
			FullName: item.DutyFullName,
			Login:    item.DutyLogin,
		},
	}
}
