package pipeline

import (
	"context"
	"strings"
	"testing"
	"time"
)

type UpperProcessor struct{}

func (up UpperProcessor) Process(s string) (string, error) {
	// time.Sleep(50 * time.Millisecond) // Simulate heavy work
	return strings.ToUpper(s), nil
}

func TestRunEngine(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	items := []string{"a", "b", "c", "d", "e", "f", "g"}
	p := UpperProcessor{}

	// 2. Run with concurrency of 3
	results := RunEngine(ctx, items, p, 3)

	if len(results) != len(items) {
		t.Errorf("Expected %d results, got %d", len(items), len(results))
	}

	for _, result := range results {
		if result.Err != nil {
			t.Errorf("Unexpected error: %v", result.Err)
		}
	}
}
