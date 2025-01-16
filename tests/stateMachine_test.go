package tests_test

import (
	"testing"
	stateMachine2 "tgbot/internal/stateMachine"
)

func TestState(t *testing.T) {
	t.Run("Set state", func(t *testing.T) {
		sm := stateMachine2.NewMemory()

		sm.SetStatement(1, stateMachine2.Default)
		got := sm.GetStatement(1)
		if got != stateMachine2.Default {
			t.Fatalf("Wanted state %s, but got %s", stateMachine2.Default, got)
		}
	})
	t.Run("Get state if not exists", func(t *testing.T) {
		sm := stateMachine2.NewMemory()

		got := sm.GetStatement(1)
		if got != stateMachine2.Default {
			t.Fatalf("Wanted state %s, but got %s", stateMachine2.Default, got)
		}
	})
}
