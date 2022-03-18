package strategy

import (
	"fmt"
	"order-book-match-engine/internal/event"
)

type Input struct {
	NewImageInput *event.DynamoEventMessage
	OldImageInput *event.DynamoEventMessage
	EventName     string
	TableName     string
}

func (i Input) toString() string {
	return fmt.Sprintf("%+v\n", i)
}

type Strategy interface {
	Accept(input *Input) bool
	Apply(input *Input)
}

type Validation struct {
	Strategies []Strategy
}

func (validation *Validation) Strategy(input *Input) {
	for _, strategy := range validation.Strategies {
		if strategy.Accept(input) {
			strategy.Apply(input)
		}
	}
}

func New() *Validation {

	return &Validation{
		Strategies: []Strategy{},
	}
}
