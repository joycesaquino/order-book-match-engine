package types

import (
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"log"
	"sort"
	"time"
)

const (
	Buy         = "BUY"
	Sale        = "SALE"
	MatchEngine = "MATCH_ENGINE"
	InTrade     = "IN_TRADE"
	Finished    = "FINISHED"
	REMOVE      = "REMOVE"
)

type (
	DynamoEventMessageKey struct {
		Hash  string `dynamodbav:"type"`
		Range string `dynamodbav:"id"`
	}

	DynamoEventMessage struct {
		Id              string  `dynamodbav:"id"`
		RequestId       string  `dynamodbav:"traceId"`
		Hash            string  `dynamodbav:"hash"`
		Value           float64 `dynamodbav:"value"`
		Quantity        int     `dynamodbav:"quantity"`
		OperationStatus string  `dynamodbav:"operationStatus"`
		Type            string  `dynamodbav:"operationType"`
		UserId          int     `dynamodbav:"userId"`
		Audit           Audit   `dynamodbav:"audit"`
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

func (eventMessage DynamoEventMessage) ToString() string {
	return fmt.Sprintf("%+v\n", eventMessage)
}

func (eventMessage DynamoEventMessage) GetKey() map[string]*dynamodb.AttributeValue {
	return map[string]*dynamodb.AttributeValue{
		"id":            {S: aws.String(eventMessage.Id)},
		"operationType": {S: aws.String(eventMessage.Type)},
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
