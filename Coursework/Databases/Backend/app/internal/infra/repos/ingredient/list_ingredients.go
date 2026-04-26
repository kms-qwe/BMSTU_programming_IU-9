package ingredient

import (
	"context"

	"coffee-shop-backend/app/internal/models"
	"coffee-shop-backend/app/internal/utils"

	"github.com/jackc/pgx/v5"
)

func (r *Repo) ListIngredients(ctx context.Context, includeInactive bool) ([]models.Ingredient, error) {
	builder := r.psql.
		Select(
			"i.id AS id",
			"i.name AS name",
			"i.unit AS unit",
			"i.current_stock AS current_stock",
			"NOT i.is_deleted AS is_active",
			"COUNT(ri.id) AS used_in_recipes_count",
		).
		From(ingredientsTable+" i").
		LeftJoin("recipe_ingredients ri ON ri.ingredient_id = i.id").
		Where("(? OR NOT i.is_deleted)", includeInactive).
		GroupBy("i.id", "i.name", "i.unit", "i.current_stock", "i.is_deleted").
		OrderBy("i.name ASC")

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, utils.InternalError("failed to build ingredients query", err)
	}

	rows, err := r.provider.GetExecutor(ctx).Query(ctx, query, args...)
	if err != nil {
		return nil, utils.InternalError("failed to list ingredients", err)
	}
	defer rows.Close()

	dbItems, err := pgx.CollectRows(rows, pgx.RowToStructByName[ingredientDBModel])
	if err != nil {
		return nil, utils.InternalError("failed to collect ingredients", err)
	}

	return toServiceIngredients(dbItems), nil
}
