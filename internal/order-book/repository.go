package order_book

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/caarlos0/env"
	"github.com/pkg/errors"
	"log"
	"order-book-match-engine/internal/types"
	"time"
)

type (
	DynamoDBAPI interface {
		TransactWriteItemsWithContext(ctx aws.Context, input *dynamodb.TransactWriteItemsInput, opts ...request.Option) (*dynamodb.TransactWriteItemsOutput, error)
		QueryWithContext(aws.Context, *dynamodb.QueryInput, ...request.Option) (*dynamodb.QueryOutput, error)
	}

	Config struct {
		TableName string `env:"ORDER_BOOK_TABLE_NAME,required"`
	}

	OperationRepository interface {
		FindAll(ctx context.Context, orderType string, status string) (types.Messages, error)
		Update(ctx context.Context, matchOrders types.MatchOrders, operation *types.DynamoEventMessage, status string) error
	}

	operationRepository struct {
		cfg *Config
		db  DynamoDBAPI
	}
)

func (r operationRepository) Update(ctx context.Context, matchOrders types.MatchOrders, operation *types.DynamoEventMessage, status string) error {
	now, err := dynamodbattribute.Marshal(time.Now())
	if err != nil {
		return errors.Wrap(err, "Marshal TimeNow to AttributeValue")
	}

	updateExpression := aws.String("SET operationStatus = :updatedStatus , audit.updatedAt = :now, audit.updatedBy = :updatedBy")
	expressionAttributeValues := map[string]*dynamodb.AttributeValue{
		":updatedStatus": {S: aws.String(status)},
		":updatedBy":     {S: aws.String(types.MatchEngine)},
		":now":           now,
	}

	var transactions []*dynamodb.TransactWriteItem
	transactions = append(transactions, &dynamodb.TransactWriteItem{
		Update: &dynamodb.Update{
			Key:                       operation.GetKey(),
			ConditionExpression:       nil,
			ExpressionAttributeValues: expressionAttributeValues,
			TableName:                 aws.String(r.cfg.TableName),
			UpdateExpression:          updateExpression,
		},
	})

	for _, match := range matchOrders[operation.Id] {
		transactions = append(transactions, &dynamodb.TransactWriteItem{
			Update: &dynamodb.Update{
				ConditionExpression:       nil,
				TableName:                 aws.String(r.cfg.TableName),
				UpdateExpression:          updateExpression,
				ExpressionAttributeValues: expressionAttributeValues,
				Key:                       match.GetKey(),
			},
		})
	}

	if _, err = r.db.TransactWriteItemsWithContext(ctx, &dynamodb.TransactWriteItemsInput{
		TransactItems: transactions,
	}); err != nil {
		return err
	}

	return nil
}

func (r operationRepository) FindAll(ctx context.Context, orderType string, status string) (types.Messages, error) {

	query := &dynamodb.QueryInput{
		KeyConditions: map[string]*dynamodb.Condition{
			"type": {
				AttributeValueList: []*dynamodb.AttributeValue{{S: aws.String(orderType)}},
				ComparisonOperator: aws.String(dynamodb.ComparisonOperatorEq),
			},
			"id": {
				AttributeValueList: []*dynamodb.AttributeValue{{S: aws.String(status)}},
				ComparisonOperator: aws.String(dynamodb.ComparisonOperatorBeginsWith),
			},
		},
		TableName: aws.String(r.cfg.TableName),
	}

	output, err := r.db.QueryWithContext(ctx, query)
	if err != nil {
		return nil, errors.Wrapf(err, "Get Key on Table(%s) By Id(%s) And Status(%s)", r.cfg.TableName, orderType, status)
	}

	var operations []*types.DynamoEventMessage
	err = dynamodbattribute.UnmarshalListOfMaps(output.Items, &operations)
	if err != nil {
		fmt.Printf("[ERROR] - Error when unmarshal operations")
	}

	return operations, nil
}

func NewOperationRepository(client DynamoDBAPI, config *Config) OperationRepository {

	if config == nil {
		config = new(Config)
		if err := env.Parse(config); err != nil {
			log.Fatalf("[ERROR] Missing configuration for operation repository: %s", err)
		}
	}

	return &operationRepository{
		cfg: config,
		db:  client,
	}
}
