package config

import (
	"fmt"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	//"github.com/joho/godotenv"
)

type Config struct {
	Server  ServerConfig  `yaml:"server"`
	Browser BrowserConfig `yaml:"browser"`
	Options OptionsConfig `yaml:"options"`
}

type ServerConfig struct {
	Env             string        `yaml:"env" env:"SERVER_ENV" env-required:"true"`
	HTTPAddr        string        `yaml:"http_addr" env:"SERVER_HTTP_ADDR" env-required:"true"`
	KuperCfg        KuperConfig   `yaml:"kuper_config"`
	RequestTimeout  time.Duration `yaml:"request_timeout" env:"SERVER_REQUEST_TIMEOUT" env-default:"180000ms"`
	ShutdownTimeout time.Duration `yaml:"shutdown_timeout" env:"SERVER_SHUTDOWN_TIMEOUT" env-default:"15000ms"`
}

type BrowserConfig struct {
	WsURL                   string        `yaml:"ws_url" env:"BROWSER_WS_URL" env-required:"true"`
	Headless                bool          `yaml:"headless"`
	Proxy                   ProxyConfig   `yaml:"proxy"`
	Referer                 string        `yaml:"referer" env-default:"https://google.com"`
	AcceptLanguage          string        `yaml:"accept_language" env-default:"ru-RU,ru;q=0.9"`
	HeadlessMode            bool          `yaml:"headless_mode" env:"BROWSER_HEADLESS_MODE" env-default:"true"`
	SessionTimeout          time.Duration `yaml:"session_timeout" env:"BROWSER_SESSION_TIMEOUT" env-default:"180000ms"`
	HumanLikeMode           bool          `yaml:"human_like_mode"`
	TestParserMode          bool          `yaml:"test_parser_mode"`
	WorkTimeout             time.Duration `yaml:"work_timeout" env-default:"5000ms"`
	WaitStableDuration      time.Duration `yaml:"wait_stable_duration" env:"BROWSER_WAIT_STABLE_DURATION" env-default:"500ms"`
	WaitDOMStableDuration   time.Duration `yaml:"wait_dom_stable_duration" env:"BROWSER_WAIT_DOM_STABLE_DURATION" env-default:"300ms"`
	WaitDOMStableDiff       float64       `yaml:"wait_dom_stable_diff" env:"BROWSER_WAIT_DOM_STABLE_DIFF" env-default:"0.85"`
	WaitRequestIdleDuration time.Duration `yaml:"wait_request_idle_duration" env:"BROWSER_WAIT_IDLE_DURATION" env-default:"500ms"`
}

type OptionsConfig struct {
	LoggerTimeFormat string `yaml:"logger_time_format" env:"OPTIONS_LOGGER_TIME_FORMAT" env-default:"02-01-2006 15:04:05"`
}

type KuperConfig struct {
	BaseURL                      *string `yaml:"base_url" env-required:"true"`
	ApiProductsPath              *string `yaml:"api_products_path" env-required:"true"`
	CaptchaCheckBox              *string `yaml:"captcha_check_box" env-required:"true"`
	SmartCaptchaSelector         *string `yaml:"smart_captcha_selector" env-required:"true"`
	CurrentAddressSelector       *string `yaml:"current_address_selector" env-required:"true"`
	AddressButtonSelector        *string `yaml:"address_button_selector" env-required:"true"`
	AddressCheckAttributeValue   *string `yaml:"address_check_attribute_value" env-required:"true"`
	AddressInputSelector         *string `yaml:"address_input_selector" env-required:"true"`
	AddressInputDropDownSelector *string `yaml:"address_input_drop_down_selector" env-required:"true"`
	AddressSaveButtonSelector    *string `yaml:"address_save_button_selector" env-required:"true"`
	MarketSelector               *string `yaml:"market_selector" env-required:"true"`
	AllProdsSelector             *string `yaml:"all_prods_selector" env-required:"true"`
	LastPageSelector             *string `yaml:"last_page_selector" env-required:"true"`
	LastPageText                 *string `yaml:"last_page_text" env-required:"true"`
	NextPageSelector             *string `yaml:"next_page_selector" env-required:"true"`
}

type ProxyConfig struct {
	IP       string `yaml:"ip" env:"BROWSER_PROXY_IP"`
	Port     string `yaml:"port" env:"BROWSER_PROXY_PORT"`
	Login    string `yaml:"login" env:"BROWSER_PROXY_LOGIN"`
	Password string `yaml:"password" env:"BROWSER_PROXY_PASSWORD"`
}

func LoadConfig() (*Config, error) {
	var cfg Config

	/*
		if err := godotenv.Load(".env"); err != nil {
			return nil, fmt.Errorf("load env-file: %w", err)
		}
	*/

	path := os.Getenv("CONFIG_PATH")
	if path == "" {
		return nil, fmt.Errorf("CONFIG_PATH not set")
	}

	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	return &cfg, nil
}
