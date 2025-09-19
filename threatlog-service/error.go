package threatlogservice

import "errors"

var (
	ErrNotFound   = errors.New("event not found")
	ErrBadRequest = errors.New("bad request")
)
