package service

import (
	"context"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	orderBook "order-book-match-engine/internal/order-book"
	"order-book-match-engine/internal/types"
	walletIntegration "order-book-match-engine/internal/wallet-integration"
)

type Match struct {
	repository orderBook.OperationRepository
	queue      *walletIntegration.Queue
}

func NewMatchEngine(sess *session.Session) *Match {
	db := dynamodb.New(sess)
	repository := orderBook.NewOperationRepository(db, nil)

	return &Match{
		repository: repository,
		queue:      walletIntegration.NewQueue(sess),
	}
}

func (m Match) Match(ctx context.Context, record *types.DynamoRecord) {

	newImage, _, err := record.ConverterEventRaw()
	if err != nil {
		return
	}

	sales, err := m.repository.FindAll(ctx, getOperationType(newImage), newImage.Status)
	if err != nil {
		return
	}

	matchOrders, itsAMatch := match(newImage, sales)
	if itsAMatch {
		if err := m.repository.Update(ctx, matchOrders, newImage, types.Finished); err != nil {
			return
		}
		if err := m.queue.Send(ctx, buildOrders(newImage, sales)); err != nil {
			return
		}
	}

	if !itsAMatch {
		if err := m.repository.Update(ctx, matchOrders, newImage, types.InTrade); err != nil {
			return
		}
	}

}

func getOperationType(newImage *types.DynamoEventMessage) string {
	var operationType string

	switch newImage.Type {
	case types.Sale:
		operationType = types.Buy
		break
	case types.Buy:
		operationType = types.Sale
		break
	}
	return operationType
}

func match(operation *types.DynamoEventMessage, orders types.Messages) (map[string][]*types.DynamoEventMessage, bool) {

	var mod = operation.Quantity
	var value int

	matchOrders := make(map[string][]*types.DynamoEventMessage)
	for _, order := range orders.SortByCreatedAt() {

		if operation.Quantity >= order.Quantity {

			if (mod - order.Quantity) < 0 {
				mod = mod + order.Quantity
				continue
			}

			mod = mod - order.Quantity
			if mod == 0 {

				matchOrders[operation.Id] = append(matchOrders[operation.Id], order)
				value = value + order.Quantity
				return matchOrders, operation.Quantity == value
			}

			matchOrders[operation.Id] = append(matchOrders[operation.Id], order)
			value = value + order.Quantity
			continue
		}

	}

	return matchOrders, operation.Quantity == value
}

func buildOrders(buy *types.DynamoEventMessage, sales []*types.DynamoEventMessage) (orders []*types.Order) {

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
