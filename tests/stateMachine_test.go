package tests_test

import (
	"testing"
	"tgbot/stateMachine"
)

func TestState(t *testing.T) {
	t.Run("Set state", func(t *testing.T) {
		sm := stateMachine.NewMemStateMachine()

		sm.SetStatement(1, stateMachine.Default)
		got := sm.GetStatement(1)
		if got != stateMachine.Default {
			t.Fatalf("Wanted state %s, but got %s", stateMachine.Default, got)
		}
	})
	t.Run("Get state if not exists", func(t *testing.T) {
		sm := stateMachine.NewMemStateMachine()

		got := sm.GetStatement(1)
		if got != stateMachine.Default {
			t.Fatalf("Wanted state %s, but got %s", stateMachine.Default, got)
		}
	})
}
