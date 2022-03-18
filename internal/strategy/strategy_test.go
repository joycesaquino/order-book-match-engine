package strategy

import (
	"github.com/stretchr/testify/mock"
	"order-book-match-engine/internal/event"
	"testing"
)

const QUEUE = "https://error.sa-east-1.amazonaws.com/562821017172/offer-update-hml-error-flow"
const tableName = "offer-update-hml-prices"

type StrategyMock struct {
	mock     mock.Mock
	strategy Strategy
	name     string
}

func (m *StrategyMock) Accept(i *Input) bool {
	m.mock.Called(i)
	return m.strategy.Accept(i)
}

func (m *StrategyMock) Apply(i *Input) {
	m.mock.Called(i)
}

func TestValidation_Strategy(t *testing.T) {

	var buy, sell, updateStatus, deleteOperation StrategyMock
	var strategies = map[string]*StrategyMock{}

	strategies["buy"] = &buy
	strategies["sell"] = &sell
	strategies["updateStatus"] = &updateStatus
	strategies["deleteOperation"] = &deleteOperation

	mockObject := &Validation{
		Strategies: []Strategy{
			&buy,
			&sell,
			&updateStatus,
			&deleteOperation,
		},
	}

	tests := []struct {
		name       string
		strategies map[string]bool
		input      *Input
	}{
		{name: "Success strategy", strategies: map[string]bool{"updateOffer": true},
			input: &Input{
				NewImageInput: &event.DynamoEventMessage{},
				OldImageInput: nil,
				TableName:     tableName,
				StatusChange:  true,
			},
		},
		{name: "ChangedSku Strategy", strategies: map[string]bool{"skuChange": true}, input: &Input{
			NewImageInput: &event.DynamoEventMessage{},
			OldImageInput: &event.DynamoEventMessage{},
			TableName:     tableName,
			StatusChange:  true,
		},
		},
		{
			name: "Disable offer Strategy", strategies: map[string]bool{"disableOffer": true},
			input: &Input{
				NewImageInput: &event.DynamoEventMessage{},
				OldImageInput: nil,
				TableName:     tableName,
				StatusChange:  true,
			},
		},
		{name: "Delete offer Strategy", strategies: map[string]bool{"deleteOffer": true},
			input: &Input{
				EventName:     "REMOVE",
				NewImageInput: nil,
				OldImageInput: &event.DynamoEventMessage{},
				TableName:     tableName,
				StatusChange:  true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			//buy = StrategyMock{mock: mock.Mock{}, strategy: &SuccessUpdate{}, name: "buy"}
			//sell = StrategyMock{mock: mock.Mock{}, strategy: &SkuHasChanged{}, name: "sell"}
			//updateStatus = StrategyMock{mock: mock.Mock{}, strategy: &DisableOffer{}, name: "updateStatus"}
			//deleteOperation = StrategyMock{mock: mock.Mock{}, strategy: &DeleteOffer{}, name: "deleteOperation"}

			for _, obj := range strategies {
				obj.mock.On("Apply", tt.input).Return()
				obj.mock.On("Accept", tt.input).Return()
			}

			//call method
			mockObject.Strategy(tt.input)

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
