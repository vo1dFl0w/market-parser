package repository

import (
	"context"

	"github.com/go-rod/rod"
)

type BrowserRepository interface {
	Connect(ctx context.Context) (*rod.Browser, error)
	NewPage(ctx context.Context, markeyURL string) (Page, error)
}