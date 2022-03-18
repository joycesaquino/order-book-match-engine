package order_book

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/pkg/errors"
	"order-book-match-engine/internal/event"
)

type (
	DynamoDBAPI interface {
		PutItemWithContext(aws.Context, *dynamodb.PutItemInput, ...request.Option) (*dynamodb.PutItemOutput, error)
		BatchGetItemWithContext(aws.Context, *dynamodb.BatchGetItemInput, ...request.Option) (*dynamodb.BatchGetItemOutput, error)
		BatchWriteItemWithContext(aws.Context, *dynamodb.BatchWriteItemInput, ...request.Option) (*dynamodb.BatchWriteItemOutput, error)
		QueryWithContext(aws.Context, *dynamodb.QueryInput, ...request.Option) (*dynamodb.QueryOutput, error)
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

func (r oderRepository) FindAll(ctx context.Context, keys event.DynamoEventMessageKey, status string) ([]* event.DynamoEventMessage, error) {

	query := &dynamodb.QueryInput{
		KeyConditions: map[string]*dynamodb.Condition{
			"type": {
				AttributeValueList: []*dynamodb.AttributeValue{{S: aws.String(keys.Hash)}},
				ComparisonOperator: aws.String(dynamodb.ComparisonOperatorEq),
			},
			"id": {
				AttributeValueList: []*dynamodb.AttributeValue{{S: aws.String(status)}},
				ComparisonOperator: aws.String(dynamodb.ComparisonOperatorContains),
			},
		},
		TableName: aws.String(r.cfg.TableName),
	}

	output, err := r.db.QueryWithContext(ctx, query)
	if err != nil {
		return nil, errors.Wrapf(err, "Get Key on Table(%s) By Id(%s) And Status(%s)", r.cfg.TableName, keys.Hash, status)
	}

	var operations []*event.DynamoEventMessage
	err = dynamodbattribute.UnmarshalListOfMaps(output.Items, &operations)
	if err != nil {
		fmt.Printf("[ERROR] - Error when unmarshal operations")
	}

	return operations, nil
}
