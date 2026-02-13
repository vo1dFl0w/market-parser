package http

import (
	"errors"
	"net/http"

	"github.com/vo1dFl0w/market-parser/internal/domain"
	"github.com/vo1dFl0w/market-parser/internal/transport/http/httpgen"
)

const (
	StatusClientClosedRequest = 499
)

var (
	ErrBadRequest          = errors.New("bad request")
	ErrClientClosedRequest = errors.New("client closed request")
	ErrGatewayTimeout      = errors.New("gateway timeout")
	ErrInternalServerError = errors.New("internal server error")
)

type HTTPError struct {
	Message string
	Status  int
}

func (e *HTTPError) Error() string {
	return e.Message
}

func (e *HTTPError) ToParseErrRes() httpgen.APIV1MarketParserParseGetRes {
	switch e.Status {
	case http.StatusBadRequest:
		return &httpgen.APIV1MarketParserParseGetBadRequest{Message: e.Message, Status: e.Status}
	case StatusClientClosedRequest:
		return &httpgen.APIV1MarketParserParseGetCode499{Message: e.Message, Status: e.Status}
	case http.StatusGatewayTimeout:
		return &httpgen.APIV1MarketParserParseGetGatewayTimeout{Message: e.Message, Status: e.Status}
	default:
		return &httpgen.APIV1MarketParserParseGetInternalServerError{Message: e.Message, Status: e.Status}
	}
}

func MapError(err error) *HTTPError {
	switch {
	case errors.Is(err, domain.ErrEmptyCategory):
		return &HTTPError{Message: ErrBadRequest.Error(), Status: http.StatusBadRequest}
	case errors.Is(err, domain.ErrEmptyAddress):
		return &HTTPError{Message: ErrBadRequest.Error(), Status: http.StatusBadRequest}
	case errors.Is(err, domain.ErrEmptyMarket):
		return &HTTPError{Message: ErrBadRequest.Error(), Status: http.StatusBadRequest}
	case errors.Is(err, domain.ErrClientClosedRequest):
		return &HTTPError{Message: ErrClientClosedRequest.Error(), Status: StatusClientClosedRequest}
	case errors.Is(err, domain.ErrGatewayTimeout):
		return &HTTPError{Message: ErrGatewayTimeout.Error(), Status: http.StatusGatewayTimeout}
	default:
		return &HTTPError{Message: ErrInternalServerError.Error(), Status: http.StatusInternalServerError}
	}
}