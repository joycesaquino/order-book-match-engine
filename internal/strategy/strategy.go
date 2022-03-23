package strategy

import (
	"context"
	"fmt"
	"order-book-match-engine/internal/types"
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

func New() *Validation {

	return &Validation{
		Strategies: []Strategy{},
	}
}