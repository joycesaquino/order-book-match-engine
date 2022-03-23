package main

import (
	"context"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"order-book-match-engine/internal/service"
	"order-book-match-engine/internal/types"
)

func main() {
	lambda.Start(Handler)
}

// Handler TODO Match Engine for sale
// Handler TODO Add TTL em casos de Match para deleção
// Handler TODO Cloud formation
// Handler TODO Local Stack

func Handler(ctx context.Context, dynamoEvent types.DynamoEvent) error {
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
