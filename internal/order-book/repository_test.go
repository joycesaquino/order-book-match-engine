package order_book

import (
	"context"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"order-book-match-engine/internal/types"
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

func Test_operationRepository_Update(t *testing.T) {
	type fields struct {
		repository OperationRepository
	}
	type args struct {
		ctx    context.Context
		keys   types.DynamoEventMessageKey
		status string
	}
	_ = []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{name: "Should update status", fields: fields{repository: mockRepository()}, args: args{
			ctx: context.Background(),
			keys: types.DynamoEventMessageKey{
				Hash:  "BUY",
				Range: "AVAILABLE|334455|daa42855-7786-4a27-a9c6-8917dd26e2a7",
			},
			status: types.InNegotiation,
		}, wantErr: false},
	}
	//for _, tt := range tests {
	//	t.Run(tt.name, func(t *testing.T) {
	//		r := tt.fields.repository
	//		if err := r.Update(tt.args.ctx, tt.args.keys, tt.args.status); (err != nil) != tt.wantErr {
	//			t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
	//		}
	//	})
	//}
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
			status:    types.InOffer,
		}, want: []*types.DynamoEventMessage{
			{
				Id:        "IN_OFFER|334455|028deebb-4007-4a9a-a3a9-2a1d9d5a7865",
				RequestId: "",
				Hash:      "",
				Value:     0,
				Quantity:  0,
				Status:    "",
				Type:      "",
				UserId:    0,
				Audit:     types.Audit{},
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
