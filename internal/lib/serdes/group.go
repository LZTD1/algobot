package serdes

import "errors"

var (
	ErrUnrecognized = errors.New("unrecognized sertype")
)

type SerType int

const (
	GroupType SerType = iota
	UserType
)
