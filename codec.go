package stopwatch

import (
	"errors"
	"fmt"
	"time"
)

type StopWatchLabels []string

const maxSupportedDuration = ((1 << 26) - 1)

func Encode(labels StopWatchLabels, results *Results) ([]uint32, error) {
	labelMap := make(map[string]int)
	if len(labels) > 64 {
		return nil, errors.New("only up to 64 labels supported")
	}
	for i := 0; i < len(labels); i++ {
		labelMap[labels[i]] = i
	}

	out := make([]uint32, len(results.Steps))
	for i := 0; i < len(results.Steps); i++ {
		step := results.Steps[i]
		labelIndex, found := labelMap[step.Label]
		if !found {
			return nil, fmt.Errorf("label %s not found in encode label list", step.Label)
		}
		dur := step.Duration / time.Microsecond
		if dur > maxSupportedDuration {
			return nil, fmt.Errorf("label %s duration of %dus exceeds the max support duration %dus", step.Label, dur, maxSupportedDuration)
		}

		val := uint32(uint32(labelIndex<<26) + uint32(dur))
		out[i] = val
	}

	return out, nil
}

func Decode(labels StopWatchLabels, data []uint32) (*Results, error) {
	labelMap := make(map[int]string)
	if len(labels) > 64 {
		return nil, errors.New("only up to 64 labels supported")
	}
	for i := 0; i < len(labels); i++ {
		labelMap[i] = labels[i]
	}

	out := &Results{
		Steps: make([]Step, len(data)),
	}

	for i := 0; i < len(data); i++ {
		dur := int64(data[i]&maxSupportedDuration) * int64(time.Microsecond)
		index := int(data[i] >> 26)

		label, found := labelMap[index]
		if !found {
			return nil, fmt.Errorf("label for index %d not found", index)
		}

		out.Steps[i] = Step{
			Label:    label,
			Duration: time.Duration(dur),
		}
	}

	return out, nil
}
