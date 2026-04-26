package report

import (
	"context"

	"coffee-shop-backend/app/internal/models"
	"coffee-shop-backend/app/internal/utils"

	"github.com/jackc/pgx/v5"
)

func (r *Repo) GetInventoryReport(ctx context.Context, filter models.InventoryReportFilter) (*models.InventoryReportData, error) {
	summaryRows, err := r.provider.GetExecutor(ctx).Query(ctx, `
		WITH filtered_ops AS (
			SELECT io.*
			FROM inventory_operations io
			JOIN shifts sh ON sh.id = io.shift_id
			WHERE io.created_at BETWEEN $1 AND $2
			  AND ($3::uuid IS NULL OR sh.duty = $3::uuid)
			  AND ($4::uuid IS NULL OR io.ingredient_id = $4::uuid)
			  AND ($5::text IS NULL OR io.operation_type = $5)
		)
		SELECT
			i.name AS ingredient_name,
			i.unit AS unit,
			COALESCE(SUM(CASE WHEN fo.operation_type IN ('MANUAL_RESTOCK', 'INITIAL_STOCK') AND fo.change_amount > 0 THEN fo.change_amount ELSE 0 END), 0) AS restocked,
			COALESCE(SUM(CASE WHEN fo.operation_type = 'AUTO_ORDER_WRITE_OFF' THEN ABS(fo.change_amount) ELSE 0 END), 0) AS auto_written_off,
			COALESCE(SUM(CASE WHEN fo.operation_type = 'MANUAL_WRITE_OFF' THEN ABS(fo.change_amount) ELSE 0 END), 0) AS manual_written_off,
			i.current_stock AS current_stock
		FROM ingredients i
		LEFT JOIN filtered_ops fo ON fo.ingredient_id = i.id
		WHERE ($4::uuid IS NULL OR i.id = $4::uuid)
		GROUP BY i.id, i.name, i.unit, i.current_stock
		ORDER BY i.name ASC
	`, filter.From, filter.To, filter.EmployeeID, filter.IngredientID, filter.OperationType)
	if err != nil {
		return nil, utils.InternalError("failed to load inventory summary report", err)
	}
	defer summaryRows.Close()

	summaryDBRows, err := pgx.CollectRows(summaryRows, pgx.RowToStructByName[inventorySummaryRowDBModel])
	if err != nil {
		return nil, utils.InternalError("failed to collect inventory summary rows", err)
	}

	operationRows, err := r.provider.GetExecutor(ctx).Query(ctx, `
		SELECT
			io.id AS id,
			io.created_at AS created_at,
			io.operation_type AS operation_type,
			io.change_amount AS change_amount,
			i.id AS ingredient_id,
			i.name AS ingredient_name,
			i.unit AS ingredient_unit,
			i.current_stock AS ingredient_current_stock,
			o.id AS order_id,
			sh.id AS shift_id,
			e.id AS employee_id,
			e.full_name AS employee_full_name
		FROM inventory_operations io
		JOIN ingredients i ON i.id = io.ingredient_id
		LEFT JOIN orders o ON o.id = io.order_id
		JOIN shifts sh ON sh.id = io.shift_id
		JOIN employees e ON e.id = sh.duty
		WHERE io.created_at BETWEEN $1 AND $2
		  AND ($3::uuid IS NULL OR e.id = $3::uuid)
		  AND ($4::uuid IS NULL OR i.id = $4::uuid)
		  AND ($5::text IS NULL OR io.operation_type = $5)
		ORDER BY io.created_at DESC
	`, filter.From, filter.To, filter.EmployeeID, filter.IngredientID, filter.OperationType)
	if err != nil {
		return nil, utils.InternalError("failed to load inventory operations report", err)
	}
	defer operationRows.Close()

	operationDBRows, err := pgx.CollectRows(operationRows, pgx.RowToStructByNameLax[reportInventoryOperationDBModel])
	if err != nil {
		return nil, utils.InternalError("failed to collect inventory operations report", err)
	}

	return &models.InventoryReportData{
		Summaries:  toServiceInventorySummaryRows(summaryDBRows),
		Operations: toServiceReportInventoryOperations(operationDBRows),
	}, nil
}
