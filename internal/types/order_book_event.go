package types

import (
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"log"
	"sort"
	"strings"
	"time"
)

const (
	Buy           = "BUY"
	Sale          = "SALE"
	MatchEngine   = "MATCH_ENGINE"
	InTrade       = "IN_TRADE"
	InNegotiation = "IN_NEGOTIATION"
	InOffer       = "IN_OFFER"
	Finished      = "FINISHED"
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

	MatchOrders map[string][]*DynamoEventMessage
)

func (eventMessage DynamoEventMessage) ToString() string {
	return fmt.Sprintf("%+v\n", eventMessage)
}

func (dynamoRecord DynamoRecord) GetTableName() (string, error) {
	sourceArn, err := arn.Parse(dynamoRecord.EventSourceArn)
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

func (keys DynamoEventMessageKey) GetKey() map[string]*dynamodb.AttributeValue {
	return map[string]*dynamodb.AttributeValue{
		"id":   {S: aws.String(keys.Range)},
		"type": {S: aws.String(keys.Hash)},
	}
}

func (eventMessage DynamoEventMessage) GetKey() DynamoEventMessageKey {

	return DynamoEventMessageKey{
		eventMessage.Type,
		eventMessage.Id,
	}
}

type Messages []*DynamoEventMessage

func (messages Messages) SortByCreatedAt() []*DynamoEventMessage {
	sort.Slice(messages, func(i, j int) bool {
		return messages[i].Audit.CreatedAt.Before(messages[j].Audit.CreatedAt)
	})

	return messages
}

func (dynamoRecord DynamoRecord) ConverterEventRaw() (new *DynamoEventMessage, old *DynamoEventMessage, err error) {

	if dynamoRecord.Change.OldImage != nil {
		err = dynamodbattribute.UnmarshalMap(dynamoRecord.Change.OldImage, &old)
		if err != nil {
			log.Printf("Error on convert message %s", err)
			return
		}
	}

	if dynamoRecord.Change.NewImage != nil {
		err = dynamodbattribute.UnmarshalMap(dynamoRecord.Change.NewImage, &new)
		if err != nil {
			log.Printf("Error on convert message %s", err)
			return
		}
	}
	return
}

func (dynamoRecord DynamoRecord) Key() *DynamoEventMessageKey {
	var key DynamoEventMessageKey
	_ = dynamodbattribute.UnmarshalMap(dynamoRecord.Change.Keys, &key)
	return &key
}
