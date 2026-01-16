package watcher

import (
	"context"
	"testing"
	"time"
)

func TestWatch(t *testing.T) {
	httpChecker := &MockChecker{
		Delay: 10 * time.Millisecond,
	}

	urls := []string{
		"http://example.com",
		"https//google.com",
		"https//yewintnaing.dev",
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	results := Watch(ctx, urls, httpChecker)

	if len(results) == len(urls) {
		t.Logf("Expected %d results, got %d", len(urls), len(results))
	}
}

func TestWatch_Timeout(t *testing.T) {
	httpChecker := &MockChecker{
		Delay: 10 * time.Second,
	}

	urls := []string{
		"http://example.com",
		"https//google.com",
		"https//yewintnaing.dev",
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	results := Watch(ctx, urls, httpChecker)

	if len(results) < len(urls) {
		t.Logf("Success: Watch returned early with %d results due to timeout", len(results))
	}
}

func TestWatchTableDriven(t *testing.T) {
	// 1. Define the "Table" of scenarios
	tests := []struct {
		name          string
		timeout       time.Duration
		checkerDelay  time.Duration
		expectedCount int
	}{
		{
			name:          "Full Success",
			timeout:       2 * time.Second,
			checkerDelay:  100 * time.Millisecond,
			expectedCount: 2,
		},
		{
			name:          "Complete Timeout",
			timeout:       100 * time.Millisecond,
			checkerDelay:  2 * time.Second,
			expectedCount: 0,
		},
	}

	urls := []string{"site1.com", "site2.com"}

	// 2. Iterate over the table
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup context and mock for this specific case
			ctx, cancel := context.WithTimeout(context.Background(), tt.timeout)
			defer cancel()

			mock := &MockChecker{Delay: tt.checkerDelay}

			// Execute
			results := Watch(ctx, urls, mock)

			// Verify
			if len(results) != tt.expectedCount {
				t.Errorf("got %d results, want %d", len(results), tt.expectedCount)
			}
		})
	}
}
