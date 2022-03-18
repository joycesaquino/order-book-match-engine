package strategy

import (
	"context"
	"github.com/stretchr/testify/mock"
	"order-book-match-engine/internal/event"
	"testing"
)

const tableName = "order-book-table-name"

type StrategyMock struct {
	mock     mock.Mock
	strategy Strategy
	name     string
}

func (m *StrategyMock) Accept(i *Input) bool {
	m.mock.Called(i)
	return m.strategy.Accept(i)
}

func (m *StrategyMock) Apply(ctx context.Context, i *Input) {
	m.mock.Called(i)
}

func TestValidation_Strategy(t *testing.T) {

	var buyStrategy, sell, updateStatus, deleteOperation StrategyMock
	var strategies = map[string]*StrategyMock{}

	strategies["buy"] = &buyStrategy
	strategies["sell"] = &sell
	strategies["updateStatus"] = &updateStatus
	strategies["deleteOperation"] = &deleteOperation

	mockObject := &Validation{
		Strategies: []Strategy{
			&buyStrategy,
			//&sell,
			//&updateStatus,
			//&deleteOperation,
		},
	}

	tests := []struct {
		name       string
		strategies map[string]bool
		input      *Input
	}{
		{name: "Buy strategy", strategies: map[string]bool{"buyStrategy": true},
			input: &Input{
				NewImageInput: &event.DynamoEventMessage{},
				OldImageInput: nil,
				TableName:     tableName,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			buyStrategy = StrategyMock{mock: mock.Mock{}, strategy: &Buy{}, name: "buy"}
			//sell = StrategyMock{mock: mock.Mock{}, strategy: &SkuHasChanged{}, name: "sell"}
			//updateStatus = StrategyMock{mock: mock.Mock{}, strategy: &DisableOffer{}, name: "updateStatus"}
			//deleteOperation = StrategyMock{mock: mock.Mock{}, strategy: &DeleteOffer{}, name: "deleteOperation"}

			for _, obj := range strategies {
				obj.mock.On("Apply", tt.input).Return()
				obj.mock.On("Accept", tt.input).Return()
			}

			//call method
			mockObject.Strategy(context.Background(), tt.input)

			// verification
			for s, mockObj := range strategies {
				if tt.strategies[s] {
					mockObj.mock.AssertNumberOfCalls(t, "Apply", 1)
				} else {
					mockObj.mock.AssertNotCalled(t, "Apply")
				}
			}
		})
	}
}
