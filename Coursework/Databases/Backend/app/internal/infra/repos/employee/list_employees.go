package employee

import (
	"context"

	"coffee-shop-backend/app/internal/models"
	"coffee-shop-backend/app/internal/utils"

	"github.com/jackc/pgx/v5"
)

func (r *Repo) ListEmployees(ctx context.Context) ([]models.Employee, error) {
	builder := r.psql.
		Select(
			"id AS id",
			"full_name AS full_name",
			"login AS login",
			"phone AS phone",
		).
		From(employeesTable).
		OrderBy("full_name ASC")

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, utils.InternalError("failed to build employee list query", err)
	}

	rows, err := r.provider.GetExecutor(ctx).Query(ctx, query, args...)
	if err != nil {
		return nil, utils.InternalError("failed to list employees", err)
	}
	defer rows.Close()

	dbItems, err := pgx.CollectRows(rows, pgx.RowToStructByName[employeeDBModel])
	if err != nil {
		return nil, utils.InternalError("failed to collect employees", err)
	}

	return toServiceEmployees(dbItems), nil
}
