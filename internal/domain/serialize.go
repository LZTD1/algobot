package domain

type SerType int

const (
	GroupType SerType = iota
	UserType
)

type SerializeMessage struct {
	Type SerType
	Data []string
}
