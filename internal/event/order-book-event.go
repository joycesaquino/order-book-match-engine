package event

import (
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"log"
	"strings"
	"time"
)

type (
	DynamoEventMessageKey struct {
		Hash  string `dynamodbav:"type"`
		Range string `dynamodbav:"id"`
	}

	DynamoEventMessage struct {
		Id        string  `dynamodbav:"id"`
		RequestId string  `dynamodbav:"requestId"`
		Hash      string  `dynamodbav:"hash"`
		Value     float64 `dynamodbav:"value"`
		Quantity  int     `dynamodbav:"quantity"`
		Status    string  `dynamodbav:"status"`
		Type      string  `dynamodbav:"type"`
		UserId    int     `dynamodbav:"userId"`
		Audit     Audit   `dynamodbav:"audit"`
	}

	Audit struct {
		CreatedAt time.Time `dynamodbav:":createdAt,unixtime"`
		UpdatedAt time.Time `dynamodbav:":updatedAt,unixtime"`
		UpdatedBy string    `dynamodbav:"updatedBy"`
	}

	DynamoEvent struct {
		Records DynamoRecords `json:"Records"`
	}

	DynamoRecords []*DynamoRecord
	DynamoRecord  struct {
		EventSourceArn string `json:"eventSourceARN"`
		EventName      string `json:"eventName"`
		Change         Change `json:"dynamodb"`
	}
	Change struct {
		ApproximateCreationDateTime events.SecondsEpochTime             `json:"ApproximateCreationDateTime,omitempty"`
		Keys                        map[string]*dynamodb.AttributeValue `json:"Keys,omitempty"`
		NewImage                    map[string]*dynamodb.AttributeValue `json:"NewImage,omitempty"`
		OldImage                    map[string]*dynamodb.AttributeValue `json:"OldImage,omitempty"`
	}
)

func (dem DynamoEventMessage) ToString() string {
	return fmt.Sprintf("%+v\n", dem)
}

func (dr DynamoRecord) GetTableName() (string, error) {
	sourceArn, err := arn.Parse(dr.EventSourceArn)
	if err != nil {
		return "", err
	}
	if sourceArn.Service != dynamodb.ServiceName {
		return "", fmt.Errorf("SourceArn Service is not DynamoDB: %s", sourceArn.Service)
	}
	resource := sourceArn.Resource
	tableName := strings.TrimPrefix(resource, "table/")
	index := strings.Index(tableName, "/stream")
	if index > -1 {
		return strings.Split(tableName, "/stream")[0], nil
	}
	return tableName, nil
}

func (dem *DynamoEventMessage) GetKey() map[string]string {
	key := map[string]string{
		"id":   dem.Id,
		"type": dem.Type,
	}
	return key
}

func (dr DynamoRecord) ConverterEventRaw() (new *DynamoEventMessage, old *DynamoEventMessage, err error) {

	if dr.Change.OldImage != nil {
		err = dynamodbattribute.UnmarshalMap(dr.Change.OldImage, &old)
		if err != nil {
			log.Printf("Error on convert message %s", err)
			return
		}
	}

	if dr.Change.NewImage != nil {
		err = dynamodbattribute.UnmarshalMap(dr.Change.NewImage, &new)
		if err != nil {
			log.Printf("Error on convert message %s", err)
			return
		}
	}
	return
}

func (dr DynamoRecord) Key() *DynamoEventMessageKey {
	var key DynamoEventMessageKey
	_ = dynamodbattribute.UnmarshalMap(dr.Change.Keys, &key)
	return &key
}
