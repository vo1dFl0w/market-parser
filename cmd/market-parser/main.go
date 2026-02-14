package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/vo1dFl0w/market-parser/internal/adapters/browser/chromium"
	"github.com/vo1dFl0w/market-parser/internal/adapters/parsers"
	"github.com/vo1dFl0w/market-parser/internal/config"
	ht "github.com/vo1dFl0w/market-parser/internal/transport/http"
	"github.com/vo1dFl0w/market-parser/internal/transport/http/httpgen"
	"github.com/vo1dFl0w/market-parser/internal/usecase"
	"github.com/vo1dFl0w/market-parser/pkg/logger"
)

// 1. зайти на главную страницу купер
// 2. ввести адрес доставки
// 3. параллельно открыть несколько страниц с магазинами (лента, магнит, пятёрочка и так далее)
// 4. на каждой странице нажать на категорию справа (категорию получаем в url get-запроса)
// 5. параллельное открытие нескольких страниц для каждого магазина (page=1, page=2 и так далее)
// 6. парсинг страницы из полученного json

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := run(ctx); err != nil {
		log.Println(ctx, "startup", "err", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	loggerCfg := logger.NewLoggerConfig(cfg.Server.Env, cfg.Options.LoggerTimeFormat)
	logger := logger.LoadLogger(loggerCfg)

	chromiumRepo := chromium.NewChromium(cfg, logger)
	browserRepo := chromium.NewBrowser(chromiumRepo)
	kuperParser := parsers.NewKuperParser(cfg, logger, browserRepo.Chromium())
	parserSrv := usecase.NewParserService(kuperParser)

	handler := ht.NewHandler(logger, parserSrv, cfg.Server.RequestTimeout)

	srv, err := httpgen.NewServer(handler)
	if err != nil {
		return fmt.Errorf("new server: %w", err)
	}

	withMiddlewares := handler.CORSMiddleware(handler.RequestTimeoutMiddleware(handler.LoggerMiddleware(srv)))

	httpServer := http.Server{
		Addr:    cfg.Server.HTTPAddr,
		Handler: withMiddlewares,
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)
	serverErr := make(chan error, 1)

	go func() {
		logger.Info("server started", "host", httpServer.Addr)
		if err := httpServer.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			serverErr <- err
		} else {
			serverErr <- nil
		}
	}()

	select {
	case e := <-serverErr:
		return fmt.Errorf("server error: %w", e)
	case s := <-sig:
		logger.Info("initialization gracefull shutdown", "signal", s)

		shutdownCtx, cancel := context.WithTimeout(ctx, cfg.Server.ShutdownTimeout)
		defer cancel()

		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("shutdown http server: %w", err)
		}

		logger.Info("server gracefully stopped")
		return nil
	}
}
