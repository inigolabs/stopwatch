package stopwatch

import (
	"encoding/json"
	"io"
	"time"
)

func (r *Results) UnmarshalGQL(v interface{}) error {
	steps, ok := v.(map[string]interface{})["steps"]
	if ok && steps != nil {
		steps := steps.([]interface{})
		for _, step := range steps {
			step := step.(map[string]interface{})
			r.Steps = append(r.Steps, Step{
				Label:    step["label"].(string),
				Duration: step["duration"].(time.Duration),
			})
		}
	}
	return nil
}

func (r *Results) MarshalGQL(w io.Writer) {
	encoder := json.NewEncoder(w)
	err := encoder.Encode(r)
	if err != nil {
		panic(err)
	}
}
