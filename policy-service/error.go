package policyservice

import "errors"

var (
	ErrUnauthorized = errors.New("user is blocked by policy")
)
