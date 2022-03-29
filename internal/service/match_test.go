package service

import (
	"context"
	mockRepository "github.com/joycesaquino/order-book-match-engine/internal/order-book/mocks"
	mocksQueue "github.com/joycesaquino/order-book-match-engine/internal/queue/mocks"
	"github.com/joycesaquino/order-book-match-engine/internal/types"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestMatch_Match(t *testing.T) {

	buyOperation := &types.DynamoEventMessage{
		Id:       "001",
		Quantity: 10,
		Type:     types.Buy,
	}

	// test case 1
	var orders1 types.Messages

	orders1 = []*types.DynamoEventMessage{
		{Quantity: 5},
		{Quantity: 2},
		{Quantity: 6},
		{Quantity: 3},
	}

	var matchedOrders1 types.Messages

	matchedOrders1 = []*types.DynamoEventMessage{
		{Quantity: 5},
		{Quantity: 2},
		{Quantity: 3},
	}

	// test case 2

	var orders2 types.Messages

	orders2 = []*types.DynamoEventMessage{
		{Quantity: 10},
		{Quantity: 2},
		{Quantity: 6},
		{Quantity: 3},
	}

	var matchedOrders2 types.Messages

	matchedOrders2 = []*types.DynamoEventMessage{
		{Quantity: 10},
	}

	// test case 3

	var orders3 types.Messages

	orders3 = []*types.DynamoEventMessage{
		{Quantity: 6},
		{Quantity: 3},
	}

	var matchedOrders3 types.Messages

	matchedOrders3 = []*types.DynamoEventMessage{
		{Quantity: 6},
		{Quantity: 3},
	}

	type fields struct {
		repository *mockRepository.OperationRepository
		queue      *mocksQueue.Queue
	}
	type args struct {
		ctx       context.Context
		operation *types.DynamoEventMessage
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		mock    func(repository *mockRepository.OperationRepository, queue *mocksQueue.Queue)
	}{
		{
			name: "Should Match orders1 considering partial match",
			fields: fields{
				repository: new(mockRepository.OperationRepository),
				queue:      new(mocksQueue.Queue),
			},
			args: args{
				ctx:       context.TODO(),
				operation: buyOperation,
			},
			mock: func(repository *mockRepository.OperationRepository, queue *mocksQueue.Queue) {
				repository.
					On("FindAll", mock.Anything, types.Sale, types.InTrade).
					Return(orders1, nil).
					Times(1)

				repository.
					On("Update", mock.Anything, matchedOrders1, buyOperation, types.Finished).
					Return(nil).
					Times(1)

				queue.
					On("Send", mock.Anything, []*types.Order{{
						Quantity:      10,
						OperationType: types.Buy,
					}, {
						Quantity:      5,
						OperationType: types.Sale,
					}, {
						Quantity:      2,
						OperationType: types.Sale,
					}, {
						Quantity:      3,
						OperationType: types.Sale,
					}}).
					Return(nil).
					Times(1)
			},
		},
		{
			name: "Should Match orders considering full match",
			fields: fields{
				repository: new(mockRepository.OperationRepository),
				queue:      new(mocksQueue.Queue),
			},
			args: args{
				ctx:       context.TODO(),
				operation: buyOperation,
			},
			mock: func(repository *mockRepository.OperationRepository, queue *mocksQueue.Queue) {
				repository.
					On("FindAll", mock.Anything, types.Sale, types.InTrade).
					Return(orders2, nil).
					Times(1)

				repository.
					On("Update", mock.Anything, matchedOrders2, buyOperation, types.Finished).
					Return(nil).
					Times(1)

				queue.
					On("Send", mock.Anything, []*types.Order{{
						Quantity:      10,
						OperationType: types.Buy,
					}, {
						Quantity:      10,
						OperationType: types.Sale,
					}}).
					Return(nil).
					Times(1)
			},
		},
		{
			name: "Should NOT Match orders",
			fields: fields{
				repository: new(mockRepository.OperationRepository),
				queue:      new(mocksQueue.Queue),
			},
			args: args{
				ctx:       context.TODO(),
				operation: buyOperation,
			},
			mock: func(repository *mockRepository.OperationRepository, queue *mocksQueue.Queue) {
				repository.
					On("FindAll", mock.Anything, types.Sale, types.InTrade).
					Return(orders3, nil).
					Times(1)

				repository.
					On("Update", mock.Anything, matchedOrders3, buyOperation, types.InTrade).
					Return(nil).
					Times(1)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock(tt.fields.repository, tt.fields.queue)

			m := Match{
				repository: tt.fields.repository,
				queue:      tt.fields.queue,
			}

			m.Match(tt.args.ctx, tt.args.operation)

			tt.fields.repository.AssertExpectations(t)
			tt.fields.queue.AssertExpectations(t)
		})
	}
}
