package report

import (
	"context"

	"coffee-shop-backend/app/internal/models"
	"coffee-shop-backend/app/internal/utils"

	"github.com/jackc/pgx/v5"
)

func (r *Repo) GetSalesReport(ctx context.Context, filter models.ReportFilter) (*models.SalesReportData, error) {
	summaryRows, err := r.provider.GetExecutor(ctx).Query(ctx, `
		WITH order_totals AS (
			SELECT o.id, SUM(oi.price * oi.quantity) AS total_price
			FROM orders o
			JOIN shifts sh ON sh.id = o.shift_id
			JOIN order_items oi ON oi.order_id = o.id
			WHERE o.created_at BETWEEN $1 AND $2
			  AND ($3::uuid IS NULL OR sh.duty = $3::uuid)
			GROUP BY o.id
		)
		SELECT
			COUNT(*) AS orders_count,
			COALESCE(SUM(total_price), 0) AS revenue,
			COALESCE(AVG(total_price), 0) AS average_bill
		FROM order_totals
	`, filter.From, filter.To, filter.EmployeeID)
	if err != nil {
		return nil, utils.InternalError("failed to load sales summary", err)
	}
	defer summaryRows.Close()

	summaryDB, err := pgx.CollectOneRow(summaryRows, pgx.RowToStructByName[salesSummaryDBModel])
	if err != nil {
		return nil, utils.InternalError("failed to collect sales summary", err)
	}

	menuRows, err := r.provider.GetExecutor(ctx).Query(ctx, `
		WITH sales AS (
			SELECT
				m.name AS menu_item_name,
				SUM(oi.quantity) AS quantity,
				SUM(oi.price * oi.quantity) AS revenue
			FROM order_items oi
			JOIN orders o ON o.id = oi.order_id
			JOIN shifts sh ON sh.id = o.shift_id
			JOIN menu_items m ON m.id = oi.menu_item_id
			WHERE o.created_at BETWEEN $1 AND $2
			  AND ($3::uuid IS NULL OR sh.duty = $3::uuid)
			GROUP BY m.name
		)
		SELECT
			menu_item_name,
			quantity,
			revenue,
			CASE WHEN SUM(revenue) OVER () = 0 THEN 0 ELSE revenue * 100.0 / SUM(revenue) OVER () END AS sales_share_pct
		FROM sales
		ORDER BY revenue DESC, menu_item_name ASC
	`, filter.From, filter.To, filter.EmployeeID)
	if err != nil {
		return nil, utils.InternalError("failed to load sales menu items", err)
	}
	defer menuRows.Close()

	menuDBRows, err := pgx.CollectRows(menuRows, pgx.RowToStructByName[salesMenuItemRowDBModel])
	if err != nil {
		return nil, utils.InternalError("failed to collect sales menu items", err)
	}

	orderRows, err := r.provider.GetExecutor(ctx).Query(ctx, `
		SELECT
			o.created_at AS created_at,
			o.id AS order_id,
			e.full_name AS employee_name,
			SUM(oi.price * oi.quantity) AS total_price
		FROM orders o
		JOIN shifts sh ON sh.id = o.shift_id
		JOIN employees e ON e.id = sh.duty
		JOIN order_items oi ON oi.order_id = o.id
		WHERE o.created_at BETWEEN $1 AND $2
		  AND ($3::uuid IS NULL OR sh.duty = $3::uuid)
		GROUP BY o.created_at, o.id, e.full_name
		ORDER BY o.created_at DESC
	`, filter.From, filter.To, filter.EmployeeID)
	if err != nil {
		return nil, utils.InternalError("failed to load sales orders", err)
	}
	defer orderRows.Close()

	orderDBRows, err := pgx.CollectRows(orderRows, pgx.RowToStructByName[salesOrderRowDBModel])
	if err != nil {
		return nil, utils.InternalError("failed to collect sales orders", err)
	}

	return &models.SalesReportData{
		Summary:   toServiceSalesSummary(summaryDB),
		MenuItems: toServiceSalesMenuRows(menuDBRows),
		Orders:    toServiceSalesOrderRows(orderDBRows),
	}, nil
}
