package xerrors

import (
	"errors"
	"net/http"
	"testing"
)

func TestWrapAndCodeLookup(t *testing.T) {
	root := errors.New("db down")
	err := Wrap(root, "DB_ERROR", "query failed")

	if !errors.Is(err, root) {
		t.Fatalf("wrapped error should keep root cause")
	}
	if got := CodeOf(err); got != "DB_ERROR" {
		t.Fatalf("CodeOf() = %q, want %q", got, "DB_ERROR")
	}
	if !IsCode(err, "DB_ERROR") {
		t.Fatalf("IsCode should match wrapped code")
	}
}

func TestHTTPStatusMapping(t *testing.T) {
	notFoundErr := New("NOT_FOUND", "resource missing")
	if got := HTTPStatus(notFoundErr); got != http.StatusNotFound {
		t.Fatalf("HTTPStatus() = %d, want %d", got, http.StatusNotFound)
	}

	customErr := New("RATE_LIMITED", "too many requests")
	RegisterHTTPStatus("RATE_LIMITED", http.StatusTooManyRequests)
	if got := HTTPStatus(customErr); got != http.StatusTooManyRequests {
		t.Fatalf("HTTPStatus() = %d, want %d", got, http.StatusTooManyRequests)
	}
}

func TestCodeOfPlainError(t *testing.T) {
	if got := CodeOf(errors.New("plain")); got != "" {
		t.Fatalf("CodeOf(plain) = %q, want empty", got)
	}
	if got := HTTPStatus(errors.New("plain")); got != http.StatusInternalServerError {
		t.Fatalf("HTTPStatus(plain) = %d, want %d", got, http.StatusInternalServerError)
	}
}
