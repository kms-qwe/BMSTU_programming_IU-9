package shift

import (
	"context"

	"coffee-shop-backend/app/internal/models"
	"coffee-shop-backend/app/internal/utils"

	"github.com/jackc/pgx/v5"
)

func (r *Repo) getShiftByID(ctx context.Context, shiftID string) (*models.Shift, error) {
	builder := r.psql.
		Select(
			"s.id AS id",
			"s.opened_at AS opened_at",
			"s.closed_at AS closed_at",
			"s.is_active AS is_active",
			"e.id AS duty_id",
			"e.full_name AS duty_full_name",
			"e.login AS duty_login",
		).
		From(shiftsTable+" s").
		Join("employees e ON e.id = s.duty").
		Where("s.id = ?", shiftID)

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, utils.InternalError("failed to build shift by id query", err)
	}

	rows, err := r.provider.GetExecutor(ctx).Query(ctx, query, args...)
	if err != nil {
		return nil, utils.InternalError("failed to get shift by id", err)
	}
	defer rows.Close()

	dbItem, err := pgx.CollectOneRow(rows, pgx.RowToStructByNameLax[shiftDBModel])
	if err != nil {
		return nil, utils.InternalError("failed to collect shift by id", err)
	}

	serviceItem := toServiceShift(dbItem)
	return &serviceItem, nil
}
