package service

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/caarlos0/env"
	orderBook "github.com/joycesaquino/order-book-match-engine/internal/order-book"
	"github.com/joycesaquino/order-book-match-engine/internal/queue"
	"github.com/joycesaquino/order-book-match-engine/internal/types"
	"log"
)

type Config struct {
	WalletQueue string `env:"WALLET_INTEGRATION_QUEUE_URL" envDefault:"order-book-wallet-integration-queue"`
	Region      string `env:"AWS_REGION" envDefault:"sa-east-1"`
}

type Match struct {
	repository orderBook.OperationRepository
	queue      queue.Queue
	config     Config
}

func NewMatchEngine() *Match {

	var config Config
	if err := env.Parse(&config); err != nil {
		log.Fatalf("[ERROR] Missing configuration for Match engine: %s", err)
	}

	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String(config.Region),
	})

	db := dynamodb.New(sess)
	repository := orderBook.NewOperationRepository(db)

	return &Match{
		repository: repository,
		queue:      queue.NewSQSQueue(sess),
	}
}

func (m Match) Match(ctx context.Context, operation *types.DynamoEventMessage) error {

	orders, err := m.repository.FindAll(ctx, getOperationType(operation), types.InTrade)
	if err != nil {
		return fmt.Errorf("[ERROR] - error on trying find orders: %v", err)
	}

	if len(orders) == 0 {
		fmt.Printf("[INFO] - Cant find any orders to match process for order id %s", operation.Id)
		return nil
	}

	matchOrders, itsAMatch := match(operation, orders)
	if itsAMatch {
		if err := m.repository.UpdateAll(ctx, matchOrders, operation, types.Finished); err != nil {
			return fmt.Errorf("[ERROR] - error on update match orders: %v", err)
		}

		if err := m.queue.Send(ctx, buildOrders(operation, matchOrders)); err != nil {
			return fmt.Errorf("[ERROR] - error on send match orders to wallet integration: %v", err)
		}
	}

	if !itsAMatch {
		fmt.Printf("[INFO] - Cant find a valid match for order %s orders to match process", operation.Id)
		if err := m.repository.Update(ctx, operation, types.InTrade); err != nil {
			return nil
		}
	}
	for _, order := range matchOrders {
		fmt.Printf("Match with sucess for operation type %s id %s order match type %s, id %s", operation.Type, operation.Id, order.Type, order.Id)
	}

	return nil
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

func match(operation *types.DynamoEventMessage, orders types.Messages) ([]*types.DynamoEventMessage, bool) {

	var limit = operation.Quantity
	var value int

	var matchOrders []*types.DynamoEventMessage
	for _, order := range orders.SortByCreatedAt() {

		if (limit - order.Quantity) < 0 {
			continue
		}

		limit = limit - order.Quantity
		if limit == 0 {
			matchOrders = append(matchOrders, order)
			value = value + order.Quantity
			return matchOrders, operation.Quantity == value
		}

		matchOrders = append(matchOrders, order)
		value = value + order.Quantity
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
