package domain

import "errors"

var (
	ErrEmptyCategory       = errors.New("empty category")
	ErrEmptyAddress        = errors.New("empty address")
	ErrEmptyMarket         = errors.New("empty market")
	ErrGatewayTimeout      = errors.New("gateway timeout")
	ErrClientClosedRequest = errors.New("client closed request")
)
