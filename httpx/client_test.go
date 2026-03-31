package httpx

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"
)

func TestDoRetriesOnServerError(t *testing.T) {
	var count int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n := atomic.AddInt32(&count, 1)
		if n < 3 {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	}))
	defer srv.Close()

	c := NewClient(Options{Timeout: time.Second, Retry: 2, RetryWait: time.Millisecond * 10})
	req, err := http.NewRequest(http.MethodGet, srv.URL, nil)
	if err != nil {
		t.Fatal(err)
	}
	resp, err := c.Do(context.Background(), req)
	if err != nil {
		t.Fatalf("Do() error = %v", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", resp.StatusCode)
	}
	if got := atomic.LoadInt32(&count); got != 3 {
		t.Fatalf("request count = %d, want 3", got)
	}
}

func TestDoDoesNotRetryOnClientError(t *testing.T) {
	var count int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&count, 1)
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer srv.Close()

	c := NewClient(Options{Timeout: time.Second, Retry: 3, RetryWait: time.Millisecond * 10})
	req, err := http.NewRequest(http.MethodGet, srv.URL, nil)
	if err != nil {
		t.Fatal(err)
	}
	resp, err := c.Do(context.Background(), req)
	if err != nil {
		t.Fatalf("Do() error = %v", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", resp.StatusCode)
	}
	if got := atomic.LoadInt32(&count); got != 1 {
		t.Fatalf("request count = %d, want 1", got)
	}
}

func TestJSONHelper(t *testing.T) {
	type inPayload struct {
		Name string `json:"name"`
	}
	type outPayload struct {
		Echo string `json:"echo"`
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var in inPayload
		if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		_ = json.NewEncoder(w).Encode(outPayload{Echo: in.Name})
	}))
	defer srv.Close()

	c := NewClient(Options{Timeout: time.Second, Retry: 1, RetryWait: time.Millisecond * 10})
	var out outPayload
	if err := c.JSON(context.Background(), http.MethodPost, srv.URL, inPayload{Name: "alice"}, nil, &out); err != nil {
		t.Fatalf("JSON() error = %v", err)
	}
	if out.Echo != "alice" {
		t.Fatalf("out.Echo = %q, want alice", out.Echo)
	}
}
