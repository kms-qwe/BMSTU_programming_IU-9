package inventory

import (
	"context"

	baserepos "coffee-shop-backend/app/internal/infra/repos"
	"coffee-shop-backend/app/internal/models"
	"coffee-shop-backend/app/internal/utils"

	"github.com/jackc/pgx/v5"
)

func (r *Repo) ListInventoryOperations(ctx context.Context, filter models.InventoryOperationFilter) ([]models.InventoryOperation, error) {
	builder := baserepos.ApplyInventoryFilters(
		r.psql.
			Select(
				"io.id AS id",
				"io.created_at AS created_at",
				"io.operation_type AS operation_type",
				"io.change_amount AS change_amount",
				"i.id AS ingredient_id",
				"i.name AS ingredient_name",
				"i.unit AS ingredient_unit",
				"i.current_stock AS ingredient_current_stock",
				"o.id AS order_id",
				"sh.id AS shift_id",
				"e.id AS employee_id",
				"e.full_name AS employee_full_name",
			).
			From(inventoryOperationsTable+" io").
			Join("ingredients i ON i.id = io.ingredient_id").
			LeftJoin("orders o ON o.id = io.order_id").
			Join("shifts sh ON sh.id = io.shift_id").
			Join("employees e ON e.id = sh.duty"),
		filter,
	).
		OrderBy("io.created_at DESC").
		Limit(uint64(filter.PageSize)).
		Offset(uint64((filter.Page - 1) * filter.PageSize))

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, utils.InternalError("failed to build list inventory operations query", err)
	}

	rows, err := r.provider.GetExecutor(ctx).Query(ctx, query, args...)
	if err != nil {
		return nil, utils.InternalError("failed to list inventory operations", err)
	}
	defer rows.Close()

	dbItems, err := pgx.CollectRows(rows, pgx.RowToStructByNameLax[inventoryOperationDBModel])
	if err != nil {
		return nil, utils.InternalError("failed to collect inventory operations", err)
	}

	return toServiceInventoryOperations(dbItems), nil
}
