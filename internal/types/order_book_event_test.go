package types

import (
	"reflect"
	"testing"
	"time"
)

func TestMessages_SortByCreatedAt(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name     string
		messages Messages
		want     []*DynamoEventMessage
	}{
		{name: "Should sort slice by createdAt", messages: []*DynamoEventMessage{
			{
				Audit: Audit{
					CreatedAt: now.Add(5),
				},
			},
			{
				Audit: Audit{
					CreatedAt: now,
				},
			},
			{
				Audit: Audit{
					CreatedAt: now.Add(3),
				},
			},
		}, want: []*DynamoEventMessage{
			{
				Audit: Audit{
					CreatedAt: now,
				},
			},
			{
				Audit: Audit{
					CreatedAt: now.Add(3),
				},
			},
			{
				Audit: Audit{
					CreatedAt: now.Add(5),
				},
			},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.messages.SortByCreatedAt(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SortByCreatedAt() = %v, want %v", got, tt.want)
			}
		})
	}
}
