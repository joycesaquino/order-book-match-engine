package service

import (
	orderBook "order-book-match-engine/internal/order-book"
)

type Match struct {
	Order orderBook.OrderRepository
}

func (m Match) MatchEngine() {}
