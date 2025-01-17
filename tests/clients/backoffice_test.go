package test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"tgbot/internal/clients"
	"time"
)

func TestBackoffice(t *testing.T) {
	boSettings := clients.BackofficeSetting{
		Retry:        3,
		Timeout:      100 * time.Millisecond,
		RetryTimeout: 50 * time.Millisecond,
	}

	t.Run("GetKidsNamesByGroup", func(t *testing.T) {
		t.Run("401 | Unauthorized", func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("[]"))
			}))
			defer ts.Close()

			bo := clients.NewBackoffice(ts.URL, boSettings)
			_, err := bo.GetKidsNamesByGroup("", "")
			assertError(t, err, 401, "[]")
		})
		t.Run("Servers returns error", func(t *testing.T) {
			var calls []string
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
				calls = append(calls, "GET")
			}))
			defer ts.Close()

			bo := clients.NewBackoffice(ts.URL, boSettings)
			_, err := bo.GetKidsNamesByGroup("", "")
			assertError(t, err, 500, "")
			if len(calls) != 3 {
				t.Fatalf("expected 3 calls, got %d", len(calls))
			}
		})
	})

}

func assertError(t *testing.T, err *clients.ClientError, i int, s string) {
	t.Helper()

	if err == nil {
		t.Fatal("Expected error, got nil")
	}
	if err.Code != i {
		t.Fatalf("Expected code %d, got %d", i, err.Code)
	}
	if err.Message != s {
		t.Fatalf("Expected message %s, got %s", s, err.Message)
	}
}
