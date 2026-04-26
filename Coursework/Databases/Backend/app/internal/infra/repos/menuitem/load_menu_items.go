package menuitem

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"

	"coffee-shop-backend/app/internal/models"
	"coffee-shop-backend/app/internal/utils"
)

func (r *Repo) loadMenuItems(ctx context.Context, includeInactive bool, menuItemIDs []string) ([]models.MenuItem, error) {
	builder := r.psql.
		Select(
			"m.id AS id",
			"m.name AS name",
			"m.price AS price",
			"NOT m.is_deleted AS is_active",
			"ri.ingredient_id AS ingredient_id",
			"i.name AS ingredient_name",
			"i.unit AS unit",
			"ri.amount AS amount",
			"CASE WHEN i.id IS NULL THEN TRUE ELSE NOT i.is_deleted END AS ingredient_is_active",
		).
		From(menuItemsTable+" m").
		LeftJoin("recipe_ingredients ri ON ri.menu_item_id = m.id").
		LeftJoin("ingredients i ON i.id = ri.ingredient_id").
		OrderBy("m.name ASC", "i.name ASC")

	if includeInactive {
		builder = builder.Where("TRUE")
	} else {
		builder = builder.Where("NOT m.is_deleted")
	}
	if len(menuItemIDs) > 0 {
		builder = builder.Where(sq.Eq{"m.id": menuItemIDs})
	}

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, utils.InternalError("failed to build menu items query", err)
	}

	rows, err := r.provider.GetExecutor(ctx).Query(ctx, query, args...)
	if err != nil {
		return nil, utils.InternalError("failed to load menu items", err)
	}
	defer rows.Close()

	dbItems, err := pgx.CollectRows(rows, pgx.RowToStructByNameLax[menuItemRowDBModel])
	if err != nil {
		return nil, utils.InternalError("failed to collect menu items", err)
	}

	return toServiceMenuItems(dbItems), nil
}
