package stopwatch

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

// Timer is the interface for the stopwatch timer.
type Timer interface {
	Start()
	Step(label string)
	Stop()

	WriteResults(w io.Writer) error
	ShowResults() error
	GetResults() *Results
}

type Results struct {
	Steps []Step `json:"steps"`
}

type Step struct {
	Label    string        `json:"label"`
	Duration time.Duration `json:"duration"`
}

type timer struct {
	steps []*step

	now func() time.Time

	running bool
}

type step struct {
	label     string
	startTime time.Time
	endTime   time.Time
}

// Start creates a profile timer and starts it.
func NewTimer() *timer {
	t := &timer{
		now:     time.Now,
	}
	t.Start()
	return t
}

// Start starts the stopwatch
func (t *timer) Start() {
	if t.running {
		panic("stopwatch already running")
	}
	t.running = true
	t.steps = append(t.steps, &step{
		startTime: t.now(),
	})
}

// Stop stops the stopwatch
func (t *timer) Stop() {
	if !t.running {
		panic("stopwatch already stopped")
	}
	t.running = false

	t.steps = t.steps[:len(t.steps)-1]
}

// Step is like a lap on a stopwatch, it records the amount of time since the
//  the last step and marks this step with the given label
func (t *timer) Step(label string) {
	if !t.running {
		panic("stopwatch not running")
	}
	now := t.now()

	lastStep := t.steps[len(t.steps)-1]
	lastStep.label = label
	lastStep.endTime = now

	t.steps = append(t.steps, &step{
		startTime: now,
	})
}

// WriteResults writes the timer step results to the given writer.
func (t *timer) WriteResults(w io.Writer) error {
	maxLabelLength := 0
	for _, s := range t.steps {
		if len(s.label) > maxLabelLength {
			maxLabelLength = len(s.label)
		}
	}
	fmtstr := fmt.Sprintf("%%-%ds : %%v\n", maxLabelLength)

	bw := bufio.NewWriter(w)
	longestLine := 0
	totalDuration := time.Duration(0)
	for i := 0; i < len(t.steps); i++ {
		step := t.steps[i]
		duration := step.endTime.Sub(step.startTime)
		durationStr := durationMillisecondStr(duration)
		length, err := bw.WriteString(fmt.Sprintf(fmtstr, step.label, durationStr))
		if err != nil {
			return err
		}
		if length > longestLine {
			longestLine = length
		}

		totalDuration += duration
	}

	seperator := strings.Repeat("-", longestLine) + "\n"
	_, err := bw.WriteString(seperator)
	if err != nil {
		return err
	}

	totalDurationStr := durationMillisecondStr(totalDuration)
	_, err = bw.WriteString(fmt.Sprintf(fmtstr, "total", totalDurationStr))
	if err != nil {
		return err
	}
	err = bw.Flush()
	if err != nil {
		return err
	}
	return nil
}

// ShowResults outputs the timer step results to stdout.
func (t *timer) ShowResults() error {
	return t.WriteResults(os.Stdout)
}

// GetResults returns the timer step results.
func (t *timer) GetResults() *Results {
	results := &Results{}
	for i := 0; i < len(t.steps); i++ {
		step := t.steps[i]
		duration := step.endTime.Sub(step.startTime)
		results.Steps = append(results.Steps, Step{
			Label:    step.label,
			Duration: duration,
		})
	}
	return results
}

type noopTimer struct{}

// StartNoopTimer creates a timer that does nothing. This is useful in cases
//  where you want to keep the timer step call in production code, but don't
//  want them to have any logical or performance effects.
func StartNoopTimer() *noopTimer {
	return &noopTimer{}
}

// Start for noop timer does nothing
func (t *noopTimer) Start() {}

// Stop for noop timer does nothing
func (t *noopTimer) Stop() {}

// Step for noop timer does nothing
func (t *noopTimer) Step(label string) {}

// WriteResults for noop timer does nothing
func (t *noopTimer) WriteResults(w io.Writer) error { return nil }

// ShowResults for noop timer does nothing
func (t *noopTimer) ShowResults() error { return nil }

// GetResults for noop timer does nothing
func (t *noopTimer) GetResults() *Results { return nil }

func durationMillisecondStr(d time.Duration) string {
	ms := float64(d) / float64(time.Millisecond)
	return fmt.Sprintf("%fms", ms)
}
