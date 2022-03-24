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

func (strategy *Buy) Accept(input *Input) bool {
	return input.NewImage.Type == types.Buy
}

func (strategy *Buy) Apply(ctx context.Context, input *Input) {
	buy := input.NewImage
	eventMessages, err := strategy.repository.FindAll(ctx, input.Key, input.NewImage.Status)
	if err != nil {
		return
	}
	_ = Match(buy, eventMessages)

}

func (strategy *Buy) FullMatch(ctx context.Context, buy *types.DynamoEventMessage, sale *types.DynamoEventMessage) {
	if err := strategy.repository.Update(ctx, sale.GetKey(), types.Unavailable); err != nil {
		return
	}

	if err := strategy.repository.Update(ctx, buy.GetKey(), types.Unavailable); err != nil {
		return
	}

	err := strategy.queue.Send(ctx, BuildOrders(buy, sale))
	if err != nil {
		return
	}
}

func Match(buy *types.DynamoEventMessage, sales types.Messages) map[string][]*types.DynamoEventMessage {

	var mod = buy.Quantity
	matchOrders := make(map[string][]*types.DynamoEventMessage)
	for _, sale := range sales.SortByCreatedAt() {

		if buy.Quantity > sale.Quantity {

			if mod-sale.Quantity < 0 {
				mod = mod + sale.Quantity
				continue
			}

			mod = mod - sale.Quantity
			if mod == 0 {
				matchOrders[buy.Id] = append(matchOrders[buy.Id], sale)
				break
			}

			matchOrders[buy.Id] = append(matchOrders[buy.Id], sale)
			continue
		}

	}

	return matchOrders
}

func BuildOrders(buy *types.DynamoEventMessage, sale *types.DynamoEventMessage) (orders []*types.Order) {

	orders = append(orders,
		&types.Order{
			Value:         buy.Value,
			Quantity:      buy.Quantity,
			OperationType: types.Buy,
			UserId:        buy.UserId,
			RequestId:     buy.RequestId,
		}, &types.Order{
			Value:         sale.Value,
			Quantity:      sale.Quantity,
			OperationType: types.Sale,
			UserId:        sale.UserId,
			RequestId:     sale.RequestId,
		})

	return orders

}
