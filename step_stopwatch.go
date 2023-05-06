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
	label string
	time  time.Time
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
			label: t.steps[i].label,
			time:  t.steps[i].time,
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
		time: tm,
	})
}

// Step is like a lap on a stopwatch, it records the amount of time since the
//
//	the last step and marks this step with the given label
func (t *stopwatch) Step(label string) {
	if !t.running {
		panic("stopwatch not running")
	}

	t.steps = append(t.steps, &step{
		label: label,
		time:  t.now(),
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
	prevStep := t.steps[0]
	for i := 1; i < len(t.steps); i++ {
		currStep := t.steps[i]
		duration := currStep.time.Sub(prevStep.time)
		durationStr := durationMillisecondStr(duration)
		length, err := bw.WriteString(fmt.Sprintf(fmtstr, currStep.label, durationStr))
		if err != nil {
			return err
		}
		if length > longestLine {
			longestLine = length
		}

		prevStep = currStep
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
	if len(t.steps) < 2 {
		return &Results{}
	}

	results := &Results{
		Steps: make([]Step, len(t.steps)-1),
	}
	prevStep := t.steps[0]
	for i := 1; i < len(t.steps); i++ {
		currStep := t.steps[i]
		duration := currStep.time.Sub(prevStep.time)
		results.Steps[i-1].Label = currStep.label
		results.Steps[i-1].Duration = duration
		prevStep = currStep
	}
	return results
}

// GetResultMap returns the stopwatch step results in a native map format.
func (t *stopwatch) GetResultMap() []map[string]int64 {
	if len(t.steps) < 2 {
		return nil
	}

	results := make([]map[string]int64, len(t.steps)-1)

	prevStep := t.steps[0]
	for i := 1; i < len(t.steps); i++ {
		currStep := t.steps[i]
		duration := currStep.time.Sub(prevStep.time)
		results[i-1] = map[string]int64{
			currStep.label: int64(duration),
		}
		prevStep = currStep
	}
	return results
}

func durationMillisecondStr(d time.Duration) string {
	ms := float64(d) / float64(time.Millisecond)
	return fmt.Sprintf("%fms", ms)
}
