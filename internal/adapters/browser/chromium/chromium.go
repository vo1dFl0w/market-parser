package chromium

import (
	"context"
	"fmt"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
	"github.com/vo1dFl0w/market-parser/internal/config"
	"github.com/vo1dFl0w/market-parser/internal/repository"
)

const (
	//userAgentWindows = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/128.0.0.0 Safari/537.36"
	//platformWindows  = "Win32"

	userAgentLinux = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/128.0.0.0 Safari/537.36"
	platformLinux  = "Linux x86_64"

	// userAgentMacOS = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/128.0.0.0 Safari/537.36"
	// platoformMacOS = "MacIntel"
)

type Chromium struct {
	cfg *Config
}

func NewChromium(cfg *config.Config) *Chromium {
	return &Chromium{cfg: NewConfigs(cfg)}
}

func (ch *Chromium) Connect(ctx context.Context) (*rod.Browser, error) {
	var browser *rod.Browser

	fmt.Println("work timeout", ch.cfg.WorkTimeout)
	// docker-compose
	// must set headless=true in configs/confgi.yaml
	// must set http_addr=market-parser:8080 in configs/config.yaml
	if ch.cfg.Headless {
		l, err := launcher.NewManaged(ch.cfg.WsURL)
		if err != nil {
			return nil, fmt.Errorf("new managed: %w", err)
		}

		l.UserDataDir(ch.cfg.UserDataDir).
			ProfileDir(ch.cfg.ProfileDir).
			Leakless(true).
			KeepUserDataDir().
			HeadlessNew(ch.cfg.Headless).
			Set("user-agent", userAgentLinux).
			Set("disable-blink-features", "AutomationControlled").
			Set("disable-infobars").
			Set("disable-dev-shm-usage").
			Set("disable-background-timer-throttling").
			Set("lang", "ru-RU").
			Set("accept-lang", "ru-RU,ru;q=0.9,en-US;q=0.8,en;q=0.7").
			Set("disable-features", "IsolateOrigins,site-per-process").
			Set("window-size", "1920", "1080")

		if ch.cfg.Proxy.IP != "" || ch.cfg.Proxy.Port != "" || ch.cfg.Proxy.Login != "" || ch.cfg.Proxy.Password != "" {
			proxyAddr := fmt.Sprintf("%s:%s", ch.cfg.Proxy.IP, ch.cfg.Proxy.Port)
			l.Proxy(proxyAddr)
		}

		c, err := l.Client()
		if err != nil {
			return nil, fmt.Errorf("client: %w", err)
		}

		browser = rod.New().Client(c).Trace(true).Timeout(ch.cfg.SessionTimeout).Context(ctx)
		if err := browser.Connect(); err != nil {
			return nil, fmt.Errorf("connect browser: %w", err)
		}

	} else {
		// launch by makefile or go run
		// must set headless=false in configs/confgi.yaml
		// must set http_addr=localhost:8080 in configs/config.yaml
		l := launcher.New().
			UserDataDir(ch.cfg.UserDataDirLocal).
			ProfileDir(ch.cfg.ProfileDir).
			Leakless(true).
			HeadlessNew(ch.cfg.Headless).
			Set("disable-blink-features", "AutomationControlled").
			Set("disable-infobars").
			Set("disable-dev-shm-usage").
			Set("disable-background-timer-throttling").
			Set("lang", "ru-RU").
			Set("accept-lang", "ru-RU,ru;q=0.9,en-US;q=0.8,en;q=0.7").
			Set("disable-features", "IsolateOrigins,site-per-process").
			Set("window-size", "1920", "1080")

		if ch.cfg.Proxy.IP != "" || ch.cfg.Proxy.Port != "" || ch.cfg.Proxy.Login != "" || ch.cfg.Proxy.Password != "" {
			proxyAddr := fmt.Sprintf("%s:%s", ch.cfg.Proxy.IP, ch.cfg.Proxy.Port)
			l.Proxy(proxyAddr)
		}

		url, err := l.Launch()
		if err != nil {
			return nil, fmt.Errorf("launch: %w", err)
		}

		// set trace=true to get more logs
		browser = rod.New().ControlURL(url).Trace(true).Timeout(ch.cfg.SessionTimeout)

		if err := browser.Connect(); err != nil {
			return nil, fmt.Errorf("connect browser: %w", err)
		}
	}

	return browser, nil
}

func (ch *Chromium) NewPage(ctx context.Context) (repository.Page, error) {
	browser, err := ch.Connect(ctx)
	if err != nil {
		return nil, fmt.Errorf("connect browser: %w", err)
	}

	if ch.cfg.Proxy.IP != "" || ch.cfg.Proxy.Port != "" || ch.cfg.Proxy.Login != "" || ch.cfg.Proxy.Password != "" {
		go browser.HandleAuth(ch.cfg.Proxy.Login, ch.cfg.Proxy.Password)()
	}

	page, err := browser.Page(proto.TargetCreateTarget{})
	if err != nil {
		return nil, fmt.Errorf("page: %w", err)
	}

	if ch.cfg.Headless {
		if err := page.SetUserAgent(&proto.NetworkSetUserAgentOverride{
			UserAgent:      userAgentLinux,
			AcceptLanguage: "ru-RU,ru;q=0.9,en-US;q=0.8,en;q=0.7",
			Platform:       platformLinux,
		}); err != nil {
			return nil, fmt.Errorf("set user agent: %w", err)
		}
	} else {
		if err := page.SetUserAgent(&proto.NetworkSetUserAgentOverride{
			UserAgent:      userAgentLinux,
			AcceptLanguage: "ru-RU,ru;q=0.9,en-US;q=0.8,en;q=0.7",
			Platform:       platformLinux,
		}); err != nil {
			return nil, fmt.Errorf("set user agent: %w", err)
		}
	}

	if err := page.SetViewport(&proto.EmulationSetDeviceMetricsOverride{
		Width:             1920,
		Height:            1080,
		DeviceScaleFactor: 1,
		Mobile:            false,
	}); err != nil {
		fmt.Printf("set view port: %s\n", err.Error())
	}

	_, err = proto.PageNavigate{
		URL: ch.cfg.Referrer,
	}.Call(page)
	if err != nil {
		return nil, fmt.Errorf("page navigate call: %w", err)
	}
	if err := page.WaitLoad(); err != nil {
		return nil, fmt.Errorf("wait load: %w", err)
	}

	return &rodPage{page: page, browser: browser, cfg: ch.cfg}, nil
}
