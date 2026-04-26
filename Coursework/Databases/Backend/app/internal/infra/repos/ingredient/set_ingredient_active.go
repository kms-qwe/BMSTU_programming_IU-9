package ingredient

import (
	"context"
	"errors"

	"coffee-shop-backend/app/internal/models"
	"coffee-shop-backend/app/internal/utils"

	"github.com/jackc/pgx/v5"
)

func (r *Repo) SetIngredientActive(ctx context.Context, ingredientID string, isActive bool) (*models.Ingredient, error) {
	builder := r.psql.
		Update(ingredientsTable).
		Set("is_deleted", !isActive).
		Where("id = ?", ingredientID).
		Suffix("RETURNING id")

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, utils.InternalError("failed to build update ingredient activity query", err)
	}

	rows, err := r.provider.GetExecutor(ctx).Query(ctx, query, args...)
	if err != nil {
		return nil, utils.InternalError("failed to set ingredient activity", err)
	}
	defer rows.Close()

	dbItem, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[createdIngredientDBModel])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, utils.InternalError("failed to collect updated ingredient", err)
	}

	return r.GetIngredient(ctx, dbItem.ID)
}
