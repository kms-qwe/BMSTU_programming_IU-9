package menuitem

import (
	"context"
	"errors"

	"coffee-shop-backend/app/internal/models"
	"coffee-shop-backend/app/internal/utils"

	"github.com/jackc/pgx/v5"
)

func (r *Repo) UpdateMenuItemPrice(ctx context.Context, menuItemID string, price float64) (*models.MenuItem, error) {
	builder := r.psql.
		Update(menuItemsTable).
		Set("price", price).
		Where("id = ?", menuItemID).
		Suffix("RETURNING id")

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, utils.InternalError("failed to build update menu item price query", err)
	}

	rows, err := r.provider.GetExecutor(ctx).Query(ctx, query, args...)
	if err != nil {
		return nil, utils.InternalError("failed to update menu item price", err)
	}
	defer rows.Close()

	dbItem, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[createdMenuItemDBModel])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, utils.InternalError("failed to collect updated menu item", err)
	}

	return r.GetMenuItem(ctx, dbItem.ID)
}
