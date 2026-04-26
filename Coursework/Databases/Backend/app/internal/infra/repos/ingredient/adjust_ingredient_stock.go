package ingredient

import (
	"context"

	sq "github.com/Masterminds/squirrel"

	"coffee-shop-backend/app/internal/utils"
)

func (r *Repo) AdjustIngredientStock(ctx context.Context, ingredientID string, delta float64) error {
	builder := r.psql.
		Update(ingredientsTable).
		Set("current_stock", sq.Expr("current_stock + ?", delta)).
		Where("id = ?", ingredientID)

	query, args, err := builder.ToSql()
	if err != nil {
		return utils.InternalError("failed to build adjust ingredient stock query", err)
	}

	result, err := r.provider.GetExecutor(ctx).Exec(ctx, query, args...)
	if err != nil {
		return utils.InternalError("failed to adjust ingredient stock", err)
	}
	if result.RowsAffected() == 0 {
		return utils.NotFoundError("Ингредиент не найден")
	}

	return nil
}
