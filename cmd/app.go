package main

import (
	"context"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/joycesaquino/order-book-match-engine/internal/service"
	"github.com/joycesaquino/order-book-match-engine/internal/types"
)

func main() {
	lambda.Start(Handler)
}

// Handler TODO Add TTL em casos de Match para deleção
// Handler TODO Cloud formation
// Handler TODO Local Stack
// Handler TODO Cobertura de testes

func Handler(ctx context.Context, dynamoEvent types.DynamoEvent) error {

	sess, err := session.NewSession()
	if err != nil {
		return err
	}

	matchEngine := service.NewMatchEngine(sess)
	for _, record := range dynamoEvent.Records {

		newImage, _, err := record.ConverterEventRaw()
		if err != nil {
			return err
		}

		if newImage.Status == types.InTrade {
			matchEngine.Match(ctx, newImage)
		}
	}
	return nil

}
