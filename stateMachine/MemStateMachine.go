package stateMachine

type MemStateMachine struct {
	statements map[int64]Statement
}

func NewMemStateMachine() *MemStateMachine {
	return &MemStateMachine{
		statements: make(map[int64]Statement),
	}
}
func (m *MemStateMachine) SetStatement(uid int64, statement Statement) {
	m.statements[uid] = statement
}
func (m *MemStateMachine) GetStatement(uid int64) Statement {
	v, ok := m.statements[uid]
	if ok {
		return v
	}
	m.statements[uid] = Default
	return m.statements[uid]
}
