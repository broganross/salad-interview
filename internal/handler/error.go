package handler

import "errors"

var (
	ErrMaxRetry = errors.New("maximum retries executed")
)
