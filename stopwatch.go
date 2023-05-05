package stopwatch

import (
	"io"
	"time"
)

// StopWatch is the interface for the stopwatch
type StopWatch interface {
	Start()
	StartWithTime(tm time.Time)
	Step(label string)
	Stop()
	Copy() StopWatch

	WriteResults(w io.Writer) error
	ShowResults() error
	GetResults() *Results
	GetResultMap() []map[string]int64
}

type Results struct {
	Steps []Step `json:"steps"`
}

type Step struct {
	Label    string        `json:"label"`
	Duration time.Duration `json:"duration"`
}
