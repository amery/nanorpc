package client

import (
	"sync"
	"testing"
)

// ConcurrentTestHelper helps with concurrent test execution
type ConcurrentTestHelper struct {
	t           *testing.T
	results     []any
	errors      []error
	wg          sync.WaitGroup
	mutex       sync.Mutex
	numRoutines int
}

// NewConcurrentTestHelper creates a new concurrent test helper
func NewConcurrentTestHelper(t *testing.T, numRoutines int) *ConcurrentTestHelper {
	t.Helper()
	return &ConcurrentTestHelper{
		t:           t,
		numRoutines: numRoutines,
		results:     make([]any, numRoutines),
		errors:      make([]error, numRoutines),
	}
}

// Run executes the test function concurrently
func (h *ConcurrentTestHelper) Run(testFunc func(int) (any, error)) {
	h.t.Helper()
	h.wg.Add(h.numRoutines)

	for i := 0; i < h.numRoutines; i++ {
		go func(idx int) {
			defer h.wg.Done()
			result, err := testFunc(idx)

			h.mutex.Lock()
			h.results[idx] = result
			h.errors[idx] = err
			h.mutex.Unlock()
		}(i)
	}

	h.wg.Wait()
}

// GetResults returns all results and errors
func (h *ConcurrentTestHelper) GetResults() ([]any, []error) {
	h.t.Helper()
	return h.results, h.errors
}

// AssertNoErrors checks that no goroutines returned errors
func (h *ConcurrentTestHelper) AssertNoErrors() {
	h.t.Helper()
	for i, err := range h.errors {
		if err != nil {
			h.t.Errorf("Goroutine %d failed: %v", i, err)
		}
	}
}

// GetResult returns the result at index with type assertion
func GetResult[T any](values []any, index int) (T, bool) {
	if index < 0 || index >= len(values) {
		var zero T
		return zero, false
	}
	result, ok := values[index].(T)
	return result, ok
}
