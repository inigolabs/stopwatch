package profiletimer

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

// Timer is the interface for profile or stopwatch timer.
type Timer interface {
	Step(label string)
	WriteResults(w io.Writer) error
	ShowResults() error
}

type profileTimer struct {
	steps []*step
}

type step struct {
	label string
	time  time.Time
}

// StartProfileTimer creates a profile timer and starts it.
func StartProfileTimer() *profileTimer {
	t := &profileTimer{}
	t.steps = append(t.steps, &step{"init", time.Now()})
	return t
}

// Step is like a lap on a stopwatch, it records the amount of time since the
//  the last step and marks this step with the given label
func (t *profileTimer) Step(label string) {
	t.steps = append(t.steps, &step{label, time.Now()})
}

// WriteResults writes the timer step results to the given writer.
func (t *profileTimer) WriteResults(w io.Writer) error {
	maxLabelLength := 0
	for _, s := range t.steps {
		if len(s.label) > maxLabelLength {
			maxLabelLength = len(s.label)
		}
	}
	fmtstr := fmt.Sprintf("%%-%ds : %%v\n", maxLabelLength)

	bw := bufio.NewWriter(w)
	longestLine := 0
	for i := 1; i < len(t.steps); i++ {
		step := t.steps[i]
		duration := step.time.Sub(t.steps[i-1].time)
		durationStr := durationMillisecondStr(duration)
		length, err := bw.WriteString(fmt.Sprintf(fmtstr, step.label, durationStr))
		if err != nil {
			return err
		}
		if length > longestLine {
			longestLine = length
		}
	}
	seperator := strings.Repeat("-", longestLine) + "\n"
	_, err := bw.WriteString(seperator)
	if err != nil {
		return err
	}

	lastStep := t.steps[len(t.steps)-1]
	totalDuration := lastStep.time.Sub(t.steps[0].time)
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
func (t *profileTimer) ShowResults() error {
	return t.WriteResults(os.Stdout)
}

type noopTimer struct{}

// StartNoopTimer creates a timer that does nothing. This is useful in cases
//  where you want to keep the timer step call in production code, but don't
//  want them to have any logical or performance effects.
func StartNoopTimer() *noopTimer {
	return &noopTimer{}
}

// Step for noop timer does nothing
func (t *noopTimer) Step(label string) {}

// WriteResults for noop timer does nothing
func (t *noopTimer) WriteResults(w io.Writer) error { return nil }

// ShowResults for noop timer does nothing
func (t *noopTimer) ShowResults() error { return nil }

func durationMillisecondStr(d time.Duration) string {
	ms := float64(d) / float64(time.Millisecond)
	return fmt.Sprintf("%fms", ms)
}
