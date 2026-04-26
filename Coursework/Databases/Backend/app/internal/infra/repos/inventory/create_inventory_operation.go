package inventory

import (
	"context"

	"coffee-shop-backend/app/internal/models"
	"coffee-shop-backend/app/internal/utils"

	"github.com/jackc/pgx/v5"
)

func (r *Repo) CreateInventoryOperation(ctx context.Context, params models.InventoryOperationCreateRecord) (*models.InventoryOperation, error) {
	builder := r.psql.
		Insert(inventoryOperationsTable).
		Columns("shift_id", "ingredient_id", "order_id", "change_amount", "operation_type").
		Values(params.ShiftID, params.IngredientID, params.OrderID, params.ChangeAmount, params.OperationType).
		Suffix("RETURNING id")

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, utils.InternalError("failed to build create inventory operation query", err)
	}

	rows, err := r.provider.GetExecutor(ctx).Query(ctx, query, args...)
	if err != nil {
		return nil, utils.InternalError("failed to create inventory operation", err)
	}
	defer rows.Close()

	dbItem, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[createdInventoryOperationDBModel])
	if err != nil {
		return nil, utils.InternalError("failed to collect created inventory operation", err)
	}

	return r.getInventoryOperationByID(ctx, dbItem.ID)
}
