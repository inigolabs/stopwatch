package stopwatch

import (
	"io"
	"time"
)

type testStopwatch struct {
	sw *stopwatch
}

// Start creates a profile stopwatch
func NewTestStopWatch() StopWatch {
	tsw := &testStopwatch{}
	tsw.sw = &stopwatch{
		running: false,
		now:     tsw.mockNow,
	}

	return tsw
}

// Copy copies the stopwatch
func (t *testStopwatch) Copy() StopWatch {
	return nil
}

// Start starts the stopwatch
func (t *testStopwatch) Start() {
	t.StartWithTime(time.Now())
}

// StartWithTime starts the stopwatch at a given time
func (t *testStopwatch) StartWithTime(tm time.Time) {
	if t.sw.running {
		panic("stopwatch already running")
	}
	t.sw.running = true
	t.sw.steps = append(t.sw.steps, &step{
		time: t.sw.now(),
	})
}

// Step is like a lap on a stopwatch, it records the amount of time since the
//
//	the last step and marks this step with the given label
func (t *testStopwatch) Step(label string) {
	if !t.sw.running {
		panic("stopwatch not running")
	}
	t.sw.steps = append(t.sw.steps, &step{
		label: label,
		time:  t.sw.now(),
	})
}

// WriteResults writes the stopwatch step results to the given writer.
func (t *testStopwatch) WriteResults(w io.Writer) error {
	return t.sw.WriteResults(w)
}

// ShowResults outputs the stopwatch step results to stdout.
func (t *testStopwatch) ShowResults() error {
	return t.sw.ShowResults()
}

// GetResults returns the stopwatch step results.
func (t *testStopwatch) GetResults() *Results {
	return t.sw.GetResults()
}

// GetResultMap returns the stopwatch step results in a native map format.
func (t *testStopwatch) GetResultMap() []map[string]int64 {
	return t.sw.GetResultMap()
}

func (t *testStopwatch) mockNow() time.Time {
	base := time.Date(1000, 0, 0, 0, 0, 0, 0, time.UTC)
	shift := time.Duration(len(t.sw.steps)) * time.Millisecond
	return base.Add(shift)
}
