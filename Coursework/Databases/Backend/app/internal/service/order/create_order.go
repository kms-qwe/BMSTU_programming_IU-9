package order

import (
	"context"
	"strings"

	"coffee-shop-backend/app/internal/models"
	"coffee-shop-backend/app/internal/service"
	"coffee-shop-backend/app/internal/utils"
)

func (s *Service) CreateOrder(ctx context.Context, params models.CreateOrderParams) (*models.Order, error) {
	if len(params.Items) == 0 {
		return nil, utils.ValidationError("order items are required", nil)
	}
	aggregated := map[string]int{}
	for _, item := range params.Items {
		if strings.TrimSpace(item.MenuItemID) == "" || item.Quantity <= 0 {
			return nil, utils.ValidationError("each item must contain menu_item_id and positive quantity", nil)
		}
		aggregated[item.MenuItemID] += item.Quantity
	}
	menuIDs := make([]string, 0, len(aggregated))
	for id := range aggregated {
		menuIDs = append(menuIDs, id)
	}
	var created *models.Order
	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		activeShift, err := s.shiftRepo.GetActiveShift(ctx, true)
		if err != nil {
			return err
		}
		if activeShift == nil {
			return utils.BusinessError(utils.ErrorCodeShiftNotActive, "Нельзя создать заказ без активной смены", nil)
		}
		menuItems, err := s.menuRepo.ListMenuItemsByIDs(ctx, menuIDs)
		if err != nil {
			return err
		}
		if len(menuItems) != len(menuIDs) {
			return utils.NotFoundError("Не все пункты меню найдены")
		}
		menuByID := map[string]models.MenuItem{}
		required := map[string]float64{}
		for _, item := range menuItems {
			if !item.IsActive {
				return utils.BusinessError(utils.ErrorCodeInactiveMenuItem, "Нельзя добавить в заказ неактивный пункт меню", nil)
			}
			menuByID[item.ID] = item
			for _, recipeItem := range item.Recipe {
				if !recipeItem.IngredientIsActive {
					return utils.BusinessError(utils.ErrorCodeInactiveIngredient, "Нельзя создать заказ с неактивным ингредиентом в рецепте", nil)
				}
				required[recipeItem.IngredientID] += recipeItem.Amount * float64(aggregated[item.ID])
			}
		}
		ingredientIDs := make([]string, 0, len(required))
		for id := range required {
			ingredientIDs = append(ingredientIDs, id)
		}
		ingredients, err := s.ingredientRepo.GetIngredientsByIDs(ctx, ingredientIDs, true)
		if err != nil {
			return err
		}
		ingredientByID := map[string]models.Ingredient{}
		for _, item := range ingredients {
			ingredientByID[item.ID] = item
		}
		problems := make([]map[string]any, 0)
		for id, amount := range required {
			ing, ok := ingredientByID[id]
			if !ok {
				return utils.NotFoundError("Не все ингредиенты найдены")
			}
			if ing.CurrentStock+1e-9 < amount {
				problems = append(problems, map[string]any{"ingredient_id": ing.ID, "ingredient_name": ing.Name, "required": amount, "available": ing.CurrentStock, "unit": ing.Unit})
			}
		}
		if len(problems) > 0 {
			return utils.BusinessError(utils.ErrorCodeInsufficientStock, "Недостаточно ингредиентов для создания заказа", map[string]any{"ingredients": problems})
		}
		created, err = s.orderRepo.CreateOrder(ctx, activeShift.ID)
		if err != nil {
			return err
		}
		orderItems := make([]models.OrderItemCreateRecord, 0, len(menuIDs))
		for _, id := range menuIDs {
			item := menuByID[id]
			orderItems = append(orderItems, models.OrderItemCreateRecord{MenuItemID: id, Quantity: aggregated[id], Price: item.Price})
		}
		if err := s.orderRepo.AddOrderItems(ctx, created.ID, orderItems); err != nil {
			return err
		}
		for id, amount := range required {
			if err := s.ingredientRepo.AdjustIngredientStock(ctx, id, -amount); err != nil {
				return err
			}
			orderID := created.ID
			if _, err := s.inventoryRepo.CreateInventoryOperation(ctx, models.InventoryOperationCreateRecord{ShiftID: activeShift.ID, IngredientID: id, OrderID: &orderID, ChangeAmount: -amount, OperationType: service.OperationTypeAutoOrderWriteOff}); err != nil {
				return err
			}
		}
		created, err = s.orderRepo.GetOrder(ctx, created.ID)
		return err
	})
	if err != nil {
		return nil, err
	}
	return created, nil
}
