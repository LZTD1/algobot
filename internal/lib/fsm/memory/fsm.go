package memory

import (
	"algobot/internal/lib/fsm"
	"sync"
)

type Memory struct {
	mu         sync.Mutex
	statements map[int64]fsm.State
}

func New() *Memory {
	return &Memory{
		statements: make(map[int64]fsm.State),
	}
}
func (m *Memory) SetState(uid int64, statement fsm.State) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.statements[uid] = statement
}
func (m *Memory) State(uid int64) fsm.State {
	m.mu.Lock()
	defer m.mu.Unlock()

	v, ok := m.statements[uid]
	if ok {
		return v
	}
	m.statements[uid] = fsm.Default
	return m.statements[uid]
}
