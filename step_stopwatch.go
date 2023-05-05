package stopwatch

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

type stopwatch struct {
	steps []*step

	now func() time.Time

	running bool
}

type step struct {
	label     string
	startTime time.Time
	endTime   time.Time
}

// Start creates a profile stopwatch
func NewStopWatch() StopWatch {
	return &stopwatch{
		now:     time.Now,
		running: false,
	}
}

// Copy copies the stopwatch
func (t *stopwatch) Copy() StopWatch {
	if t == nil {
		return nil
	}

	var steps = make([]*step, len(t.steps))

	for i := range steps {
		steps[i] = &step{
			label:     t.steps[i].label,
			startTime: t.steps[i].startTime,
			endTime:   t.steps[i].endTime,
		}
	}

	return &stopwatch{
		now:     t.now,
		running: t.running,
		steps:   steps,
	}
}

// Start starts the stopwatch
func (t *stopwatch) Start() {
	t.StartWithTime(t.now())
}

// StartWithTime starts the stopwatch at a given time
func (t *stopwatch) StartWithTime(tm time.Time) {
	if t.running {
		panic("stopwatch already running")
	}
	t.running = true
	t.steps = append(t.steps, &step{
		startTime: tm,
	})
}

// Stop stops the stopwatch
func (t *stopwatch) Stop() {
	if !t.running {
		panic("stopwatch already stopped")
	}
	t.running = false

	t.steps = t.steps[:len(t.steps)-1]
}

// Step is like a lap on a stopwatch, it records the amount of time since the
//
//	the last step and marks this step with the given label
func (t *stopwatch) Step(label string) {
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

// WriteResults writes the stopwatch step results to the given writer.
func (t *stopwatch) WriteResults(w io.Writer) error {
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

// ShowResults outputs the stopwatch step results to stdout.
func (t *stopwatch) ShowResults() error {
	return t.WriteResults(os.Stdout)
}

// GetResults returns the stopwatch step results.
func (t *stopwatch) GetResults() *Results {
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

// GetResultMap returns the stopwatch step results in a native map format.
func (t *stopwatch) GetResultMap() []map[string]int64 {
	results := []map[string]int64{}
	for i := 0; i < len(t.steps); i++ {
		step := t.steps[i]
		duration := step.endTime.Sub(step.startTime)
		results = append(results, map[string]int64{
			step.label: int64(duration),
		})
	}
	return results
}

func durationMillisecondStr(d time.Duration) string {
	ms := float64(d) / float64(time.Millisecond)
	return fmt.Sprintf("%fms", ms)
}
