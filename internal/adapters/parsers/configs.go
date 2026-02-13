package parsers

import (
	"github.com/vo1dFl0w/market-parser/internal/config"
)

type KuperConfig struct {
	TestParserMode  bool
	HumanLikeMode   bool
	ApiProductsPath string
	BaseURL         string
	Referrer        string
	Selectors       *KuperSelectors
}

type KuperSelectors struct {
	CurrentAddressSelector       string
	AddressButtonSelector        string
	AddressCheckAttributeValue   string
	AddressInputSelector         string
	AddressInputDropDownSelector string
	AddressSaveButtonSelector    string
	MarketSelector               string
	AllProdsSelector             string
	LastPageSelector             string
	LastPageText                 string
	NextPageSelector             string
}

func NewKuperConfig(cfg *config.Config) *KuperConfig {
	return &KuperConfig{
		TestParserMode:  cfg.Browser.TestParserMode,
		HumanLikeMode:   cfg.Browser.HumanLikeMode,
		ApiProductsPath: *cfg.Server.KuperCfg.ApiProductsPath,
		BaseURL:         *cfg.Server.KuperCfg.BaseURL,
		Referrer:        cfg.Browser.Referer,
		Selectors: &KuperSelectors{
			CurrentAddressSelector:       *cfg.Server.KuperCfg.CurrentAddressSelector,
			AddressButtonSelector:        *cfg.Server.KuperCfg.AddressButtonSelector,
			AddressCheckAttributeValue:   *cfg.Server.KuperCfg.AddressCheckAttributeValue,
			AddressInputSelector:         *cfg.Server.KuperCfg.AddressInputSelector,
			AddressInputDropDownSelector: *cfg.Server.KuperCfg.AddressInputDropDownSelector,
			AddressSaveButtonSelector:    *cfg.Server.KuperCfg.AddressSaveButtonSelector,
			MarketSelector:               *cfg.Server.KuperCfg.MarketSelector,
			AllProdsSelector:             *cfg.Server.KuperCfg.AllProdsSelector,
			LastPageSelector:             *cfg.Server.KuperCfg.LastPageSelector,
			LastPageText:                 *cfg.Server.KuperCfg.LastPageText,
			NextPageSelector:             *cfg.Server.KuperCfg.NextPageSelector,
		},
	}
}
