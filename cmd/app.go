package main

import (
	"context"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"order-book-match-engine/internal/event"
	"order-book-match-engine/internal/service"
)

func main() {
	lambda.Start(Handler)
}

func Handler(ctx context.Context, dynamoEvent event.DynamoEvent) error {
	sess, err := session.NewSession()
	if err != nil {
		return err
	}

	matchEngine := service.NewMatchEngine(sess)
	for _, record := range dynamoEvent.Records {
		go matchEngine.Match(ctx, record)
	}
	return nil

}
