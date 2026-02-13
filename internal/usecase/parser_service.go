package usecase

import (
	"context"
	"fmt"

	"github.com/vo1dFl0w/market-parser/internal/domain"
	"github.com/vo1dFl0w/market-parser/internal/repository"
)

type ParserService interface {
	ParseProductsByCategory(ctx context.Context, category string, address string, market string) ([]domain.Products, error)
}

type parserService struct {
	parserRepo repository.ParserRepository
}

func NewParserService(parserRepo repository.ParserRepository) *parserService {
	return &parserService{parserRepo: parserRepo}
}

func (s *parserService) ParseProductsByCategory(ctx context.Context, category string, address string, market string) ([]domain.Products, error) {
	if category == "" {
		return nil, domain.ErrEmptyCategory
	}

	if address == "" {
		return nil, domain.ErrEmptyAddress
	}

	if market == "" {
		return nil, domain.ErrEmptyMarket
	}

	res, err := s.parserRepo.GetAllProductsByCategory(ctx, category, address, market)
	if err != nil {
		return nil, fmt.Errorf("get all products by category: %w", err)
	}
	return res, nil
}
