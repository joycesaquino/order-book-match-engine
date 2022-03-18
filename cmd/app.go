package main

import (
	"context"
	"github.com/aws/aws-lambda-go/lambda"
	"order-book-match-engine/internal/event"
)

func main() {
	lambda.Start(Handler)
}

func Handler(ctx context.Context, dynamoEvent event.DynamoEvent) error {

	return nil

}
