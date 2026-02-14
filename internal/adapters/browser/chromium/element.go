package chromium

import (
	"context"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
	"github.com/vo1dFl0w/market-parser/internal/repository"
)

var workTimeout time.Duration = time.Millisecond * 5000

type rodElement struct {
	element *rod.Element
}

func (re *rodElement) Element(ctx context.Context, selector string) (repository.Element, error) {
	elem, err := re.element.Timeout(workTimeout).Context(ctx).Element(selector)
	if err != nil {
		return nil, err
	}
	return &rodElement{element: elem}, nil
}

func (re *rodElement) Click(ctx context.Context) error {
	if err := re.element.Timeout(workTimeout).Context(ctx).Click(proto.InputMouseButtonLeft, 1); err != nil {
		return err
	}

	return nil
}

func (re *rodElement) Input(ctx context.Context, text string) error {
	if err := re.element.Timeout(workTimeout).Context(ctx).Input(text); err != nil {
		return err
	}

	return nil
}

func (re *rodElement) ScrollIntoView(ctx context.Context) error {
	if err := re.element.Timeout(workTimeout).Context(ctx).ScrollIntoView(); err != nil {
		return err
	}

	return nil
}

func (re *rodElement) Attribute(ctx context.Context, attribute string) (string, error) {
	v, err := re.element.Timeout(workTimeout).Context(ctx).Attribute(attribute)
	if err != nil {
		return "", err
	}

	return *v, nil
}

func (re *rodElement) Text(ctx context.Context) (string, error) {
	t, err := re.element.Timeout(workTimeout).Context(ctx).Text()
	if err != nil {
		return "", err
	}
	return t, nil
}
