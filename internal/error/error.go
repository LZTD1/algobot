package appError

import (
	"errors"
)

var ErrNotValid = errors.New("not valid")
var ErrHasNone = errors.New("not found")
