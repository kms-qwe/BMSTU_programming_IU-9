package report

import (
	"context"

	"coffee-shop-backend/app/internal/models"
	"coffee-shop-backend/app/internal/utils"

	"github.com/jackc/pgx/v5"
)

func (r *Repo) GetMenuPopularityReport(ctx context.Context, filter models.ReportFilter) (*models.MenuPopularityReportData, error) {
	rows, err := r.provider.GetExecutor(ctx).Query(ctx, `
		WITH sales AS (
			SELECT
				m.name AS menu_item_name,
				SUM(oi.quantity) AS sales_count,
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
			ROW_NUMBER() OVER (ORDER BY revenue DESC, menu_item_name ASC) AS place,
			menu_item_name,
			sales_count,
			revenue,
			CASE WHEN SUM(revenue) OVER () = 0 THEN 0 ELSE revenue * 100.0 / SUM(revenue) OVER () END AS revenue_share_pct
		FROM sales
		ORDER BY revenue DESC, menu_item_name ASC
	`, filter.From, filter.To, filter.EmployeeID)
	if err != nil {
		return nil, utils.InternalError("failed to load menu popularity report", err)
	}
	defer rows.Close()

	dbRows, err := pgx.CollectRows(rows, pgx.RowToStructByName[menuPopularityRowDBModel])
	if err != nil {
		return nil, utils.InternalError("failed to collect menu popularity rows", err)
	}

	return &models.MenuPopularityReportData{
		Rows: toServiceMenuPopularityRows(dbRows),
	}, nil
}
