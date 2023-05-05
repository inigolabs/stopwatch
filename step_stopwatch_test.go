package stopwatch

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/google/go-cmp/cmp"
)

type timeMock struct {
	time time.Time
}

func (t *timeMock) now() time.Time {
	return t.time
}

func TestStopwatchGetResults(t *testing.T) {
	mock := &timeMock{}
	testStopWatch := &stopwatch{
		now: mock.now,
	}

	mock.time = time.Date(2022, 2, 2, 22, 22, 22, 0, time.UTC)
	testStopWatch.Start()
	mock.time = time.Date(2022, 2, 2, 22, 22, 23, 0, time.UTC)
	testStopWatch.Step("a")
	mock.time = time.Date(2022, 2, 2, 22, 22, 25, 0, time.UTC)
	testStopWatch.Step("b")
	mock.time = time.Date(2022, 2, 2, 22, 22, 28, 0, time.UTC)
	testStopWatch.Step("c")
	testStopWatch.Stop()
	actual := testStopWatch.GetResults()

	expect := &Results{
		Steps: []Step{
			{Label: "a", Duration: 1 * time.Second},
			{Label: "b", Duration: 2 * time.Second},
			{Label: "c", Duration: 3 * time.Second},
		},
	}

	diff := cmp.Diff(expect, actual)
	if diff != "" {
		t.Fatal(diff)
	}
}

func TestStopwatchGetResultMap(t *testing.T) {
	mock := &timeMock{}
	testStopWatch := &stopwatch{
		now: mock.now,
	}

	mock.time = time.Date(2022, 2, 2, 22, 22, 22, 0, time.UTC)
	testStopWatch.Start()
	mock.time = time.Date(2022, 2, 2, 22, 22, 23, 0, time.UTC)
	testStopWatch.Step("a")
	mock.time = time.Date(2022, 2, 2, 22, 22, 25, 0, time.UTC)
	testStopWatch.Step("b")
	mock.time = time.Date(2022, 2, 2, 22, 22, 28, 0, time.UTC)
	testStopWatch.Step("c")
	testStopWatch.Stop()
	actual := testStopWatch.GetResultMap()

	expect := []map[string]int64{
		{
			"a": int64(1 * time.Second),
		},
		{
			"b": int64(2 * time.Second),
		},
		{
			"c": int64(3 * time.Second),
		},
	}

	diff := cmp.Diff(expect, actual)
	if diff != "" {
		t.Fatal(diff)
	}
}

func TestStopwatchCopy(t *testing.T) {
	sw := NewStopWatch()

	sw.Start()
	sw.Step("one")
	sw.Step("two")
	sw.Stop()

	swc := sw.Copy()

	swc.Start()
	swc.Step("three")
	swc.Step("four")
	swc.Stop()

	sw.Start()
	sw.Step("three")
	sw.Step("four")
	sw.Stop()

	t1, t2 := sw.GetResults(), swc.GetResults()

	require.Equal(t, len(t1.Steps), len(t2.Steps))
	require.Equal(t, t1.Steps[0], t2.Steps[0])
	require.Equal(t, t1.Steps[1], t2.Steps[1])
	require.NotEqual(t, t1.Steps[2], t2.Steps[2])
	require.NotEqual(t, t1.Steps[3], t2.Steps[3])
}
