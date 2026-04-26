package order

import (
	"context"

	baserepos "coffee-shop-backend/app/internal/infra/repos"
	"coffee-shop-backend/app/internal/models"
	"coffee-shop-backend/app/internal/utils"

	"github.com/jackc/pgx/v5"
)

func (r *Repo) ListOrders(ctx context.Context, filter models.OrderFilter) ([]models.Order, error) {
	builder := baserepos.ApplyOrderFilters(
		r.psql.
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
			LeftJoin("order_items oi ON oi.order_id = o.id"),
		filter,
	).
		GroupBy("o.id", "o.created_at", "s.id", "s.opened_at", "e.id", "e.full_name").
		OrderBy("o.created_at DESC").
		Limit(uint64(filter.PageSize)).
		Offset(uint64((filter.Page - 1) * filter.PageSize))

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, utils.InternalError("failed to build list orders query", err)
	}

	rows, err := r.provider.GetExecutor(ctx).Query(ctx, query, args...)
	if err != nil {
		return nil, utils.InternalError("failed to list orders", err)
	}
	defer rows.Close()

	dbItems, err := pgx.CollectRows(rows, pgx.RowToStructByNameLax[orderDBModel])
	if err != nil {
		return nil, utils.InternalError("failed to collect orders", err)
	}

	orders := make([]models.Order, 0, len(dbItems))
	orderIDs := make([]string, 0, len(dbItems))
	for _, item := range dbItems {
		orders = append(orders, toServiceOrder(item))
		orderIDs = append(orderIDs, item.ID)
	}

	if len(orderIDs) == 0 {
		return orders, nil
	}

	itemsByOrderID, err := r.loadOrderItems(ctx, orderIDs)
	if err != nil {
		return nil, err
	}

	for idx := range orders {
		orders[idx].Items = itemsByOrderID[orders[idx].ID]
	}

	return orders, nil
}
