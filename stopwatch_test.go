package stopwatch

import (
	"testing"
	"time"

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
	testTimer := &timer{
		now: mock.now,
	}

	mock.time = time.Date(2022, 2, 2, 22, 22, 22, 0, time.UTC)
	testTimer.Start()
	mock.time = time.Date(2022, 2, 2, 22, 22, 23, 0, time.UTC)
	testTimer.Step("a")
	mock.time = time.Date(2022, 2, 2, 22, 22, 25, 0, time.UTC)
	testTimer.Step("b")
	mock.time = time.Date(2022, 2, 2, 22, 22, 28, 0, time.UTC)
	testTimer.Step("c")
	testTimer.Stop()
	actual := testTimer.GetResults()

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
