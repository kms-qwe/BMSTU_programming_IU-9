package menuitem

import (
	"context"
	"errors"

	"coffee-shop-backend/app/internal/models"
	"coffee-shop-backend/app/internal/utils"

	"github.com/jackc/pgx/v5"
)

func (r *Repo) SetMenuItemActive(ctx context.Context, menuItemID string, isActive bool) (*models.MenuItem, error) {
	builder := r.psql.
		Update(menuItemsTable).
		Set("is_deleted", !isActive).
		Where("id = ?", menuItemID).
		Suffix("RETURNING id")

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, utils.InternalError("failed to build update menu item activity query", err)
	}

	rows, err := r.provider.GetExecutor(ctx).Query(ctx, query, args...)
	if err != nil {
		return nil, utils.InternalError("failed to update menu item activity", err)
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
