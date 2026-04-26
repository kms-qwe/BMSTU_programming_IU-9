package ingredient

import (
	"context"
	"errors"

	"coffee-shop-backend/app/internal/models"
	"coffee-shop-backend/app/internal/utils"

	"github.com/jackc/pgx/v5"
)

func (r *Repo) GetIngredient(ctx context.Context, ingredientID string) (*models.Ingredient, error) {
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
		Where("i.id = ?", ingredientID).
		GroupBy("i.id", "i.name", "i.unit", "i.current_stock", "i.is_deleted")

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, utils.InternalError("failed to build ingredient query", err)
	}

	rows, err := r.provider.GetExecutor(ctx).Query(ctx, query, args...)
	if err != nil {
		return nil, utils.InternalError("failed to get ingredient", err)
	}
	defer rows.Close()

	dbItem, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[ingredientDBModel])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, utils.InternalError("failed to collect ingredient", err)
	}

	serviceItem := toServiceIngredient(dbItem)
	return &serviceItem, nil
}
