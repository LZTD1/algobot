package backoffice

import (
	"algobot/internal/config"
	"algobot/test/mocks"
	"context"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestRequester(t *testing.T) {
	t.Run("test timeout", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			time.Sleep(1 * time.Second)
			rw.WriteHeader(http.StatusOK)
		}))
		defer server.Close()

		bo := NewBackoffice(&config.Backoffice{
			Retries:         1,
			RetriesTimeout:  time.Millisecond,
			ResponseTimeout: 10 * time.Millisecond,
		}, WithLogger(mocks.NewMockLogger()), WithURL(server.URL))

		req, _ := http.NewRequest(http.MethodGet, server.URL, nil)
		_, err := bo.doReq(req)

		assert.Error(t, err)
		assert.ErrorIs(t, err, context.DeadlineExceeded)
	})
	t.Run("test retry not 200 OK", func(t *testing.T) {
		t.Run("if received 4xx", func(t *testing.T) {
			calls := make([]int, 0)
			server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
				calls = append(calls, 1)
				rw.WriteHeader(http.StatusNotFound)
			}))
			defer server.Close()

			bo := NewBackoffice(&config.Backoffice{
				Retries:         5,
				RetriesTimeout:  time.Millisecond,
				ResponseTimeout: 100 * time.Millisecond,
			}, WithLogger(mocks.NewMockLogger()), WithURL(server.URL))

			req, _ := http.NewRequest(http.MethodGet, server.URL, nil)
			_, err := bo.doReq(req)

			assert.Error(t, err)
			assert.Equal(t, 1, len(calls))
			assert.ErrorIs(t, err, Err4xxStatus)
		})
		t.Run("if recivied not 200 OK", func(t *testing.T) {
			calls := make([]int, 0)
			server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
				calls = append(calls, 1)
				rw.WriteHeader(http.StatusHTTPVersionNotSupported)
			}))
			defer server.Close()

			bo := NewBackoffice(&config.Backoffice{
				Retries:         5,
				RetriesTimeout:  time.Millisecond,
				ResponseTimeout: 100 * time.Millisecond,
			}, WithLogger(mocks.NewMockLogger()), WithURL(server.URL))

			req, _ := http.NewRequest(http.MethodGet, server.URL, nil)
			_, err := bo.doReq(req)

			assert.Error(t, err)
			assert.Equal(t, 5, len(calls))
			assert.ErrorIs(t, err, ErrBadCode)
		})
	})
	t.Run("test retries count", func(t *testing.T) {
		counter := 0
		server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			counter++
			time.Sleep(1 * time.Second)
			rw.WriteHeader(http.StatusOK)
		}))
		defer server.Close()

		bo := NewBackoffice(&config.Backoffice{
			Retries:         5,
			RetriesTimeout:  time.Millisecond,
			ResponseTimeout: 10 * time.Millisecond,
		}, WithLogger(mocks.NewMockLogger()), WithURL(server.URL))

		req, _ := http.NewRequest(http.MethodGet, server.URL, nil)
		_, err := bo.doReq(req)

		assert.Error(t, err)
		assert.Equal(t, 5, counter)
		assert.ErrorIs(t, err, context.DeadlineExceeded)
	})
	t.Run("test retries timeout", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			time.Sleep(1 * time.Second)
			rw.WriteHeader(http.StatusOK)
		}))
		defer server.Close()

		bo := NewBackoffice(&config.Backoffice{
			Retries:         5,
			RetriesTimeout:  5 * time.Millisecond,
			ResponseTimeout: 10 * time.Millisecond,
		}, WithLogger(mocks.NewMockLogger()), WithURL(server.URL))

		req, _ := http.NewRequest(http.MethodGet, server.URL, nil)

		timeStart := time.Now()
		_, err := bo.doReq(req)
		timeEnd := time.Now()

		assert.Error(t, err)
		assert.InDelta(t, 75, timeEnd.Sub(timeStart).Milliseconds(), 5)
	})
	t.Run("happy path", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			rw.WriteHeader(http.StatusOK)
		}))
		defer server.Close()

		bo := NewBackoffice(&config.Backoffice{
			Retries:         5,
			RetriesTimeout:  50 * time.Millisecond,
			ResponseTimeout: 100 * time.Millisecond,
		}, WithLogger(mocks.NewMockLogger()), WithURL(server.URL))
		req, _ := http.NewRequest(http.MethodGet, server.URL, nil)

		resp, err := bo.doReq(req)
		assert.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)
	})
}
