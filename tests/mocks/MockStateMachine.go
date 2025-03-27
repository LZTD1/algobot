package mocks

import (
	"algobot/internal_old/stateMachine"
	"strconv"
)

type MockStateMachine struct {
	Calls   []string
	Current stateMachine.Statement
}

func (m *MockStateMachine) GetStatement(uid int64) stateMachine.Statement {
	m.Calls = append(m.Calls, "GetStatement")
	m.Calls = append(m.Calls, strconv.FormatInt(uid, 10))
	return m.Current
}

func (m *MockStateMachine) SetStatement(uid int64, statement stateMachine.Statement) {
	m.Calls = append(m.Calls, "SetStatement")
	m.Calls = append(m.Calls, strconv.FormatInt(uid, 10))
	m.Calls = append(m.Calls, statement.String())
	m.Current = statement
}
