package adapter

import (
	"errors"
	"net/http"
	"testing"
)

func TestMapError(t *testing.T) {
	status, resp := MapError(http.StatusTooManyRequests, errors.New("too many requests"))
	if status != http.StatusTooManyRequests {
		t.Fatalf("unexpected status %d", status)
	}
	if resp.Error.Type != "rate_limit_exceeded" {
		t.Fatalf("unexpected type %s", resp.Error.Type)
	}
}
