package strategy

import (
	orderBook "order-book-match-engine/internal/order-book"
	"order-book-match-engine/internal/types"
	walletIntegration "order-book-match-engine/internal/wallet-integration"
)

type Sale struct {
	repository orderBook.OperationRepository
	queue      *walletIntegration.Queue
}

func (strategy *Sale) Accept(input *Input) bool {
	return input.NewImage.Type == types.Sale
}
