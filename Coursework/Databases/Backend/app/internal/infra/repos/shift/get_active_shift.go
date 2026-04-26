package shift

import (
	"context"
	"errors"

	"coffee-shop-backend/app/internal/models"
	"coffee-shop-backend/app/internal/utils"

	"github.com/jackc/pgx/v5"
)

func (r *Repo) GetActiveShift(ctx context.Context, forUpdate bool) (*models.Shift, error) {
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
		From(shiftsTable + " s").
		Join("employees e ON e.id = s.duty").
		Where("s.is_active = TRUE").
		OrderBy("s.opened_at DESC").
		Limit(1)

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, utils.InternalError("failed to build active shift query", err)
	}
	if forUpdate {
		query += " FOR UPDATE OF s"
	}

	rows, err := r.provider.GetExecutor(ctx).Query(ctx, query, args...)
	if err != nil {
		return nil, utils.InternalError("failed to get active shift", err)
	}
	defer rows.Close()

	dbItem, err := pgx.CollectOneRow(rows, pgx.RowToStructByNameLax[shiftDBModel])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, utils.InternalError("failed to collect active shift", err)
	}

	serviceItem := toServiceShift(dbItem)
	return &serviceItem, nil
}
