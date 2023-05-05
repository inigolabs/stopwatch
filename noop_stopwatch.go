package stopwatch

import (
	"io"
	"time"
)

type noopStopWatch struct{}

// StartNoopStopWatch creates a stopwatch that does nothing.
//
// This is useful in cases	where you want to keep the stopwatch step call in production code,
// but don't want them to have any logical or performance effects.
func StartNoopStopWatch() *noopStopWatch {
	return &noopStopWatch{}
}

// Copy for noop stopwatch does nothing
func (t *noopStopWatch) Copy() StopWatch { return &noopStopWatch{} }

// Start for noop stopwatch does nothing
func (t *noopStopWatch) Start() {}

// StartWithTime for noop stopwatch does nothing
func (t *noopStopWatch) StartWithTime(tm time.Time) {}

// Stop for noop stopwatch does nothing
func (t *noopStopWatch) Stop() {}

// Step for noop stopwatch does nothing
func (t *noopStopWatch) Step(label string) {}

// WriteResults for noop stopwatch does nothing
func (t *noopStopWatch) WriteResults(w io.Writer) error { return nil }

// ShowResults for noop stopwatch does nothing
func (t *noopStopWatch) ShowResults() error { return nil }

// GetResults for noop stopwatch does nothing
func (t *noopStopWatch) GetResults() *Results { return nil }

// GetResultsMap for noop stopwatch does nothing
func (t *noopStopWatch) GetResultMap() []map[string]int64 { return nil }
