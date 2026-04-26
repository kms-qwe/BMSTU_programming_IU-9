package repos

import (
	"errors"

	service "coffee-shop-backend/app/internal/models"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgconn"
)

func NewPSQL() sq.StatementBuilderType {
	return sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
}

func ApplyOrderFilters(builder sq.SelectBuilder, filter service.OrderFilter) sq.SelectBuilder {
	if filter.EmployeeID != "" {
		builder = builder.Where(sq.Eq{"e.id": filter.EmployeeID})
	}
	if filter.From != nil {
		builder = builder.Where(sq.GtOrEq{"o.created_at": *filter.From})
	}
	if filter.To != nil {
		builder = builder.Where(sq.LtOrEq{"o.created_at": *filter.To})
	}

	return builder
}

func ApplyInventoryFilters(builder sq.SelectBuilder, filter service.InventoryOperationFilter) sq.SelectBuilder {
	if filter.EmployeeID != "" {
		builder = builder.Where(sq.Eq{"e.id": filter.EmployeeID})
	}
	if filter.IngredientID != "" {
		builder = builder.Where(sq.Eq{"i.id": filter.IngredientID})
	}
	if filter.OperationType != "" {
		builder = builder.Where(sq.Eq{"io.operation_type": filter.OperationType})
	}
	if filter.From != nil {
		builder = builder.Where(sq.GtOrEq{"io.created_at": *filter.From})
	}
	if filter.To != nil {
		builder = builder.Where(sq.LtOrEq{"io.created_at": *filter.To})
	}

	return builder
}

func IsUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505"
	}
	return false
}
