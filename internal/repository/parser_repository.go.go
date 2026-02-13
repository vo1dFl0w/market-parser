package repository

import (
	"context"

	"github.com/vo1dFl0w/market-parser/internal/domain"
)

type ParserRepository interface {
	GetAllProductsByCategory(ctx context.Context, category string, address string, market string) ([]domain.Products, error)
}
