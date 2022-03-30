package main

import (
	"context"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/joycesaquino/order-book-match-engine/internal/types"
	"testing"
)

func TestHandler(t *testing.T) {
	type args struct {
		ctx         context.Context
		dynamoEvent types.DynamoEvent
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "Should match engine with orders", args: args{
			ctx: context.Background(),
			dynamoEvent: types.DynamoEvent{
				Records: []*types.DynamoRecord{
					{
						EventName: "INSERT",
						Change: types.Change{NewImage: map[string]*dynamodb.AttributeValue{
							"operationType":   {S: aws.String("BUY")},
							"userId":          {N: aws.String("2233998")},
							"id":              {S: aws.String("2233998|e99fc025-eb52-410b-b379-3288a4e712a4")},
							"quantity":        {N: aws.String("100")},
							"value":           {N: aws.String("1334.99")},
							"hash":            {S: aws.String("a6c23d29bced9632f76ec807d763f5d0")},
							"operationStatus": {S: aws.String("IN_TRADE")},
							"audit": {M: map[string]*dynamodb.AttributeValue{
								"createdAt": {S: aws.String("2022-03-29T22:31:43.821Z")},
								"updatedAt": {S: aws.String("2022-03-29T22:31:43.821Z")},
								"updatedBy": {S: aws.String("API")},
							}},
							"traceId": {S: aws.String("e99fc025-eb52-410b-b379-3288a4e712a4")},
						}},
					},
				},
			},
		}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Handler(tt.args.ctx, tt.args.dynamoEvent); (err != nil) != tt.wantErr {
				t.Errorf("Handler() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
