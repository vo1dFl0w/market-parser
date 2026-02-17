package chromium

import (
	"time"

	"github.com/vo1dFl0w/market-parser/internal/config"
)

type ProxyConfig struct {
	IP       string
	Port     string
	Login    string
	Password string
}

type CaptchaSelectors struct {
	KuperSmartCaptcha    string
	KuperCaptchaCheckBox string
}

type Config struct {
	WsURL                 string
	Headless              bool
	TraceMode             bool
	UserAgent             string
	Platoform             string
	Referrer              string
	AcceptLanguage        string
	Proxy                 *ProxyConfig
	CaptchaSelectors      *CaptchaSelectors
	SessionTimeout        time.Duration
	WorkTimeout           time.Duration
	WaitStableDuration    time.Duration
	WaitDOMStableDuration time.Duration
	WaitDOMStableDiff     float64
}

func NewConfigs(cfg *config.Config) *Config {
	proxy := &ProxyConfig{
		IP:       cfg.Browser.Proxy.IP,
		Port:     cfg.Browser.Proxy.Port,
		Login:    cfg.Browser.Proxy.Login,
		Password: cfg.Browser.Proxy.Password,
	}

	captcha := &CaptchaSelectors{
		KuperSmartCaptcha:    *cfg.Server.KuperCfg.SmartCaptchaSelector,
		KuperCaptchaCheckBox: *cfg.Server.KuperCfg.CaptchaCheckBox,
	}

	return &Config{
		WsURL:                 cfg.Browser.WsURL,
		Headless:              cfg.Browser.Headless,
		TraceMode:             cfg.Browser.TraceMode,
		UserAgent:             cfg.Browser.UserAgent,
		Platoform:             cfg.Browser.Platform,
		Referrer:              cfg.Browser.Referer,
		AcceptLanguage:        cfg.Browser.AcceptLanguage,
		Proxy:                 proxy,
		CaptchaSelectors:      captcha,
		SessionTimeout:        cfg.Browser.SessionTimeout,
		WorkTimeout:           cfg.Browser.WorkTimeout,
		WaitStableDuration:    cfg.Browser.WaitStableDuration,
		WaitDOMStableDuration: cfg.Browser.WaitDOMStableDuration,
		WaitDOMStableDiff:     cfg.Browser.WaitDOMStableDiff,
	}
}
