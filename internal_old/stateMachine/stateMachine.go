package stateMachine

type Statement string

const (
	Default       Statement = "default"
	SendingCookie Statement = "sendingCookie"
	ChattingAI    Statement = "chattingAI"
)

func (s Statement) String() string {
	return string(s)
}

type StateMachine interface {
	GetStatement(uid int64) Statement
	SetStatement(uid int64, statement Statement)
}
