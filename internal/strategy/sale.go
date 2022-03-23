package strategy

import (
	"context"
	"log"
	orderBook "order-book-match-engine/internal/order-book"
	"order-book-match-engine/internal/types"
	walletIntegration "order-book-match-engine/internal/wallet-integration"
)

type Sale struct {
	repository orderBook.OperationRepository
	queue      *walletIntegration.Queue
}

func (strategy *Sale) Apply(ctx context.Context, input *Input) {
	log.Println("implement me")
}

func (strategy *Sale) Accept(input *Input) bool {
	return input.NewImage.Type == types.Sale
}
