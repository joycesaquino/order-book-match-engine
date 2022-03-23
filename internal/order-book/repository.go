package order_book

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
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
		UpdateItemWithContext(aws.Context, *dynamodb.UpdateItemInput, ...request.Option) (*dynamodb.UpdateItemOutput, error)
		QueryWithContext(aws.Context, *dynamodb.QueryInput, ...request.Option) (*dynamodb.QueryOutput, error)
	}

	Config struct {
		TableName string `env:"ORDER_BOOK_TABLE_NAME,required"`
	}

	OperationRepository interface {
		FindAll(ctx context.Context, keys types.DynamoEventMessageKey, status string) (types.Messages, error)
		Update(ctx context.Context, keys types.DynamoEventMessageKey, status string) error
	}

	operationRepository struct {
		cfg *Config
		db  DynamoDBAPI
	}
)

func (r operationRepository) Update(ctx context.Context, keys types.DynamoEventMessageKey, status string) error {
	now, err := dynamodbattribute.Marshal(time.Now())
	if err != nil {
		return errors.Wrap(err, "Marshal TimeNow to AttributeValue")
	}

	input := &dynamodb.UpdateItemInput{
		Key:                 keys.GetKey(),
		ReturnValues:        aws.String(dynamodb.ReturnValueNone),
		UpdateExpression:    aws.String("SET status = :status,audit.updatedAt = :now, audit.updatedBy = :updatedBy"),
		ConditionExpression: aws.String("status = :available"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":status":    {S: aws.String(status)},
			":available": {S: aws.String(status)},
			":updatedBy": {S: aws.String(types.MatchEngine)},
			":now":       now,
		},
	}

	_, err = r.db.UpdateItemWithContext(ctx, input)
	if err != nil {
		if err, ok := err.(awserr.Error); ok && err.Code() == dynamodb.ErrCodeConditionalCheckFailedException {
			log.Printf("[ConditionalCheck] - error : %s ", err)
			return nil
		}
		return errors.Wrapf(err, "Update Error")
	}
	return nil
}

func (r operationRepository) FindAll(ctx context.Context, keys types.DynamoEventMessageKey, status string) (types.Messages, error) {

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
