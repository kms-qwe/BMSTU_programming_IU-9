package ingredient

import (
	"context"

	"coffee-shop-backend/app/internal/models"
	"coffee-shop-backend/app/internal/utils"

	"github.com/jackc/pgx/v5"
)

func (r *Repo) ListIngredientUsage(ctx context.Context, ingredientID string) ([]models.IngredientUsedInRecipe, error) {
	builder := r.psql.
		Select(
			"m.id AS menu_item_id",
			"m.name AS menu_item_name",
			"ri.amount AS amount",
		).
		From("recipe_ingredients ri").
		Join("menu_items m ON m.id = ri.menu_item_id").
		Where("ri.ingredient_id = ?", ingredientID).
		OrderBy("m.name ASC")

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, utils.InternalError("failed to build ingredient usage query", err)
	}

	rows, err := r.provider.GetExecutor(ctx).Query(ctx, query, args...)
	if err != nil {
		return nil, utils.InternalError("failed to list ingredient usage", err)
	}
	defer rows.Close()

	dbItems, err := pgx.CollectRows(rows, pgx.RowToStructByName[ingredientUsageDBModel])
	if err != nil {
		return nil, utils.InternalError("failed to collect ingredient usage", err)
	}

	return toServiceIngredientUsage(dbItems), nil
}
