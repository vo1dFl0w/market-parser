package repository

import (
	"context"

	"github.com/go-rod/rod/lib/input"
	"github.com/vo1dFl0w/market-parser/internal/domain"
)

type Page interface {
	// high-level operations
	CheckCaptcha(ctx context.Context, captchaCheckBox string, smartCaptchaSelector string) error
	FindCategoryElement(ctx context.Context, categorySelector string) error
	FindAllProductsBar(ctx context.Context, allProductsBarSelector string) error
	FindLastPageNum(ctx context.Context, lastPageSelector string, lastPageText string) (int, error)
	FindAddressButton(ctx context.Context, addressButtonSelector string) error
	InputAddress(ctx context.Context, address string, addressInputSelector string) error
	ClickDropDownAddress(ctx context.Context, addressInputDropDownSelector string) error
	SaveDeliveryAddress(ctx context.Context, addressSaveButtonSelector string) error
	ParsePages(ctx context.Context, lastPageNum int) ([]domain.Products, error)

	// low-level methods
	// navigation
	Navigate(ctx context.Context, targetURL string) error
	NavigateWithReferrer(ctx context.Context, marketURL string) error

	// operations with page elements
	Element(ctx context.Context, selector string) (Element, error)
	Has(ctx context.Context, selector string) (bool, Element, error)
	HTML(ctx context.Context) (string, error)
	EachEvent(ctx context.Context) (<-chan domain.Products, <-chan error, func())
	GetPageURL(ctx context.Context) (string, error)
	MoveCursorToElement(ctx context.Context, selector string) error
	KeyboardType(ctx context.Context, key ...input.Key) error


	// wait opertaions
	WaitStable(ctx context.Context) error
	WaitLoad(ctx context.Context) error
	WaitDOMStable(ctx context.Context) error
	WaitVisible(ctx context.Context, selector string) error

	// close operations
	ClosePage() error
	CloseBrowser() error
}
