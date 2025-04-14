package test

import (
	"algobot/internal/lib/fsm"
	"algobot/internal/lib/fsm/memory"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Memory(t *testing.T) {
	t.Run("Set state", func(t *testing.T) {
		sm := memory.New()

		sm.SetState(1, fsm.Default)
		got := sm.State(1)
		assert.Equal(t, fsm.Default, got)

		sm.SetState(1, fsm.ChattingAI)
		got = sm.State(1)
		assert.Equal(t, fsm.ChattingAI, got)
	})
	t.Run("Get state if not exists", func(t *testing.T) {
		sm := memory.New()

		got := sm.State(1)
		assert.Equal(t, fsm.Default, got)
	})
}
