package order

import (
	"context"
	"errors"

	"coffee-shop-backend/app/internal/models"
	"coffee-shop-backend/app/internal/utils"

	"github.com/jackc/pgx/v5"
)

func (r *Repo) GetOrder(ctx context.Context, orderID string) (*models.Order, error) {
	builder := r.psql.
		Select(
			"o.id AS id",
			"o.created_at AS created_at",
			"s.id AS shift_id",
			"s.opened_at AS shift_opened_at",
			"e.id AS employee_id",
			"e.full_name AS employee_full_name",
			"COALESCE(SUM(oi.price * oi.quantity), 0) AS total_price",
		).
		From(ordersTable+" o").
		Join("shifts s ON s.id = o.shift_id").
		Join("employees e ON e.id = s.duty").
		LeftJoin("order_items oi ON oi.order_id = o.id").
		Where("o.id = ?", orderID).
		GroupBy("o.id", "o.created_at", "s.id", "s.opened_at", "e.id", "e.full_name")

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, utils.InternalError("failed to build get order query", err)
	}

	rows, err := r.provider.GetExecutor(ctx).Query(ctx, query, args...)
	if err != nil {
		return nil, utils.InternalError("failed to get order", err)
	}
	defer rows.Close()

	dbItem, err := pgx.CollectOneRow(rows, pgx.RowToStructByNameLax[orderDBModel])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, utils.InternalError("failed to collect order", err)
	}

	order := toServiceOrder(dbItem)

	itemsByOrderID, err := r.loadOrderItems(ctx, []string{orderID})
	if err != nil {
		return nil, err
	}
	order.Items = itemsByOrderID[orderID]

	order.InventoryOperations, err = r.loadInventoryOperationsByOrderID(ctx, orderID)
	if err != nil {
		return nil, err
	}

	return &order, nil
}
