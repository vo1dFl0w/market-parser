package parsers

import (
	"context"
	"fmt"
	"strings"

	"github.com/vo1dFl0w/market-parser/internal/config"
	"github.com/vo1dFl0w/market-parser/internal/domain"
	"github.com/vo1dFl0w/market-parser/internal/repository"
	"github.com/vo1dFl0w/market-parser/pkg/logger"
)

const testDefaultLastPageNum int = 3

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

	// создание и переход на сайт kuper.ru
	page, err := kp.browser.NewPage(ctx, kp.cfg.BaseURL)
	if err != nil {
		return nil, fmt.Errorf("new page: %w", err)
	}
	defer page.CloseBrowser()
	defer page.ClosePage()

	if err := page.WaitLoad(ctx); err != nil {
		return nil, fmt.Errorf("wait dom stable: %w", err)
	}
	if err := page.WaitStable(ctx); err != nil {
		return nil, fmt.Errorf("wait dom stable: %w", err)
	}
	if err := page.CheckCaptcha(ctx, kp.cfg.Selectors.CaptchaCheckBox, kp.cfg.Selectors.SmartCaptchaSelector); err != nil {
		return nil, fmt.Errorf("check captcha: %w", err)
	}

	// переходим по url к заданному market
	marketPageURL := fmt.Sprintf("%s/%s", kp.cfg.BaseURL, market)
	if err := page.Navigate(ctx, marketPageURL); err != nil {
		return nil, fmt.Errorf("navigate with referrer %s: %w", marketPageURL, err)
	}
	if err := page.WaitLoad(ctx); err != nil {
		return nil, fmt.Errorf("wait dom stable: %w", err)
	}
	if err := page.CheckCaptcha(ctx, kp.cfg.Selectors.CaptchaCheckBox, kp.cfg.Selectors.SmartCaptchaSelector); err != nil {
		return nil, fmt.Errorf("check captcha: %w", err)
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
			// адрес не установлен, пытаемся внести
			if err := page.FindAddressButton(ctx, selector.AddressButtonSelector); err != nil {
				return nil, fmt.Errorf("find address button: %w", err)
			}

			// ввести адрес
			if err := page.InputAddress(ctx, address, selector.AddressInputSelector); err != nil {
				return nil, fmt.Errorf("input address: %w", err)
			}

			// нажать на вспылвший адрес
			if err := page.ClickDropDownAddress(ctx, selector.AddressInputDropDownSelector); err != nil {
				return nil, fmt.Errorf("click drop down address: %w", err)
			}

			// сохранить адрес
			if err := page.SaveDeliveryAddress(ctx, kp.cfg.Selectors.AddressSaveButtonSelector); err != nil {
				return nil, fmt.Errorf("save delivery address: %w", err)
			}
		}
	} else {
		// адрес не установлен, пытаемся внести
		if err := page.FindAddressButton(ctx, selector.AddressButtonSelector); err != nil {
			return nil, fmt.Errorf("find address button: %w", err)
		}

		// ввести адрес
		if err := page.InputAddress(ctx, address, selector.AddressInputSelector); err != nil {
			return nil, fmt.Errorf("input address: %w", err)
		}

		// нажать на вспылвший адрес
		if err := page.ClickDropDownAddress(ctx, selector.AddressInputDropDownSelector); err != nil {
			return nil, fmt.Errorf("click drop down address: %w", err)
		}

		// сохранить адрес
		if err := page.SaveDeliveryAddress(ctx, kp.cfg.Selectors.AddressSaveButtonSelector); err != nil {
			return nil, fmt.Errorf("save delivery address: %w", err)
		}
	}

	// находим селектор с категорией
	categorySelector := fmt.Sprintf("span[title='%s']", category)
	if err := page.FindCategoryElement(ctx, categorySelector); err != nil {
		return nil, fmt.Errorf("find category element")
	}

	//  находим строку "все товары категории"
	if err := page.FindAllProductsBar(ctx, selector.AllProdsSelector); err != nil {
		return nil, fmt.Errorf("find all products bar: %w", err)
	}

	// ищем значение последней страницы
	lastPageNum, err := page.FindLastPageNum(ctx, kp.cfg.Selectors.LastPageSelector, kp.cfg.Selectors.LastPageText)
	if err != nil {
		return nil, fmt.Errorf("find last page num: %w", err)
	}

	// если test parser mode, то lastPageNum=testDefaultLastPageNum для тестирования функционала
	if kp.cfg.TestParserMode {
		lastPageNum = testDefaultLastPageNum
	}

	res, err := page.ParsePages(ctx, lastPageNum)
	if err != nil {
		return nil, fmt.Errorf("parse pages: %w", err)
	}

	return res, nil
}
