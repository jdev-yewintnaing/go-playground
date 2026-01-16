package watcher

import (
	"context"
	"net/http"
	"time"
)

type CheckResult struct {
	URL     string
	Latency time.Duration
	Status  int
	Err     error
}

type StatusChecker interface {
	Check(ctx context.Context, url string) CheckResult
}

type MockChecker struct {
	Delay time.Duration
}

func (m MockChecker) Check(ctx context.Context, url string) CheckResult {
	select {
	case <-ctx.Done():
		return CheckResult{
			URL: url,
			Err: ctx.Err(),
		}
	case <-time.After(m.Delay):
		return CheckResult{URL: url, Status: http.StatusOK}
	}
}

type HTTPChecker struct{}

func (h HTTPChecker) Check(ctx context.Context, url string) CheckResult {
	start := time.Now()

	// Create a request that respects the context
	req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return CheckResult{URL: url, Err: err}
	}
	defer resp.Body.Close()

	return CheckResult{
		URL:     url,
		Latency: time.Since(start),
		Status:  resp.StatusCode,
	}
}

func Watch(ctx context.Context, urls []string, checker StatusChecker) []CheckResult {
	chanResults := make(chan CheckResult, len(urls))

	for _, url := range urls {
		go func(url string) {
			chanResults <- checker.Check(ctx, url)
		}(url)
	}

	var results []CheckResult

	for range urls {
		select {
		case result := <-chanResults:
			results = append(results, result)
		case <-ctx.Done():
			return results
		}
	}

	return results
}
