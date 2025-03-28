package fsm

type State int

const (
	Default State = iota
	SendingCookie
	ChattingAI
)
