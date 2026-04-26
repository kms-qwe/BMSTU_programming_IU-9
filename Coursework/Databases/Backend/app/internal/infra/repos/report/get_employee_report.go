package report

import (
	"context"

	"coffee-shop-backend/app/internal/models"
	"coffee-shop-backend/app/internal/utils"

	"github.com/jackc/pgx/v5"
)

func (r *Repo) GetEmployeeReport(ctx context.Context, filter models.ReportFilter) (*models.EmployeeReportData, error) {
	summaryRows, err := r.provider.GetExecutor(ctx).Query(ctx, `
		WITH shift_counts AS (
			SELECT e.id, COUNT(sh.id) AS shift_count
			FROM employees e
			LEFT JOIN shifts sh ON sh.duty = e.id
				AND sh.opened_at BETWEEN $1 AND $2
			WHERE ($3::uuid IS NULL OR e.id = $3::uuid)
			GROUP BY e.id
		),
		order_totals AS (
			SELECT
				e.id AS employee_id,
				e.full_name,
				o.id AS order_id,
				SUM(oi.price * oi.quantity) AS total_price
			FROM orders o
			JOIN shifts sh ON sh.id = o.shift_id
			JOIN employees e ON e.id = sh.duty
			JOIN order_items oi ON oi.order_id = o.id
			WHERE o.created_at BETWEEN $1 AND $2
			  AND ($3::uuid IS NULL OR e.id = $3::uuid)
			GROUP BY e.id, e.full_name, o.id
		)
		SELECT
			e.full_name AS employee_name,
			COALESCE(sc.shift_count, 0) AS shift_count,
			COUNT(ot.order_id) AS orders_count,
			COALESCE(SUM(ot.total_price), 0) AS revenue,
			COALESCE(AVG(ot.total_price), 0) AS average_bill
		FROM employees e
		LEFT JOIN shift_counts sc ON sc.id = e.id
		LEFT JOIN order_totals ot ON ot.employee_id = e.id
		WHERE ($3::uuid IS NULL OR e.id = $3::uuid)
		GROUP BY e.full_name, sc.shift_count
		ORDER BY e.full_name ASC
	`, filter.From, filter.To, filter.EmployeeID)
	if err != nil {
		return nil, utils.InternalError("failed to load employee summary report", err)
	}
	defer summaryRows.Close()

	summaryDBRows, err := pgx.CollectRows(summaryRows, pgx.RowToStructByName[employeeSummaryRowDBModel])
	if err != nil {
		return nil, utils.InternalError("failed to collect employee summary rows", err)
	}

	orderRows, err := r.provider.GetExecutor(ctx).Query(ctx, `
		SELECT
			o.created_at AS created_at,
			e.full_name AS employee_name,
			o.id AS order_id,
			SUM(oi.price * oi.quantity) AS total_price
		FROM orders o
		JOIN shifts sh ON sh.id = o.shift_id
		JOIN employees e ON e.id = sh.duty
		JOIN order_items oi ON oi.order_id = o.id
		WHERE o.created_at BETWEEN $1 AND $2
		  AND ($3::uuid IS NULL OR e.id = $3::uuid)
		GROUP BY o.created_at, e.full_name, o.id
		ORDER BY o.created_at DESC
	`, filter.From, filter.To, filter.EmployeeID)
	if err != nil {
		return nil, utils.InternalError("failed to load employee order report", err)
	}
	defer orderRows.Close()

	orderDBRows, err := pgx.CollectRows(orderRows, pgx.RowToStructByName[employeeOrderRowDBModel])
	if err != nil {
		return nil, utils.InternalError("failed to collect employee order rows", err)
	}

	shiftRows, err := r.provider.GetExecutor(ctx).Query(ctx, `
		SELECT
			e.full_name AS employee_name,
			sh.opened_at AS opened_at,
			sh.closed_at AS closed_at
		FROM shifts sh
		JOIN employees e ON e.id = sh.duty
		WHERE sh.opened_at BETWEEN $1 AND $2
		  AND ($3::uuid IS NULL OR e.id = $3::uuid)
		ORDER BY sh.opened_at DESC
	`, filter.From, filter.To, filter.EmployeeID)
	if err != nil {
		return nil, utils.InternalError("failed to load employee shifts report", err)
	}
	defer shiftRows.Close()

	shiftDBRows, err := pgx.CollectRows(shiftRows, pgx.RowToStructByName[employeeShiftRowDBModel])
	if err != nil {
		return nil, utils.InternalError("failed to collect employee shift rows", err)
	}

	return &models.EmployeeReportData{
		Summaries: toServiceEmployeeSummaryRows(summaryDBRows),
		Orders:    toServiceEmployeeOrderRows(orderDBRows),
		Shifts:    toServiceEmployeeShiftRows(shiftDBRows),
	}, nil
}
