package repository

import (
	"context"

	"github.com/go-rod/rod/lib/input"
	"github.com/vo1dFl0w/market-parser/internal/domain"
)

type Page interface {
	Navigate(ctx context.Context, targetURL string) error
	NavigateWithReferrer(ctx context.Context, marketURL string) error

	Element(ctx context.Context, selector string) (Element, error)
	Has(ctx context.Context, selector string) (bool, Element, error)
	HTML(ctx context.Context) (string, error)
	EachEvent(ctx context.Context) (<-chan domain.Products, <-chan error, func())
	GetPageURL(ctx context.Context) (string, error)
	MoveCursorToElement(ctx context.Context, selector string) error
	KeyboardType(ctx context.Context, key ...input.Key) error
	
	ClosePage() error
	CloseBrowser() error

	WaitStable(ctx context.Context) error
	WaitLoad(ctx context.Context) error
	WaitDOMStable(ctx context.Context) error
	WaitVisible(ctx context.Context, selector string) error
}
