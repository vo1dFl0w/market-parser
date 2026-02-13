package chromium

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/input"
	"github.com/go-rod/rod/lib/proto"
	"github.com/vo1dFl0w/market-parser/internal/domain"
	"github.com/vo1dFl0w/market-parser/internal/repository"
)

type rodPage struct {
	browser *rod.Browser
	page    *rod.Page
	cfg     *Config
}

func (rp *rodPage) Navigate(ctx context.Context, targetURL string) error {
	if err := rp.page.Timeout(rp.cfg.WorkTimeout).Context(ctx).Navigate(targetURL); err != nil {
		return err
	}
	return nil
}

func (rp *rodPage) NavigateWithReferrer(ctx context.Context, marketURL string) error {
	_, err := proto.PageNavigate{
		URL:      marketURL,
		Referrer: rp.cfg.Referrer,
	}.Call(rp.page)
	if err != nil {
		return err
	}

	return nil
}

func (rp *rodPage) Element(ctx context.Context, selector string) (repository.Element, error) {
	elem, err := rp.page.Timeout(rp.cfg.WorkTimeout).Context(ctx).Element(selector)
	if err != nil {
		return nil, err
	}

	return &rodElement{
		element: elem,
	}, err
}

func (rp *rodPage) MoveCursorToElement(ctx context.Context, selector string) error {
	elem, err := rp.page.Timeout(rp.cfg.WorkTimeout).Context(ctx).Element(selector)
	if err != nil {
		return err
	}

	shape, err := elem.Timeout(rp.cfg.WorkTimeout).Context(ctx).Shape()
	if err != nil {
		return err
	}

	box := shape.Box()
	centerX := box.X + box.Width/2
	centerY := box.Y + box.Height/2

	if err := rp.page.Mouse.MoveLinear(proto.Point{X: centerX, Y: centerY}, 15); err != nil {
		return err
	}

	return nil
}

func (rp *rodPage) KeyboardType(ctx context.Context, key ...input.Key) error {
	if err := rp.page.Timeout(rp.cfg.WorkTimeout).Context(ctx).KeyActions().Press(key...).Do(); err != nil {
		return err
	}

	return nil
}

func (rp *rodPage) EachEvent(ctx context.Context) (<-chan domain.Products, <-chan error, func()) {
	ctxEvent, cancel := context.WithCancel(ctx)

	resCh := make(chan domain.Products, 100)
	errCh := make(chan error, 1)
	wg := &sync.WaitGroup{}

	err := proto.NetworkEnable{}.Call(rp.page)
	if err != nil {
		errCh <- fmt.Errorf("network enable: %w", err)
	}

	go func() {
		defer close(resCh)
		defer close(errCh)

		waitFn := rp.page.Timeout(rp.cfg.WorkTimeout).Context(ctxEvent).EachEvent(func(r *proto.NetworkResponseReceived) bool {
			//fmt.Printf("[EachEvent] Запущен...\n\n")
			if strings.Contains(r.Response.URL, "products") {
				wg.Add(1)

				if strings.Contains(r.Response.URL, "products") {
					res, err := proto.NetworkGetResponseBody{RequestID: r.RequestID}.Call(rp.page)
					if err != nil {
						if strings.Contains(err.Error(), "-32000") {
							return false
						}
						errCh <- fmt.Errorf("network get reponse body: %w", err)
						return false
					}

					if res.Body != "" {
						count := 0
						data := &ProductsResponse{}
						if err := json.Unmarshal([]byte(res.Body), &data); err == nil {
							for _, p := range data.Prods {
								count++
								resCh <- domain.Products{Name: p.Name, Price: p.Price, URL: p.CanonicalURL}
								//fmt.Printf("[EachEvent] nдобалено в page.go %v\n\n", p)
							}
							//fmt.Printf("[EachEvent] Обработано данных в page.go %d\n", count)
							return true
						}
					}
				}
			}
			return false
		})
		//fmt.Printf("[EachEvent] останавливаем слушатель...\n\n")
		waitFn()
	}()

	go func() {
		wg.Wait()
	}()

	stopListeningFn := func() {
		cancel()
	}

	return resCh, errCh, stopListeningFn
}

func (rp *rodPage) GetPageURL(ctx context.Context) (string, error) {
	info, err := rp.page.Timeout(rp.cfg.WorkTimeout).Context(ctx).Info()
	if err != nil {
		return "", err
	}
	return info.URL, nil
}

func (rp *rodPage) Has(ctx context.Context, selector string) (bool, repository.Element, error) {
	b, elem, err := rp.page.Timeout(rp.cfg.WorkTimeout).Context(ctx).Has(selector)
	if err != nil {
		return false, nil, err
	}
	if b {
		return b, &rodElement{element: elem}, nil
	}

	return false, nil, nil
}

func (rp *rodPage) HTML(ctx context.Context) (string, error) {
	html, err := rp.page.Timeout(rp.cfg.WorkTimeout).Context(ctx).HTML()
	if err != nil {
		return "", err
	}
	return html, nil
}

func (rp *rodPage) WaitStable(ctx context.Context) error {
	if err := rp.page.Timeout(rp.cfg.WorkTimeout).Context(ctx).WaitStable(rp.cfg.WaitStableDuration); err != nil {
		return fmt.Errorf("wait stable page: %w", err)
	}

	return nil
}

func (rp *rodPage) WaitLoad(ctx context.Context) error {
	if err := rp.page.Timeout(rp.cfg.WorkTimeout).Context(ctx).WaitLoad(); err != nil {
		return fmt.Errorf("wait load page: %w", err)
	}

	return nil
}

func (rp *rodPage) WaitDOMStable(ctx context.Context) error {
	if err := rp.page.Timeout(rp.cfg.WorkTimeout).Context(ctx).WaitDOMStable(rp.cfg.WaitStableDuration, rp.cfg.WaitDOMStableDiff); err != nil {
		return fmt.Errorf("wait dom stable page: %w", err)
	}

	return nil
}

func (rp *rodPage) WaitVisible(ctx context.Context, selector string) error {
	el, err := rp.page.Timeout(rp.cfg.WorkTimeout).Context(ctx).Element(selector)
	if err != nil {
		return err
	}
	if err := el.Timeout(rp.cfg.WorkTimeout).Context(ctx).WaitVisible(); err != nil {
		return err
	}
	return nil
}

func (rp *rodPage) ClosePage() error {
	if rp.page != nil {
		if err := rp.page.Close(); err != nil {
			return err
		}
	}

	return nil
}

func (rp *rodPage) CloseBrowser() error {
	if rp.page != nil {
		if err := rp.browser.Close(); err != nil {
			return err
		}
	}

	return nil
}
