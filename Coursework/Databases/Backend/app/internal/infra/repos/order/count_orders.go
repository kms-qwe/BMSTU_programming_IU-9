package order

import (
	"context"

	baserepos "coffee-shop-backend/app/internal/infra/repos"
	"coffee-shop-backend/app/internal/models"
	"coffee-shop-backend/app/internal/utils"
)

func (r *Repo) CountOrders(ctx context.Context, filter models.OrderFilter) (int, error) {
	builder := baserepos.ApplyOrderFilters(
		r.psql.
			Select("COUNT(*)").
			From(ordersTable+" o").
			Join("shifts s ON s.id = o.shift_id").
			Join("employees e ON e.id = s.duty"),
		filter,
	)

	query, args, err := builder.ToSql()
	if err != nil {
		return 0, utils.InternalError("failed to build count orders query", err)
	}

	var total int
	if err := r.provider.GetExecutor(ctx).QueryRow(ctx, query, args...).Scan(&total); err != nil {
		return 0, utils.InternalError("failed to count orders", err)
	}

	return total, nil
}
