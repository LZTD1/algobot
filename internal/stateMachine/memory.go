package stateMachine

type Memory struct {
	statements map[int64]Statement
}

func NewMemory() *Memory {
	return &Memory{
		statements: make(map[int64]Statement),
	}
}
func (m *Memory) SetStatement(uid int64, statement Statement) {
	m.statements[uid] = statement
}
func (m *Memory) GetStatement(uid int64) Statement {
	v, ok := m.statements[uid]
	if ok {
		return v
	}
	m.statements[uid] = Default
	return m.statements[uid]
}
