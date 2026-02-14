package chromium

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/input"
	"github.com/go-rod/rod/lib/proto"
	"github.com/vo1dFl0w/market-parser/internal/domain"
	"github.com/vo1dFl0w/market-parser/internal/repository"
)

const defaultAttemtsToSolveCaptcha int = 10
const defaultPagesNum int = 5

type rodPage struct {
	browser *rod.Browser
	page    *rod.Page
	cfg     *Config
}

func (rp *rodPage) CheckCaptcha(ctx context.Context, captchaCheckBox string, smartCaptchaSelector string) error {
	for i := 1; i <= defaultAttemtsToSolveCaptcha; i++ {
		b, _, err := rp.page.Has(smartCaptchaSelector)
		if err != nil {
			return err
		}

		if b {
			time.Sleep(time.Duration(rand.Intn(500)+500) * time.Millisecond)
			if err := rp.WaitVisible(ctx, smartCaptchaSelector); err != nil {
				return fmt.Errorf("wait visible captcha: %w", err)
			}

			captcha, err := rp.Element(ctx, smartCaptchaSelector)
			if err != nil {
				return fmt.Errorf("element captcha: %w", err)
			}

			if err := rp.MoveCursorToElement(ctx, smartCaptchaSelector); err != nil {
				return fmt.Errorf("move cursor to element captcha: %w", err)
			}

			time.Sleep(time.Duration(rand.Intn(500)+500) * time.Millisecond)
			if err := captcha.Click(ctx); err != nil {
				return fmt.Errorf("click captcha: %w", err)
			}
			time.Sleep(time.Second * 5)
			if err := rp.WaitLoad(ctx); err != nil {
				return fmt.Errorf("wait stable: %w", err)
			}
		} else {
			if i == defaultAttemtsToSolveCaptcha {
				break
			}
			time.Sleep(time.Millisecond * 300)
		}

	}

	return nil
}

func (rp *rodPage) FindCategoryElement(ctx context.Context, categorySelector string) error {
	if err := rp.WaitVisible(ctx, categorySelector); err != nil {
		return fmt.Errorf("wait visible category selector: %w", err)
	}
	categoryButton, err := rp.Element(ctx, categorySelector)
	if err != nil {
		return fmt.Errorf("element category button: %w", err)
	}

	if err := categoryButton.ScrollIntoView(ctx); err != nil {
		return fmt.Errorf("scroll into view category button: %w", err)
	}
	if err := rp.MoveCursorToElement(ctx, categorySelector); err != nil {
		return fmt.Errorf("move cursor to element category butten: %w", err)
	}
	if err := rp.WaitStable(ctx); err != nil {
		return fmt.Errorf("wait stable: %w", err)
	}

	if err := categoryButton.Click(ctx); err != nil {
		return fmt.Errorf("click category button: %w", err)
	}

	return nil
}

func (rp *rodPage) FindAllProductsBar(ctx context.Context, allProductsSelector string) error {
	if err := rp.WaitVisible(ctx, allProductsSelector); err != nil {
		return fmt.Errorf("wait visible all products selector: %w", err)
	}

	allProdsBar, err := rp.Element(ctx, allProductsSelector)
	if err != nil {
		return fmt.Errorf("element all prods bar: %w", err)
	}

	if err := allProdsBar.ScrollIntoView(ctx); err != nil {
		return fmt.Errorf("scroll into view all prods bar: %w", err)
	}
	if err := rp.MoveCursorToElement(ctx, allProductsSelector); err != nil {
		return fmt.Errorf("move cursor to element all prods bar: %w", err)
	}

	if err := allProdsBar.Click(ctx); err != nil {
		return fmt.Errorf("click all prods bar: %w", err)
	}

	if err := rp.WaitDOMStable(ctx); err != nil {
		return fmt.Errorf("wait stable: %w", err)
	}

	return nil
}

func (rp *rodPage) FindLastPageNum(ctx context.Context, lastPageSelector string, lastPageText string) (int, error) {
	lastPageNum := 0

	b, lastPageElem, err := rp.Has(ctx, lastPageSelector)
	if err != nil {
		return 0, fmt.Errorf("has last page selector: %w", err)
	}

	if b {
		lastPageText, err := lastPageElem.Element(ctx, lastPageText)
		if err != nil {
			return 0, fmt.Errorf("element last page text: %w", err)
		}
		lastPageStr, err := lastPageText.Text(ctx)
		if err != nil {
			return 0, fmt.Errorf("text last page str: %w", err)
		}
		lastPageNum, err = ParseStringToInteger(lastPageStr)
		if err != nil {
			return 0, fmt.Errorf("parse string to integer last page num: %w", err)
		}
	} else {
		lastPageNum = defaultPagesNum
	}

	return lastPageNum, nil
}

func (rp *rodPage) FindAddressButton(ctx context.Context, addressButtonSelector string) error {
	addressButton, err := rp.Element(ctx, addressButtonSelector)
	if err != nil {
		return fmt.Errorf("element address button: %w", err)
	}

	if err := rp.MoveCursorToElement(ctx, addressButtonSelector); err != nil {
		return fmt.Errorf("move cursor to element address button: %w", err)
	}

	if err := addressButton.Click(ctx); err != nil {
		return fmt.Errorf("click address button: %w", err)
	}

	if err := rp.WaitDOMStable(ctx); err != nil {
		return fmt.Errorf("wait stable: %w", err)
	}

	return nil
}

func (rp *rodPage) InputAddress(ctx context.Context, address string, addressInputSelector string) error {
	addressInput, err := rp.Element(ctx, addressInputSelector)
	if err != nil {
		return fmt.Errorf("element address input: %w", err)
	}

	if err := rp.MoveCursorToElement(ctx, addressInputSelector); err != nil {
		return fmt.Errorf("move cursor to element address input: %w", err)
	}

	if err := addressInput.Click(ctx); err != nil {
		return fmt.Errorf("click address input: %w", err)
	}

	if err := rp.KeyboardType(ctx, input.ControlLeft, input.KeyA, input.Delete); err != nil {
		return fmt.Errorf("keyboard type address input: %w", err)
	}

	if err := addressInput.Input(ctx, address); err != nil {
		return fmt.Errorf("input address input: %w", err)
	}

	if err := rp.WaitDOMStable(ctx); err != nil {
		return fmt.Errorf("wait dom stable address input: %w", err)
	}

	return nil
}

func (rp *rodPage) ClickDropDownAddress(ctx context.Context, addressInputDropDownSelector string) error {
	addressDropDown, err := rp.Element(ctx, addressInputDropDownSelector)
	if err != nil {
		return fmt.Errorf("element address drop down: %w", err)
	}

	if err := rp.MoveCursorToElement(ctx, addressInputDropDownSelector); err != nil {
		return fmt.Errorf("move cursor to element address drop down: %w", err)
	}

	if err := addressDropDown.Click(ctx); err != nil {
		return fmt.Errorf("click address drop down: %w", err)
	}

	if err := rp.WaitDOMStable(ctx); err != nil {
		return fmt.Errorf("wait stable: %w", err)
	}

	return nil
}

func (rp *rodPage) SaveDeliveryAddress(ctx context.Context, addressSaveButtonSelector string) error {
	if err := rp.WaitVisible(ctx, addressSaveButtonSelector); err != nil {
		return fmt.Errorf("wait visible address save button selector: %w", err)
	}

	addressSave, err := rp.Element(ctx, addressSaveButtonSelector)
	if err != nil {
		return fmt.Errorf("element address save: %w", err)
	}

	if err := rp.MoveCursorToElement(ctx, addressSaveButtonSelector); err != nil {
		return fmt.Errorf("move cursor to element address save: %w", err)

	}

	if err := addressSave.Click(ctx); err != nil {
		return fmt.Errorf("click address save: %w", err)
	}

	if err := rp.WaitStable(ctx); err != nil {
		return fmt.Errorf("wait stable: %w", err)
	}

	return nil
}

func (rp *rodPage) ParsePages(ctx context.Context, lastPageNum int) ([]domain.Products, error) {
	result := []domain.Products{}

	// формируем базовый url, для дальнейшей навигации по страницам basePageURL+&page=1,2,3...
	basePageURL, err := rp.GetPageURL(ctx)
	if err != nil {
		return nil, fmt.Errorf("get page url: %w", err)
	}
	// собираем информацию с каждой страницы
	for i := 1; i <= lastPageNum; i++ {
		// начать перехват тела ответа запроса, который содержит данные о товарах
		//fmt.Printf("начат перехват запроса на странице %d\n", i)
		resCh, errCh, stopListeningFn := rp.EachEvent(ctx)

		targetURL := fmt.Sprintf("%s&page=%d", basePageURL, i)
		if err := rp.Navigate(ctx, targetURL); err != nil {
			stopListeningFn()
			return nil, fmt.Errorf("navigate %s: %w", targetURL, err)
		}
		count := 0
		done := false
		for !done {
			select {
			case r, ok := <-resCh:
				if !ok {
					//fmt.Printf("Канал закрыт | Страница: %d | Итого товаров: %d\n", i, len(result))
					done = true
					break
				}
				result = append(result, r)
				count++
			case err, ok := <-errCh:
				if !ok || err == nil {
					continue
				}
				//fmt.Printf("Получена ошибка: %s\n", err.Error())

				return nil, err
			case <-ctx.Done():
				//fmt.Printf("Контекст завершился раньше канала!\n")

				stopListeningFn()
				return nil, ctx.Err()
			}
		}
		stopListeningFn()
		//fmt.Printf("[Parser] Данных получено: %d | страница: %d\n", count, i)
	}

	return result, nil
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

	return &rodElement{element: elem}, err
}

func (rp *rodPage) MoveCursorToElement(ctx context.Context, selector string) error {
	// поиск элемента и его координат
	elem, err := rp.page.Timeout(rp.cfg.WorkTimeout).Context(ctx).Element(selector)
	if err != nil {
		return err
	}

	shape, err := elem.Timeout(rp.cfg.WorkTimeout).Context(ctx).Shape()
	if err != nil {
		return err
	}
	box := shape.Box()

	// целевая точка внутри элемента со случайным смещением (от 20% до 80% размера)
	targetX := box.X + box.Width*(0.2+rand.Float64()*0.6)
	targetY := box.Y + box.Height*(0.2+rand.Float64()*0.6)

	// определение начальной точки (если текущая неизвестна, берем случайную на периферии)
	pos := rp.page.Mouse.Position()
	startX, startY := pos.X, pos.Y

	// решаем, будет ли проскок (50% вероятность)
	hasOvershoot := rand.Float32() < 0.5

	points := []struct{ x, y float64 }{}

	// если проскок, сначала идем к "ошибочной" точке
	if hasOvershoot {
		dist := 5.0 + rand.Float64()*12.0 // На сколько пикселей пролетим мимо
		overX := targetX + (targetX-startX)*0.05 + dist
		overY := targetY + (targetY-startY)*0.05 + dist
		points = append(points, struct{ x, y float64 }{overX, overY})
	}

	// конечная цель всегда в списке последней
	points = append(points, struct{ x, y float64 }{targetX, targetY})

	// выполнение движений
	currentX, currentY := startX, startY
	for idx, pt := range points {
		// генерация контрольной точки для кривой Безье
		controlX := currentX + (pt.x-currentX)*rand.Float64() + float64(rand.Intn(40)-20)
		controlY := currentY + (pt.y-currentY)*rand.Float64() + float64(rand.Intn(40)-20)

		steps := 15 + rand.Intn(10)
		// для корректирующего движения после проскока нужно меньше шагов
		if idx > 0 {
			steps = 5 + rand.Intn(5)
			time.Sleep(time.Duration(rand.Intn(50)+30) * time.Millisecond) // Пауза "осознания ошибки"
		}

		for i := 1; i <= steps; i++ {
			t := float64(i) / float64(steps)

			// формула Безье
			curX := (1-t)*(1-t)*currentX + 2*(1-t)*t*controlX + t*t*pt.x
			curY := (1-t)*(1-t)*currentY + 2*(1-t)*t*controlY + t*t*pt.y

			// jitter (микро-дрожание)
			curX += rand.Float64()*1.2 - 0.6
			curY += rand.Float64()*1.2 - 0.6

			if err := rp.page.Mouse.MoveLinear(proto.Point{X: curX, Y: curY}, 1); err != nil {
				return err
			}

			// динамическая пауза (замедление к концу каждого отрезка)
			pause := 4 + int(t*12) + rand.Intn(4)
			time.Sleep(time.Duration(pause) * time.Millisecond)
		}
		// обновляем текущую позицию для следующего отрезка (если он есть)
		currentX, currentY = pt.x, pt.y
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

	wg.Add(1)
	go func() {
		defer wg.Done()
		defer close(resCh)
		defer close(errCh)

		waitFn := rp.page.Timeout(rp.cfg.WorkTimeout).Context(ctxEvent).EachEvent(func(r *proto.NetworkResponseReceived) bool {
			//fmt.Printf("[EachEvent] Запущен...\n\n")
			if strings.Contains(r.Response.URL, "products") {

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
	html, err := rp.page.Context(ctx).HTML()
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
