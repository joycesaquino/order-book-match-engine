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
		want  int
	}
	tests := []struct {
		name string
		args args
		want map[string]*types.DynamoEventMessage
	}{
		{name: "Should Match orders", args: args{
			buy: &types.DynamoEventMessage{
				Id:       "001",
				Quantity: 10,
			},
			sales: []*types.DynamoEventMessage{
				{
					Quantity: 5,
				},
				{
					Quantity: 2,
				},
				{
					Quantity: 6,
				},
				{
					Quantity: 3,
				},
			},
			want: 2,
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Match(tt.args.buy, tt.args.sales); !reflect.DeepEqual(len(got), tt.want) {
				t.Errorf("Match() = %v, want %v", len(got), tt.want)
			}

		})
	}
}
