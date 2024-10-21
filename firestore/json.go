package firestore

import (
	"encoding/json"
	"fmt"
)

type (
	JSONObject struct {
		Value map[string]any
	}
	JSONArray struct {
		Values []map[string]any
	}
)

func (j *JSONObject) UnmarshalJSON(bytes []byte) error {
	var data any
	if err := json.Unmarshal(bytes, &data); err != nil {
		return err
	}

	switch data := data.(type) {
	case map[string]any:
		j.Value = data
		return nil
	default:
		return fmt.Errorf("expected json object, got %T", data)
	}
}

func (j *JSONArray) UnmarshalJSON(bytes []byte) error {
	var data any
	if err := json.Unmarshal(bytes, &data); err != nil {
		return err
	}

	switch data := data.(type) {
	case []any:
		objects := make([]map[string]any, len(data))
		for i, value := range data {
			switch value := value.(type) {
			case map[string]any:
				objects[i] = value
			default:
				return fmt.Errorf("no json object in array at pos %d", i+1)
			}
		}

		j.Values = objects

		return nil
	default:
		return fmt.Errorf("expected json array, got %T", data)
	}
}
