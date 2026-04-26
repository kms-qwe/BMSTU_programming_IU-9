package inventory

import (
	"context"
	"errors"

	"coffee-shop-backend/app/internal/models"
	"coffee-shop-backend/app/internal/utils"

	"github.com/jackc/pgx/v5"
)

func (r *Repo) getInventoryOperationByID(ctx context.Context, operationID string) (*models.InventoryOperation, error) {
	builder := r.psql.
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
		Join("employees e ON e.id = sh.duty").
		Where("io.id = ?", operationID)

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, utils.InternalError("failed to build get inventory operation query", err)
	}

	rows, err := r.provider.GetExecutor(ctx).Query(ctx, query, args...)
	if err != nil {
		return nil, utils.InternalError("failed to get inventory operation", err)
	}
	defer rows.Close()

	dbItem, err := pgx.CollectOneRow(rows, pgx.RowToStructByNameLax[inventoryOperationDBModel])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, utils.InternalError("failed to collect inventory operation", err)
	}

	serviceItem := toServiceInventoryOperation(dbItem)
	return &serviceItem, nil
}
