package service

import (
	"context"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"order-book-match-engine/internal/event"
	orderBook "order-book-match-engine/internal/order-book"
	"order-book-match-engine/internal/strategy"
)

type Match struct {
	validation *strategy.Validation
	Order      orderBook.OperationRepository
}

func NewMatchEngine(sess *session.Session) *Match {
	db := dynamodb.New(sess)
	return &Match{Order: orderBook.NewOperationRepository(db, nil)}
}

func (m Match) Match(ctx context.Context, record *event.DynamoRecord) {
	newImage, oldImage, err := record.ConverterEventRaw()
	if err != nil {
		return
	}

	tableName, err := record.GetTableName()
	if err != nil {
		return
	}

	input := &strategy.Input{
		NewImageInput: newImage,
		OldImageInput: oldImage,
		EventName:     record.EventName,
		TableName:     tableName,
	}
	m.validation.Strategy(ctx, input)
}
