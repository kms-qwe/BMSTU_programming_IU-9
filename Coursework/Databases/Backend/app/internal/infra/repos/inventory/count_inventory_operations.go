package inventory

import (
	"context"

	baserepos "coffee-shop-backend/app/internal/infra/repos"
	"coffee-shop-backend/app/internal/models"
	"coffee-shop-backend/app/internal/utils"
)

func (r *Repo) CountInventoryOperations(ctx context.Context, filter models.InventoryOperationFilter) (int, error) {
	builder := baserepos.ApplyInventoryFilters(
		r.psql.
			Select("COUNT(*)").
			From(inventoryOperationsTable+" io").
			Join("ingredients i ON i.id = io.ingredient_id").
			LeftJoin("orders o ON o.id = io.order_id").
			Join("shifts sh ON sh.id = io.shift_id").
			Join("employees e ON e.id = sh.duty"),
		filter,
	)

	query, args, err := builder.ToSql()
	if err != nil {
		return 0, utils.InternalError("failed to build count inventory operations query", err)
	}

	var total int
	if err := r.provider.GetExecutor(ctx).QueryRow(ctx, query, args...).Scan(&total); err != nil {
		return 0, utils.InternalError("failed to count inventory operations", err)
	}

	return total, nil
}
