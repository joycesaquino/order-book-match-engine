package strategy

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	orderBook "order-book-match-engine/internal/order-book"
	"order-book-match-engine/internal/types"
	walletIntegration "order-book-match-engine/internal/wallet-integration"
)

type Input struct {
	Key       types.DynamoEventMessageKey
	NewImage  *types.DynamoEventMessage
	OldImage  *types.DynamoEventMessage
	EventName string
	TableName string
}

func (i Input) toString() string {
	return fmt.Sprintf("%+v\n", i)
}

type Strategy interface {
	Accept(input *Input) bool
	Apply(ctx context.Context, input *Input)
}

type Validation struct {
	Strategies []Strategy
}

func (validation *Validation) Strategy(ctx context.Context, input *Input) {
	for _, strategy := range validation.Strategies {
		if strategy.Accept(input) {
			strategy.Apply(ctx, input)
		}
	}
}

func New(sess *session.Session) *Validation {
	db := dynamodb.New(sess)
	repository := orderBook.NewOperationRepository(db, nil)
	return &Validation{
		Strategies: []Strategy{
			&Buy{
				repository: repository,
				queue:      walletIntegration.NewQueue(sess),
			},
			&Sale{
				repository: repository,
				queue:      walletIntegration.NewQueue(sess),
			},
		},
	}
}
