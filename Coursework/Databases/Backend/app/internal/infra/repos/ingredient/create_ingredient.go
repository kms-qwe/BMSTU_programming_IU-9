package ingredient

import (
	"context"

	baserepos "coffee-shop-backend/app/internal/infra/repos"
	"coffee-shop-backend/app/internal/models"
	"coffee-shop-backend/app/internal/utils"

	"github.com/jackc/pgx/v5"
)

func (r *Repo) CreateIngredient(ctx context.Context, params models.CreateIngredientParams) (*models.Ingredient, error) {
	builder := r.psql.
		Insert(ingredientsTable).
		Columns("name", "unit", "current_stock").
		Values(params.Name, params.Unit, params.InitialStock).
		Suffix("RETURNING id")

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, utils.InternalError("failed to build create ingredient query", err)
	}

	rows, err := r.provider.GetExecutor(ctx).Query(ctx, query, args...)
	if err != nil {
		return nil, utils.InternalError("failed to create ingredient", err)
	}
	defer rows.Close()

	dbItem, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[createdIngredientDBModel])
	if err != nil {
		if baserepos.IsUniqueViolation(err) {
			return nil, utils.ConflictError(utils.ErrorCodeDuplicateName, "Ингредиент с таким названием уже существует", nil)
		}
		return nil, utils.InternalError("failed to collect created ingredient", err)
	}

	return r.GetIngredient(ctx, dbItem.ID)
}
