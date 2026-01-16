package logger

import (
	"fmt"
	"testing"
)

func TestLogger(t *testing.T) {
	// 1. Setup the Mock
	mock := &MockWriter{}

	// 2. Initialize Logger with our pattern
	l := NewLogger(
		WithPrefix("TEST"),
		WithWriter(mock), // Injecting the mock!
	)

	fmt.Println(l.prefix)

	// 3. Run the action
	testMsg := "Hello Go"
	err := l.Log(testMsg)

	// 4. Verify results
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	expected := "[TEST] Hello Go"
	if mock.LastMessage != expected {
		t.Errorf("Expected %s, but got %s", expected, mock.LastMessage)
	}
}
