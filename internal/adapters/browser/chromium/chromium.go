package chromium

import (
	"context"
	"fmt"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
	"github.com/vo1dFl0w/market-parser/internal/config"
	"github.com/vo1dFl0w/market-parser/internal/repository"
	"github.com/vo1dFl0w/market-parser/pkg/logger"
)

type Chromium struct {
	cfg    *Config
	logger logger.Logger
}

func NewChromium(cfg *config.Config, logger logger.Logger) *Chromium {
	return &Chromium{cfg: NewConfigs(cfg), logger: logger}
}

func (ch *Chromium) Connect(ctx context.Context) (*rod.Browser, error) {
	var browser *rod.Browser
	// docker-compose
	// must set headless=true in configs/confgi.yaml
	// must set http_addr=market-parser:8080 in configs/config.yaml
	if ch.cfg.Headless {
		l, err := launcher.NewManaged(ch.cfg.WsURL)
		if err != nil {
			return nil, fmt.Errorf("new managed: %w", err)
		}

		l.HeadlessNew(ch.cfg.Headless).
			Set("user-agent", ch.cfg.UserAgent).
			Set("disable-blink-features", "AutomationControlled").
			Set("disable-infobars").
			Set("disable-dev-shm-usage").
			Set("disable-background-timer-throttling").
			Set("lang", "ru-RU").
			Set("accept-lang", ch.cfg.AcceptLanguage).
			Set("disable-features", "IsolateOrigins,site-per-process").
			Set("window-size", "1920", "1080")

		if ch.cfg.Proxy.IP != "" && ch.cfg.Proxy.Port != "" {
			proxyAddr := fmt.Sprintf("%s:%s", ch.cfg.Proxy.IP, ch.cfg.Proxy.Port)
			l.Proxy(proxyAddr)
		}

		c, err := l.Client()
		if err != nil {
			return nil, fmt.Errorf("client: %w", err)
		}

		browser = rod.New().Client(c).Trace(ch.cfg.TraceMode).Timeout(ch.cfg.SessionTimeout).Context(ctx)
		if err := browser.Connect(); err != nil {
			return nil, fmt.Errorf("connect browser: %w", err)
		}

	} else {
		// launch by makefile or go run
		// must set headless=false in configs/confgi.yaml
		// must set http_addr=localhost:8080 in configs/config.yaml
		l := launcher.New().
			HeadlessNew(ch.cfg.Headless).
			Set("user-agent", ch.cfg.UserAgent).
			Set("disable-blink-features", "AutomationControlled").
			Set("disable-infobars").
			Set("disable-dev-shm-usage").
			Set("disable-background-timer-throttling").
			Set("lang", "ru-RU").
			Set("accept-lang", ch.cfg.AcceptLanguage).
			Set("disable-features", "IsolateOrigins,site-per-process").
			Set("window-size", "1920", "1080")

		if ch.cfg.Proxy.IP != "" && ch.cfg.Proxy.Port != "" {
			proxyAddr := fmt.Sprintf("%s:%s", ch.cfg.Proxy.IP, ch.cfg.Proxy.Port)
			l.Proxy(proxyAddr)
		}

		url, err := l.Launch()
		if err != nil {
			return nil, fmt.Errorf("launch: %w", err)
		}

		// set trace=true to get more logs
		browser = rod.New().ControlURL(url).Trace(ch.cfg.TraceMode).Timeout(ch.cfg.SessionTimeout)

		if err := browser.Connect(); err != nil {
			return nil, fmt.Errorf("connect browser: %w", err)
		}
	}

	return browser, nil
}

func (ch *Chromium) NewPage(ctx context.Context, marketURL string) (repository.Page, error) {
	browser, err := ch.Connect(ctx)
	if err != nil {
		return nil, fmt.Errorf("connect browser: %w", err)
	}

	page, err := browser.Page(proto.TargetCreateTarget{})
	if err != nil {
		return nil, fmt.Errorf("page: %w", err)
	}

	if ch.cfg.Proxy.IP != "" || ch.cfg.Proxy.Port != "" || ch.cfg.Proxy.Login != "" || ch.cfg.Proxy.Password != "" {
		go browser.HandleAuth(ch.cfg.Proxy.Login, ch.cfg.Proxy.Password)()
	}

	if err := page.SetUserAgent(&proto.NetworkSetUserAgentOverride{
		UserAgent:      ch.cfg.UserAgent,
		AcceptLanguage: ch.cfg.AcceptLanguage,
		Platform:       ch.cfg.Platoform,
	}); err != nil {
		return nil, fmt.Errorf("set user agent: %w", err)
	}

	if err := page.SetViewport(&proto.EmulationSetDeviceMetricsOverride{
		Width:             1920,
		Height:            1080,
		DeviceScaleFactor: 1,
		Mobile:            false,
	}); err != nil {
		return nil, fmt.Errorf("set view port: %s\n", err)
	}

	_, err = proto.PageNavigate{
		URL:      marketURL,
		Referrer: ch.cfg.Referrer,
	}.Call(page)
	if err != nil {
		return nil, fmt.Errorf("page navigate call: %w", err)
	}
	if err := page.WaitLoad(); err != nil {
		return nil, fmt.Errorf("wait load: %w", err)
	}

	return &rodPage{page: page, browser: browser, cfg: ch.cfg}, nil
}
