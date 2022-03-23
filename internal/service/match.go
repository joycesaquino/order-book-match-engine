package service

import (
	"context"
	"github.com/aws/aws-sdk-go/aws/session"
	"order-book-match-engine/internal/strategy"
	"order-book-match-engine/internal/types"
)

type Match struct {
	validation *strategy.Validation
}

func NewMatchEngine(sess *session.Session) *Match {
	return &Match{validation: strategy.New(sess)}
}

func (m Match) Match(ctx context.Context, record *types.DynamoRecord) {
	newImage, oldImage, err := record.ConverterEventRaw()
	if err != nil {
		return
	}

	tableName, err := record.GetTableName()
	if err != nil {
		return
	}

	input := &strategy.Input{
		NewImage:  newImage,
		OldImage:  oldImage,
		EventName: record.EventName,
		TableName: tableName,
	}
	m.validation.Strategy(ctx, input)
}
