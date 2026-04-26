package order

import (
	baserepos "coffee-shop-backend/app/internal/infra/repos"
	"coffee-shop-backend/app/internal/utils"

	sq "github.com/Masterminds/squirrel"
)

const ordersTable = "orders"

type Repo struct {
	provider utils.IExecutorProvider
	psql     sq.StatementBuilderType
}

func New(provider utils.IExecutorProvider) *Repo {
	return &Repo{provider: provider, psql: baserepos.NewPSQL()}
}
