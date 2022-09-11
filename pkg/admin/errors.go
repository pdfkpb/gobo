package admin

import "errors"

var (
	ErrNotAdmin = errors.New("user not authorized to take that action")
)
