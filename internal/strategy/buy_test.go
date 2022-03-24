package strategy

import (
	"order-book-match-engine/internal/types"
	"reflect"
	"testing"
)

func TestMatch(t *testing.T) {
	type args struct {
		buy   *types.DynamoEventMessage
		sales types.Messages
	}
	tests := []struct {
		name string
		args args
		want map[string][]*types.DynamoEventMessage
	}{
		{name: "Should Match orders considering partial match", args: args{
			buy: &types.DynamoEventMessage{
				Id:       "001",
				Quantity: 10,
			},
			sales: []*types.DynamoEventMessage{
				{Quantity: 5},
				{Quantity: 2},
				{Quantity: 6},
				{Quantity: 3},
			},
		}, want: map[string][]*types.DynamoEventMessage{
			"001": {{Quantity: 5}, {Quantity: 2}, {Quantity: 3}},
		}},
		{name: "Should Match orders considering full match", args: args{
			buy: &types.DynamoEventMessage{
				Id:       "001",
				Quantity: 10,
			},
			sales: []*types.DynamoEventMessage{
				{Quantity: 10},
				{Quantity: 2},
				{Quantity: 6},
				{Quantity: 3},
			},
		}, want: map[string][]*types.DynamoEventMessage{
			"001": {{Quantity: 10}},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Match(tt.args.buy, tt.args.sales); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Match() = %v, want %v", got, tt.want)
			}

		})
	}
}
