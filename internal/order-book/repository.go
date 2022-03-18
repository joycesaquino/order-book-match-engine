package order_book

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"order-book-match-engine/internal/event"
)

type (
	DynamoDBAPI interface {
		PutItemWithContext(aws.Context, *dynamodb.PutItemInput, ...request.Option) (*dynamodb.PutItemOutput, error)
		BatchGetItemWithContext(aws.Context, *dynamodb.BatchGetItemInput, ...request.Option) (*dynamodb.BatchGetItemOutput, error)
		BatchWriteItemWithContext(aws.Context, *dynamodb.BatchWriteItemInput, ...request.Option) (*dynamodb.BatchWriteItemOutput, error)
	}

	Config struct {
		TableName string `env:"ORDER_BOOK_TABLE_NAME,required"`
	}

	OrderRepository interface {
		FindAll(ctx context.Context, keys event.DynamoEventMessageKey) ([]event.DynamoEventMessage, error)
	}

	oderRepository struct {
		cfg *Config
		db  DynamoDBAPI
	}
)

type Match struct {

}

func (r oderRepository) FindAll(ctx context.Context, keys event.DynamoEventMessageKey) ([]event.DynamoEventMessage, error) {
	var batch = new(dynamodb.BatchGetItemInput)

	batch.SetRequestItems(map[string]*dynamodb.KeysAndAttributes{
		r.cfg.TableName: {
			ConsistentRead: aws.Bool(true),
			Keys: []map[string]*dynamodb.AttributeValue{},
		},
	})

	_, err := r.db.BatchGetItemWithContext(ctx, batch)
	if err != nil {
		return nil, fmt.Errorf("an error occurred while searching prices %v at table %s: %w", keys, r.cfg.TableName, err)
	}

	return nil,nil
}