package ingredient

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"

	"coffee-shop-backend/app/internal/models"
	"coffee-shop-backend/app/internal/utils"
)

func (r *Repo) GetIngredientsByIDs(ctx context.Context, ingredientIDs []string, forUpdate bool) ([]models.Ingredient, error) {
	if len(ingredientIDs) == 0 {
		return nil, nil
	}

	builder := r.psql.
		Select(
			"i.id AS id",
			"i.name AS name",
			"i.unit AS unit",
			"i.current_stock AS current_stock",
			"NOT i.is_deleted AS is_active",
			"COALESCE(rc.used_in_recipes_count, 0) AS used_in_recipes_count",
		).
		From(ingredientsTable + " i").
		LeftJoin(`(
			SELECT ingredient_id, COUNT(*) AS used_in_recipes_count
			FROM recipe_ingredients
			GROUP BY ingredient_id
		) rc ON rc.ingredient_id = i.id`).
		Where(sq.Eq{"i.id": ingredientIDs}).
		OrderBy("i.name ASC")

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, utils.InternalError("failed to build ingredients by ids query", err)
	}
	if forUpdate {
		query += " FOR UPDATE OF i"
	}

	rows, err := r.provider.GetExecutor(ctx).Query(ctx, query, args...)
	if err != nil {
		return nil, utils.InternalError("failed to get ingredients by ids", err)
	}
	defer rows.Close()

	dbItems, err := pgx.CollectRows(rows, pgx.RowToStructByName[ingredientDBModel])
	if err != nil {
		return nil, utils.InternalError("failed to collect ingredients by ids", err)
	}

	return toServiceIngredients(dbItems), nil
}
