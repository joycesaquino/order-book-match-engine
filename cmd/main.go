package main

import (
	"context"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/joycesaquino/order-book-match-engine/internal/service"
	"github.com/joycesaquino/order-book-match-engine/internal/types"
)

func main() {
	lambda.Start(Handler)
}

func Handler(ctx context.Context, dynamoEvent types.DynamoEvent) error {

	matchEngine := service.NewMatchEngine()
	for _, record := range dynamoEvent.Records {

		newImage, _, err := record.ConverterEventRaw()
		if err != nil {
			return err
		}

		if newImage.OperationStatus == types.InTrade && record.EventName != types.REMOVE {
			err := matchEngine.Match(ctx, newImage)
			if err != nil {
				return err
			}
		}
	}
	return nil

}
