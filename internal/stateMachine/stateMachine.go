package stateMachine

type Statement string

const (
	Default       Statement = "default"
	SendingCookie Statement = "sendingCookie"
)

type StateMachine interface {
	GetStatement(uid int64) Statement
	SetStatement(uid int64, statement Statement)
}
