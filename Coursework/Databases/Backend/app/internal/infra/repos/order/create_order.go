package order

import (
	"context"

	"coffee-shop-backend/app/internal/models"
	"coffee-shop-backend/app/internal/utils"

	"github.com/jackc/pgx/v5"
)

func (r *Repo) CreateOrder(ctx context.Context, shiftID string) (*models.Order, error) {
	builder := r.psql.
		Insert(ordersTable).
		Columns("shift_id").
		Values(shiftID).
		Suffix("RETURNING id, created_at")

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, utils.InternalError("failed to build create order query", err)
	}

	rows, err := r.provider.GetExecutor(ctx).Query(ctx, query, args...)
	if err != nil {
		return nil, utils.InternalError("failed to create order", err)
	}
	defer rows.Close()

	item, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[createdOrderDBModel])
	if err != nil {
		return nil, utils.InternalError("failed to collect created order", err)
	}

	return &models.Order{
		ID:        item.ID,
		CreatedAt: item.CreatedAt,
	}, nil
}
