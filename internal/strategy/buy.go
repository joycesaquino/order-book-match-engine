package strategy

import (
	"context"
	orderBook "order-book-match-engine/internal/order-book"
	"order-book-match-engine/internal/types"
	walletIntegration "order-book-match-engine/internal/wallet-integration"
)

type Buy struct {
	repository orderBook.OperationRepository
	queue      *walletIntegration.Queue
}

func (strategy Buy) Accept(input *Input) bool {
	return input.NewImage.Type == types.Buy
}

func (strategy Buy) Apply(ctx context.Context, input *Input) {
	buy := input.NewImage
	sales, err := strategy.repository.FindAll(ctx, types.Sale, input.NewImage.Status)
	if err != nil {
		return
	}

	matchOrders, itsAMatch := Match(ctx, buy, sales)

	if itsAMatch {
		err := strategy.repository.Update(ctx, buy.GetKey(), types.Finished)
		if err != nil {
			// TODO Return error to failed lambda and retry
			return
		}

		for _, sale := range matchOrders[buy.Id] {
			if err := strategy.repository.Update(ctx, sale.GetKey(), types.Finished); err != nil {
				// TODO Rollback transactions and return error to failed lambda and retry
				return
			}
		}

		if err := strategy.queue.Send(ctx, BuildOrders(buy, sales)); err != nil {
			return
		}
	}

	if !itsAMatch {
		if err := strategy.repository.Update(ctx, buy.GetKey(), types.InTrade); err != nil {
			return
		}

		if strategy.updateMatchOrdersStatus(ctx, matchOrders, buy.Id, types.InOffer) {
			return
		}
	}

}

// TODO Update sales using transaction and return error
func (strategy Buy) updateMatchOrdersStatus(ctx context.Context, matchOrders types.MatchOrders, buyId string, status string) bool {
	for _, sale := range matchOrders[buyId] {
		err := strategy.repository.Update(ctx, sale.GetKey(), status)
		if err != nil {
			// TODO Rollback transactions and return error to failed lambda and retry
			return true
		}
	}
	return false
}

func Match(ctx context.Context, buy *types.DynamoEventMessage, sales types.Messages) (map[string][]*types.DynamoEventMessage, bool) {

	var mod = buy.Quantity
	var value int

	matchOrders := make(map[string][]*types.DynamoEventMessage)
	for _, sale := range sales.SortByCreatedAt() {

		if buy.Quantity >= sale.Quantity {

			if (mod - sale.Quantity) < 0 {
				mod = mod + sale.Quantity
				continue
			}

			mod = mod - sale.Quantity
			if mod == 0 {

				//if err := strategy.repository.Update(ctx, sale.GetKey(), types.InNegotiation); err != nil {
				//	continue
				//}
				matchOrders[buy.Id] = append(matchOrders[buy.Id], sale)
				value = value + sale.Quantity
				return matchOrders, buy.Quantity == value
			}

			//if err := strategy.repository.Update(ctx, sale.GetKey(), types.InNegotiation); err != nil {
			//	continue
			//}
			matchOrders[buy.Id] = append(matchOrders[buy.Id], sale)
			value = value + sale.Quantity
			continue
		}

	}

	return matchOrders, buy.Quantity == value
}

func BuildOrders(buy *types.DynamoEventMessage, sales []*types.DynamoEventMessage) (orders []*types.Order) {

	orders = append(orders,
		&types.Order{
			Value:         buy.Value,
			Quantity:      buy.Quantity,
			OperationType: types.Buy,
			UserId:        buy.UserId,
			RequestId:     buy.RequestId,
		})

	for _, sale := range sales {
		orders = append(orders, &types.Order{
			Value:         sale.Value,
			Quantity:      sale.Quantity,
			OperationType: types.Sale,
			UserId:        sale.UserId,
			RequestId:     sale.RequestId,
		})
	}

	return orders

}
