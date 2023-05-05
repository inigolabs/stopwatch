package stopwatch

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"
)

func TestCodec(t *testing.T) {
	results := &Results{
		Steps: []Step{
			{
				Label:    "B",
				Duration: 1 * time.Millisecond,
			},
			{
				Label:    "D",
				Duration: 2 * time.Millisecond,
			},
			{
				Label:    "A",
				Duration: 3 * time.Millisecond,
			},
			{
				Label:    "F",
				Duration: 4 * time.Millisecond,
			},
			{
				Label:    "E",
				Duration: 5 * time.Millisecond,
			},
		},
	}

	labels := StopWatchLabels{"A", "B", "C", "D", "E", "F"}

	encoded, err := Encode(labels, results)

	require.NoError(t, err)
	decoded, err := Decode(labels, encoded)
	require.NoError(t, err)
	diff := cmp.Diff(results, decoded)
	if diff != "" {
		t.Fatal(diff)
	}
}
