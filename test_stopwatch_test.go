package stopwatch

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestTestStopwatchGetResults(t *testing.T) {
	testStopWatch := NewTestStopWatch()

	testStopWatch.Start()
	testStopWatch.Step("a")
	testStopWatch.Step("b")
	testStopWatch.Step("c")
	actual := testStopWatch.GetResults()

	expect := &Results{
		Steps: []Step{
			{Label: "a", Duration: time.Millisecond},
			{Label: "b", Duration: time.Millisecond},
			{Label: "c", Duration: time.Millisecond},
		},
	}

	diff := cmp.Diff(expect, actual)
	if diff != "" {
		t.Fatal(diff)
	}
}

func TestTestStopwatchGetResultMap(t *testing.T) {
	testStopWatch := NewTestStopWatch()

	testStopWatch.Start()
	testStopWatch.Step("a")
	testStopWatch.Step("b")
	testStopWatch.Step("c")
	actual := testStopWatch.GetResultMap()

	expect := []map[string]int64{
		{
			"a": int64(time.Millisecond),
		},
		{
			"b": int64(time.Millisecond),
		},
		{
			"c": int64(time.Millisecond),
		},
	}

	diff := cmp.Diff(expect, actual)
	if diff != "" {
		t.Fatal(diff)
	}
}
