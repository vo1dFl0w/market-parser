package repository

import "context"

type Element interface {
	Element(ctx context.Context, selector string) (Element, error)
	
	Click(ctx context.Context) error
	Input(ctx context.Context, text string) error
	ScrollIntoView(ctx context.Context) error

	Attribute(ctx context.Context, attribute string) (string, error)
	Text(ctx context.Context) (string, error)
}
