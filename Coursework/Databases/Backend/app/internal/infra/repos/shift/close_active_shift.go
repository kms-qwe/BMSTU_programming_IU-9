package shift

import (
	"context"
	"errors"
	"time"

	"coffee-shop-backend/app/internal/models"
	"coffee-shop-backend/app/internal/utils"

	"github.com/jackc/pgx/v5"
)

func (r *Repo) CloseActiveShift(ctx context.Context, closedAt time.Time) (*models.Shift, error) {
	builder := r.psql.
		Update(shiftsTable).
		Set("closed_at", closedAt).
		Set("is_active", false).
		Where("is_active = TRUE").
		Suffix("RETURNING id")

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, utils.InternalError("failed to build close shift query", err)
	}

	rows, err := r.provider.GetExecutor(ctx).Query(ctx, query, args...)
	if err != nil {
		return nil, utils.InternalError("failed to close shift", err)
	}
	defer rows.Close()

	dbItem, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[createdShiftDBModel])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, utils.InternalError("failed to collect closed shift", err)
	}

	return r.getShiftByID(ctx, dbItem.ID)
}
