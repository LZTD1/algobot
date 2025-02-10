package appError

import (
	"errors"
)

var ErrNotValid = errors.New("not valid")
var ErrNotFound = errors.New("not found")
var ErrHasNone = errors.New("has no one")
var NotEnoughArgs = errors.New("not enough arguments")
