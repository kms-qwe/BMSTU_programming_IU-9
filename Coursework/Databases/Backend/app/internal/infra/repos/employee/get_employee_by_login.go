package employee

import (
	"context"
	"errors"

	"coffee-shop-backend/app/internal/models"
	"coffee-shop-backend/app/internal/utils"

	"github.com/jackc/pgx/v5"
)

func (r *Repo) GetEmployeeByLogin(ctx context.Context, login string) (*models.EmployeeCredentials, error) {
	builder := r.psql.
		Select(
			"id AS id",
			"full_name AS full_name",
			"login AS login",
			"phone AS phone",
			"password_hash AS password_hash",
		).
		From(employeesTable).
		Where("login = ?", login)

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, utils.InternalError("failed to build employee query", err)
	}

	rows, err := r.provider.GetExecutor(ctx).Query(ctx, query, args...)
	if err != nil {
		return nil, utils.InternalError("failed to get employee", err)
	}
	defer rows.Close()

	dbItem, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[employeeCredentialsDBModel])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, utils.InternalError("failed to collect employee", err)
	}

	serviceItem := toServiceEmployeeCredentials(dbItem)
	return &serviceItem, nil
}
