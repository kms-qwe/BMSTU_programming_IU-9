package menuitem

import (
	"context"

	baserepos "coffee-shop-backend/app/internal/infra/repos"
	"coffee-shop-backend/app/internal/models"
	"coffee-shop-backend/app/internal/utils"

	"github.com/jackc/pgx/v5"
)

func (r *Repo) CreateMenuItem(ctx context.Context, params models.CreateMenuItemParams) (*models.MenuItem, error) {
	builder := r.psql.
		Insert(menuItemsTable).
		Columns("name", "price").
		Values(params.Name, params.Price).
		Suffix("RETURNING id")

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, utils.InternalError("failed to build create menu item query", err)
	}

	rows, err := r.provider.GetExecutor(ctx).Query(ctx, query, args...)
	if err != nil {
		return nil, utils.InternalError("failed to create menu item", err)
	}
	defer rows.Close()

	dbItem, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[createdMenuItemDBModel])
	if err != nil {
		if baserepos.IsUniqueViolation(err) {
			return nil, utils.ConflictError(utils.ErrorCodeDuplicateName, "Пункт меню с таким названием уже существует", nil)
		}
		return nil, utils.InternalError("failed to collect created menu item", err)
	}

	for _, recipeItem := range params.Recipe {
		recipeBuilder := r.psql.
			Insert("recipe_ingredients").
			Columns("menu_item_id", "ingredient_id", "amount").
			Values(dbItem.ID, recipeItem.IngredientID, recipeItem.Amount)

		recipeQuery, recipeArgs, err := recipeBuilder.ToSql()
		if err != nil {
			return nil, utils.InternalError("failed to build recipe item query", err)
		}

		if _, err := r.provider.GetExecutor(ctx).Exec(ctx, recipeQuery, recipeArgs...); err != nil {
			return nil, utils.InternalError("failed to create recipe item", err)
		}
	}

	return r.GetMenuItem(ctx, dbItem.ID)
}
