package http

import "errors"

var (
	ErrIncorrectResponseStatusCode = errors.New("incorrect response status code")

	ErrURLIsNotAllowed = errors.New("url is not allowed")
	ErrURLIsBlocked    = errors.New("url is blocked")
)
