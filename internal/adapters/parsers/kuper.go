package parsers

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-rod/rod/lib/input"
	"github.com/vo1dFl0w/market-parser/internal/config"
	"github.com/vo1dFl0w/market-parser/internal/domain"
	"github.com/vo1dFl0w/market-parser/internal/repository"
	"github.com/vo1dFl0w/market-parser/pkg/logger"
)

type Kuper interface {
	GetAllProductsByCategory(ctx context.Context, category string, address string) ([]domain.Products, error)
}

type kuper struct {
	cfg     *KuperConfig
	browser repository.BrowserRepository
	logger  logger.Logger
}

func NewKuperParser(cfg *config.Config, logger logger.Logger, browser repository.BrowserRepository) *kuper {
	return &kuper{
		cfg:     NewKuperConfig(cfg),
		browser: browser,
		logger:  logger,
	}
}

func (kp *kuper) GetAllProductsByCategory(ctx context.Context, category string, address string, market string) ([]domain.Products, error) {
	selector := kp.cfg.Selectors

	page, err := kp.browser.NewPage(ctx)
	if err != nil {
		return nil, fmt.Errorf("new page: %w", err)
	}
	defer page.CloseBrowser()
	defer page.ClosePage()

	// переход на сайт kuper.ru
	marketPageURL := fmt.Sprintf("%s/%s", kp.cfg.BaseURL, market)
	if err := page.NavigateWithReferrer(ctx, marketPageURL); err != nil {
		return nil, fmt.Errorf("navigate with referrer %s: %w", marketPageURL, err)
	}
	if err := page.WaitLoad(ctx); err != nil {
		return nil, fmt.Errorf("wait dom stable: %w", err)
	}

	// проверяем не получен ли блок или капча после перехода
	html, err := page.HTML(ctx)
	if err != nil {
		return nil, fmt.Errorf("html page: %w", err)
	}
	if strings.Contains(html, "captcha") {
		return nil, fmt.Errorf("captcha or block detected: %w", err)
	}

	// проверяем установлен ли уже адрес на сайте
	b, currentAddrBar, err := page.Has(ctx, selector.CurrentAddressSelector)
	if err != nil {
		return nil, fmt.Errorf("current addr bar: %w", err)
	}
	// если да, то смотрим содержимое адреса
	if b {
		addrText, err := currentAddrBar.Text(ctx)
		if err != nil {
			return nil, fmt.Errorf("text addr text: %w", err)
		}
		// если адрес не совпадает, то пытаемся его поменять, если нет, то пропускаем
		if !strings.Contains(addrText, address) {
			// нажать на кнопку для задания адреса доставки address
			addressButton, err := page.Element(ctx, selector.AddressButtonSelector)
			if err != nil {
				return nil, fmt.Errorf("element address button: %w", err)
			}
			if kp.cfg.HumanLikeMode {
				if err := page.MoveCursorToElement(ctx, selector.AddressButtonSelector); err != nil {
					return nil, fmt.Errorf("move cursor to element address button: %w", err)
				}
			}
			if err := addressButton.Click(ctx); err != nil {
				return nil, fmt.Errorf("click address button: %w", err)
			}
			if err := page.WaitDOMStable(ctx); err != nil {
				return nil, fmt.Errorf("wait stable: %w", err)
			}

			// проверика совпадает ли адрес с тем, что есть в адресной строке
			addressInput, err := page.Element(ctx, selector.AddressInputSelector)
			if err != nil {
				return nil, fmt.Errorf("element address input: %w", err)
			}

			// нажать на строку для ввода адреса
			if kp.cfg.HumanLikeMode {
				if err := page.MoveCursorToElement(ctx, selector.AddressInputSelector); err != nil {
					return nil, fmt.Errorf("move cursor to element address input: %w", err)
				}
			}
			if err := addressInput.Click(ctx); err != nil {
				return nil, fmt.Errorf("click address input: %w", err)
			}
			if err := page.KeyboardType(ctx, input.ControlLeft, input.KeyA, input.Delete); err != nil {
				return nil, fmt.Errorf("keyboard type address input: %w", err)
			}
			if err := addressInput.Input(ctx, address); err != nil {
				return nil, fmt.Errorf("input address input: %w", err)
			}
			if err := page.WaitDOMStable(ctx); err != nil {
				return nil, fmt.Errorf("wait stable: %w", err)
			}

			// нажать на вспылвший адрес
			addressDropDown, err := page.Element(ctx, selector.AddressInputDropDownSelector)
			if err != nil {
				return nil, fmt.Errorf("element address drop down: %w", err)
			}
			if kp.cfg.HumanLikeMode {
				if err := page.MoveCursorToElement(ctx, selector.AddressInputDropDownSelector); err != nil {
					return nil, fmt.Errorf("move cursor to element address drop down: %w", err)
				}
			}
			if err := addressDropDown.Click(ctx); err != nil {
				return nil, fmt.Errorf("click address drop down: %w", err)
			}
			if err := page.WaitDOMStable(ctx); err != nil {
				return nil, fmt.Errorf("wait stable: %w", err)
			}

			// сохранить адрес
			if err := page.WaitVisible(ctx, selector.AddressSaveButtonSelector); err != nil {
				return nil, fmt.Errorf("wait stable: %w", err)
			}
			addressSave, err := page.Element(ctx, selector.AddressSaveButtonSelector)
			if err != nil {
				return nil, fmt.Errorf("element address save: %w", err)
			}
			if kp.cfg.HumanLikeMode {
				if err := page.MoveCursorToElement(ctx, selector.AddressSaveButtonSelector); err != nil {
					return nil, fmt.Errorf("move cursor to element address save: %w", err)
				}
			}
			if err := addressSave.Click(ctx); err != nil {
				return nil, fmt.Errorf("click address save: %w", err)
			}
			if err := page.WaitStable(ctx); err != nil {
				return nil, fmt.Errorf("wait stable: %w", err)
			}
		}
	} else {
		// адрес не установлен, пытаемся внести
		addressButton, err := page.Element(ctx, selector.AddressButtonSelector)
		if err != nil {
			return nil, fmt.Errorf("element address button: %w", err)
		}
		if kp.cfg.HumanLikeMode {
			if err := page.MoveCursorToElement(ctx, selector.AddressButtonSelector); err != nil {
				return nil, fmt.Errorf("move cursor to element address button: %w", err)
			}
		}
		if err := addressButton.Click(ctx); err != nil {
			return nil, fmt.Errorf("click address button: %w", err)
		}
		if err := page.WaitDOMStable(ctx); err != nil {
			return nil, fmt.Errorf("wait stable: %w", err)
		}

		// проверика совпадает ли адрес с тем, что есть в адресной строке
		addressInput, err := page.Element(ctx, selector.AddressInputSelector)
		if err != nil {
			return nil, fmt.Errorf("element address input: %w", err)
		}

		// нажать на строку для ввода адреса
		if kp.cfg.HumanLikeMode {
			if err := page.MoveCursorToElement(ctx, selector.AddressInputSelector); err != nil {
				return nil, fmt.Errorf("move cursor to element address input: %w", err)
			}
		}
		if err := addressInput.Click(ctx); err != nil {
			return nil, fmt.Errorf("click address input: %w", err)
		}
		if err := page.KeyboardType(ctx, input.ControlLeft, input.KeyA, input.Delete); err != nil {
			return nil, fmt.Errorf("keyboard type address input: %w", err)
		}
		if err := addressInput.Input(ctx, address); err != nil {
			return nil, fmt.Errorf("input address input: %w", err)
		}
		if err := page.WaitDOMStable(ctx); err != nil {
			return nil, fmt.Errorf("wait stable: %w", err)
		}

		// нажать на вспылвший адрес
		addressDropDown, err := page.Element(ctx, selector.AddressInputDropDownSelector)
		if err != nil {
			return nil, fmt.Errorf("element address drop down: %w", err)
		}
		if kp.cfg.HumanLikeMode {
			if err := page.MoveCursorToElement(ctx, selector.AddressInputDropDownSelector); err != nil {
				return nil, fmt.Errorf("move cursor to element address drop down: %w", err)
			}
		}
		if err := addressDropDown.Click(ctx); err != nil {
			return nil, fmt.Errorf("click address drop down: %w", err)
		}
		if err := page.WaitDOMStable(ctx); err != nil {
			return nil, fmt.Errorf("wait stable: %w", err)
		}

		// сохранить адрес
		if err := page.WaitVisible(ctx, selector.AddressSaveButtonSelector); err != nil {
			return nil, fmt.Errorf("wait stable: %w", err)
		}
		addressSave, err := page.Element(ctx, selector.AddressSaveButtonSelector)
		if err != nil {
			return nil, fmt.Errorf("element address save: %w", err)
		}
		if kp.cfg.HumanLikeMode {
			if err := page.MoveCursorToElement(ctx, selector.AddressSaveButtonSelector); err != nil {
				return nil, fmt.Errorf("move cursor to element address save: %w", err)
			}
		}
		if err := addressSave.Click(ctx); err != nil {
			return nil, fmt.Errorf("click address save: %w", err)
		}
		if err := page.WaitStable(ctx); err != nil {
			return nil, fmt.Errorf("wait stable: %w", err)
		}
	}

	// находим селектор с категорией

	categorySelector := fmt.Sprintf("span[title='%s']", category)
	categoryButton, err := page.Element(ctx, categorySelector)
	if err != nil {
		return nil, fmt.Errorf("element category button: %w", err)
	}
	if kp.cfg.HumanLikeMode {
		if err := categoryButton.ScrollIntoView(ctx); err != nil {
			return nil, fmt.Errorf("scroll into view category button: %w", err)
		}
		if err := page.MoveCursorToElement(ctx, categorySelector); err != nil {
			return nil, fmt.Errorf("move cursor to element category butten: %w", err)
		}
		if err := page.WaitStable(ctx); err != nil {
			return nil, fmt.Errorf("wait stable: %w", err)
		}
	}
	if err := categoryButton.Click(ctx); err != nil {
		return nil, fmt.Errorf("click category button: %w", err)
	}

	//  находим строку "все товары категории"
	if err := page.WaitVisible(ctx, selector.AllProdsSelector); err != nil {
		return nil, fmt.Errorf("wait stable: %w", err)
	}
	allProdsBar, err := page.Element(ctx, selector.AllProdsSelector)
	if err != nil {
		return nil, fmt.Errorf("element all prods bar: %w", err)
	}
	if kp.cfg.HumanLikeMode {
		if err := allProdsBar.ScrollIntoView(ctx); err != nil {
			return nil, fmt.Errorf("scroll into view all prods bar: %w", err)
		}
		if err := page.MoveCursorToElement(ctx, selector.AllProdsSelector); err != nil {
			return nil, fmt.Errorf("move cursor to element all prods bar: %w", err)
		}
	}
	if err := allProdsBar.Click(ctx); err != nil {
		return nil, fmt.Errorf("click all prods bar: %w", err)
	}
	if err := page.WaitDOMStable(ctx); err != nil {
		return nil, fmt.Errorf("wait stable: %w", err)
	}

	// ищем значение последней страницы
	lastPageSelector, err := page.Element(ctx, selector.LastPageSelector)
	if err != nil {
		return nil, fmt.Errorf("element last page: %w", err)
	}
	lastPageText, err := lastPageSelector.Element(ctx, selector.LastPageText)
	if err != nil {
		return nil, fmt.Errorf("element last page text: %w", err)
	}
	lastPageStr, err := lastPageText.Text(ctx)
	if err != nil {
		return nil, fmt.Errorf("text last page str: %w", err)
	}
	lastPageNum, err := ParseStringToInteger(lastPageStr)
	if err != nil {
		return nil, fmt.Errorf("parse string to integer last page num: %w", err)
	}

	// если test parser mode, то lastPageNum=3 для тестирования функционала
	if kp.cfg.TestParserMode {
		lastPageNum = 3
	}

	basePageURL, err := page.GetPageURL(ctx)
	if err != nil {
		return nil, fmt.Errorf("get page url: %w", err)
	}

	result := []domain.Products{}

	// собираем информацию с каждой страницы
	for i := 1; i <= lastPageNum; i++ {
		// начать перехват тела ответа запроса, который содержит данные о товарах
		//fmt.Printf("начат перехват запроса на странице %d\n", i)
		resCh, errCh, stopListeningFn := page.EachEvent(ctx)

		targetURL := fmt.Sprintf("%s&page=%d", basePageURL, i)
		if err := page.Navigate(ctx, targetURL); err != nil {
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

	//fmt.Printf("[Parser] Всего данных получено: %d\n", len(result))
	return result, nil
}
