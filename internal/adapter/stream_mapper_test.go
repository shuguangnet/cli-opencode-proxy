package adapter

import (
	"strings"
	"testing"
)

func TestBuildStreamChunk(t *testing.T) {
	chunk, err := BuildStreamChunk("id1", "model1", 0, "hello", true, "")
	if err != nil {
		t.Fatal(err)
	}
	text := string(chunk)
	if !strings.Contains(text, "data: ") || !strings.Contains(text, "hello") {
		t.Fatalf("unexpected chunk %s", text)
	}
}
