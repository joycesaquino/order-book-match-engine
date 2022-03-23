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

	for _, sale := range eventMessages.SortByCreatedAt() {
		if buy.Quantity == sale.Quantity {
			err := strategy.queue.Send(ctx, BuildOrders(buy, sale))
			if err != nil {
				return
			}
		}
	}
}

func BuildOrders(buyMessage *types.DynamoEventMessage, saleMessage *types.DynamoEventMessage) (orders []*types.Order) {

	orders = append(orders,
		&types.Order{
			Value:         buyMessage.Value,
			Quantity:      buyMessage.Quantity,
			OperationType: types.Buy,
			UserId:        buyMessage.UserId,
			RequestId:     buyMessage.RequestId,
		}, &types.Order{
			Value:         saleMessage.Value,
			Quantity:      saleMessage.Quantity,
			OperationType: types.Sale,
			UserId:        saleMessage.UserId,
			RequestId:     saleMessage.RequestId,
		})

	return orders

}
