package order_book

import (
	"context"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/joycesaquino/order-book-match-engine/internal/types"
	"os"
	"reflect"
	"testing"
)

func mockRepository() OperationRepository {
	_ = os.Setenv("ORDER_BOOK_TABLE_NAME", "order-book-operation")

	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String("sa-east-1"),
	})
	db := dynamodb.New(sess)
	return NewOperationRepository(db, nil)
}

func Test_operationRepository_FindAll(t *testing.T) {
	type fields struct {
		repository OperationRepository
	}
	type args struct {
		ctx       context.Context
		orderType string
		status    string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    types.Messages
		wantErr bool
	}{
		{name: "Should get all sales in database", fields: fields{repository: mockRepository()}, args: args{
			ctx:       context.Background(),
			orderType: types.Sale,
			status:    types.InTrade,
		}, want: []*types.DynamoEventMessage{
			{
				Id:              "1234221|5267039e-7b48-4320-8bed-f67c8e9a376e",
				RequestId:       "5267039e-7b48-4320-8bed-f67c8e9a376e",
				Hash:            "3355622bd33f4f0e44057e5b6b2b433",
				Value:           10,
				Quantity:        100,
				OperationStatus: "IN_TRADE",
				Type:            "SALE",
				UserId:          0,
				Audit:           types.Audit{},
			},
		}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := tt.fields.repository
			got, err := r.FindAll(tt.args.ctx, tt.args.orderType, tt.args.status)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindAll() got = %v, want %v", len(got), tt.want)
			}
		})
	}
}
