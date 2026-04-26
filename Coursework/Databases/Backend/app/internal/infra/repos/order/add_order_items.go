package order

import (
	"context"

	"coffee-shop-backend/app/internal/models"
	"coffee-shop-backend/app/internal/utils"
)

func (r *Repo) AddOrderItems(ctx context.Context, orderID string, items []models.OrderItemCreateRecord) error {
	dbItems := toDBOrderItemCreateModels(orderID, items)

	for _, item := range dbItems {
		builder := r.psql.
			Insert("order_items").
			Columns("order_id", "menu_item_id", "quantity", "price").
			Values(item.OrderID, item.MenuItemID, item.Quantity, item.Price)

		query, args, err := builder.ToSql()
		if err != nil {
			return utils.InternalError("failed to build add order item query", err)
		}

		if _, err := r.provider.GetExecutor(ctx).Exec(ctx, query, args...); err != nil {
			return utils.InternalError("failed to insert order item", err)
		}
	}

	return nil
}
