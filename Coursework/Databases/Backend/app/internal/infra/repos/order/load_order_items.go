package order

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"

	"coffee-shop-backend/app/internal/models"
	"coffee-shop-backend/app/internal/utils"
)

func (r *Repo) loadOrderItems(ctx context.Context, orderIDs []string) (map[string][]models.OrderItem, error) {
	builder := r.psql.
		Select(
			"oi.id AS id",
			"oi.order_id AS order_id",
			"oi.menu_item_id AS menu_item_id",
			"m.name AS menu_item_name",
			"oi.quantity AS quantity",
			"oi.price AS price",
		).
		From("order_items oi").
		Join("menu_items m ON m.id = oi.menu_item_id").
		Where(sq.Eq{"oi.order_id": orderIDs}).
		OrderBy("oi.order_id", "m.name")

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, utils.InternalError("failed to build load order items query", err)
	}

	rows, err := r.provider.GetExecutor(ctx).Query(ctx, query, args...)
	if err != nil {
		return nil, utils.InternalError("failed to load order items", err)
	}
	defer rows.Close()

	dbItems, err := pgx.CollectRows(rows, pgx.RowToStructByNameLax[orderItemDBModel])
	if err != nil {
		return nil, utils.InternalError("failed to collect order items", err)
	}

	result := make(map[string][]models.OrderItem)
	for _, item := range dbItems {
		result[item.OrderID] = append(result[item.OrderID], toServiceOrderItems([]orderItemDBModel{item})...)
	}

	return result, nil
}
