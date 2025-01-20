package firestore

import (
	"encoding/json"
	"math"
)

type FirestoreDoc struct {
	Value map[string]any
}

var _ json.Marshaler = &FirestoreDoc{}

func NewFirestoreDoc(value map[string]any) *FirestoreDoc {
	prepareMap(value)
	return &FirestoreDoc{value}
}

func (d *FirestoreDoc) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.Value)
}

func prepareMap(v map[string]any) {
	for key, value := range v {
		switch value := value.(type) {
		case float64:
			if math.IsNaN(value) {
				v[key] = "NaN"
			}
		case []any:
			prepareSlice(value)
		case map[string]any:
			prepareMap(value)
		}
	}
}

func prepareSlice(v []any) {
	for _, value := range v {
		switch value := value.(type) {
		case map[string]any:
			prepareMap(value)
		}
	}
}
