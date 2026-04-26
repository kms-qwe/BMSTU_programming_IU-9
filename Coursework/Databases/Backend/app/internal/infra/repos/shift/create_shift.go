package shift

import (
	"context"

	"coffee-shop-backend/app/internal/models"
	"coffee-shop-backend/app/internal/utils"

	"github.com/jackc/pgx/v5"
)

func (r *Repo) CreateShift(ctx context.Context, employeeID string) (*models.Shift, error) {
	builder := r.psql.
		Insert(shiftsTable).
		Columns("duty").
		Values(employeeID).
		Suffix("RETURNING id")

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, utils.InternalError("failed to build create shift query", err)
	}

	rows, err := r.provider.GetExecutor(ctx).Query(ctx, query, args...)
	if err != nil {
		return nil, utils.InternalError("failed to create shift", err)
	}
	defer rows.Close()

	dbItem, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[createdShiftDBModel])
	if err != nil {
		return nil, utils.InternalError("failed to collect created shift", err)
	}

	return r.getShiftByID(ctx, dbItem.ID)
}
