package http

import (
	"context"
	"net/http"
	"time"

	"github.com/vo1dFl0w/market-parser/internal/transport/http/httpgen"
	"github.com/vo1dFl0w/market-parser/internal/usecase"
	"github.com/vo1dFl0w/market-parser/pkg/logger"
)

type Handler struct {
	logger         logger.Logger
	parserSrv      usecase.ParserService
	requestTimeout time.Duration
}

func NewHandler(logger logger.Logger, parserSrv usecase.ParserService, requestTimeout time.Duration) *Handler {
	return &Handler{logger: logger, parserSrv: parserSrv, requestTimeout: requestTimeout}
}

func (h *Handler) APIV1MarketParserParseGet(ctx context.Context, params httpgen.APIV1MarketParserParseGetParams) (httpgen.APIV1MarketParserParseGetRes, error) {
	h.logger.Info("[Handler] Передаём запрос в parser service...")
	res, err := h.parserSrv.ParseProductsByCategory(ctx, params.Category, params.Address, params.Market)
	if err != nil {
		httpErr := MapError(err)
		h.LogHTTPError(ctx, err, httpErr)
		return httpErr.ToParseErrRes(), nil
	}
	

	resp := make(httpgen.ParseResponse, 0, len(res))
	for _, p := range res {
		resp = append(resp, httpgen.Product{
			Name:  p.Name,
			Link:  p.URL,
			Price: p.Price,
		})
	}

	return &resp, nil
}

func (h *Handler) LogHTTPError(ctx context.Context, err error, httpErr *HTTPError) {
	attrs := []any{
		"error", err,
		"status", httpErr.Status,
		"message", httpErr.Message,
	}

	switch {
	case httpErr.Status >= 500:
		switch httpErr.Status {
		case http.StatusGatewayTimeout:
			h.logger.Error("http_request_failed", append(attrs, "reason", "dependency_timeout")...)
		default:
			h.logger.Error("http_request_failed", append(attrs, "reason", "internal_server_error")...)
		}
	case httpErr.Status >= 400:
		h.logger.Warn("http_request_failed", append(attrs, "reason", "client_error")...)
	}
}