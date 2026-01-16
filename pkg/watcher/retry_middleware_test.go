package watcher

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"
)

type sequenceChecker struct {
	mu    sync.Mutex
	seq   map[string][]CheckResult
	calls map[string]int
}

func newSequenceChecker(seq map[string][]CheckResult) *sequenceChecker {
	return &sequenceChecker{
		seq:   seq,
		calls: make(map[string]int),
	}
}

func (s *sequenceChecker) Check(ctx context.Context, url string) CheckResult {
	s.mu.Lock()
	s.calls[url]++
	n := s.calls[url]
	results := s.seq[url]
	s.mu.Unlock()

	if len(results) == 0 {
		return CheckResult{URL: url, Err: errors.New("no sequence configured")}
	}

	idx := n - 1
	if idx >= len(results) {
		idx = len(results) - 1
	}

	r := results[idx]
	if r.URL == "" {
		r.URL = url
	}
	return r
}

func (s *sequenceChecker) Calls(url string) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.calls[url]
}

func TestRetryMiddleware_Check_TableDriven(t *testing.T) {
	type testCase struct {
		name      string
		retries   int
		url       string
		sequence  []CheckResult
		wantCalls int

		wantStatus int
		wantErr    bool
		wantErrMsg string // optional exact error check
	}

	tests := []testCase{
		{
			name:    "success on first attempt",
			retries: 3,
			url:     "site1.com",
			sequence: []CheckResult{
				{Status: 200, Err: nil},
			},
			wantCalls:  1,
			wantStatus: 200,
			wantErr:    false,
		},
		{
			name:    "fails twice then succeeds",
			retries: 5,
			url:     "site2.com",
			sequence: []CheckResult{
				{Status: 0, Err: errors.New("timeout")},
				{Status: 0, Err: errors.New("timeout")},
				{Status: 200, Err: nil},
			},
			wantCalls:  3,
			wantStatus: 200,
			wantErr:    false,
		},
		{
			name:    "all attempts fail returns last failure",
			retries: 3,
			url:     "site3.com",
			sequence: []CheckResult{
				{Status: 0, Err: errors.New("e1")},
				{Status: 0, Err: errors.New("e2")},
				{Status: 0, Err: errors.New("e3")},
			},
			wantCalls:  3,
			wantStatus: 0,
			wantErr:    true,
			wantErrMsg: "e3",
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()

			mock := newSequenceChecker(map[string][]CheckResult{
				tc.url: tc.sequence,
			})

			mw := RetryMiddleware{
				next:    mock,
				retries: tc.retries,
			}

			got := mw.Check(ctx, tc.url)

			if calls := mock.Calls(tc.url); calls != tc.wantCalls {
				t.Fatalf("calls(%q)=%d want=%d", tc.url, calls, tc.wantCalls)
			}

			if got.Status != tc.wantStatus {
				t.Fatalf("status=%d want=%d", got.Status, tc.wantStatus)
			}

			if (got.Err != nil) != tc.wantErr {
				t.Fatalf("errPresent=%v want=%v err=%v", got.Err != nil, tc.wantErr, got.Err)
			}

			if tc.wantErrMsg != "" {
				if got.Err == nil || got.Err.Error() != tc.wantErrMsg {
					t.Fatalf("err=%v want message=%q", got.Err, tc.wantErrMsg)
				}
			}
		})
	}
}
