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
	matchOrders := Match(buy, eventMessages)

	var value int
	for _, match := range matchOrders[buy.Id] {
		value = value + match.Quantity
	}

	if buy.Quantity == value {
		// Deu Match, update status de compra e venda
		// Update ok ? Manda pra fila
		// Update not ok com conditionalCheckFailed? Retry ou Update Date para estimular novamente o registro
	}

	if buy.Quantity != value {
		// Se não, update status IN_TRADE compra.
	}

}

func Match(buy *types.DynamoEventMessage, sales types.Messages) map[string][]*types.DynamoEventMessage {

	var mod = buy.Quantity
	matchOrders := make(map[string][]*types.DynamoEventMessage)
	for _, sale := range sales.SortByCreatedAt() {

		if buy.Quantity >= sale.Quantity {

			if (mod - sale.Quantity) < 0 {
				mod = mod + sale.Quantity
				continue
			}

			mod = mod - sale.Quantity
			if mod == 0 {
				matchOrders[buy.Id] = append(matchOrders[buy.Id], sale)
				return matchOrders
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
