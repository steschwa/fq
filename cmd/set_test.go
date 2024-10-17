package cmd

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJSONObjectDecoding(t *testing.T) {
	assert := assert.New(t)

	fixtures := []struct {
		data        string
		shouldError bool
	}{
		{data: `{}`, shouldError: false},
		{data: `{"foo": "bar"}`, shouldError: false},
		{data: `[]`, shouldError: true},
		{data: `[{}]`, shouldError: true},
		{data: `[{"foo": "bar"}]`, shouldError: true},
		{data: `1`, shouldError: true},
		{data: `"foo"`, shouldError: true},
		{data: `true`, shouldError: true},
		{data: `false`, shouldError: true},
		{data: `null`, shouldError: true},
	}

	for i, fixture := range fixtures {
		var j jsonObject
		err := json.Unmarshal([]byte(fixture.data), &j)

		if fixture.shouldError {
			assert.Error(err, "index %d", i)
		} else {
			assert.NoError(err, "index %d", i)
		}
	}
}

func TestJSONArrayDecoding(t *testing.T) {
	assert := assert.New(t)

	fixtures := []struct {
		data        string
		shouldError bool
	}{
		{data: `[]`, shouldError: false},
		{data: `[{}]`, shouldError: false},
		{data: `[{"foo": "bar"}]`, shouldError: false},
		{data: `[1]`, shouldError: true},
		{data: `["foo"]`, shouldError: true},
		{data: `[true]`, shouldError: true},
		{data: `[false]`, shouldError: true},
		{data: `[null]`, shouldError: true},
		{data: `{}`, shouldError: true},
		{data: `{"foo": "bar"}`, shouldError: true},
		{data: `1`, shouldError: true},
		{data: `"foo"`, shouldError: true},
		{data: `true`, shouldError: true},
		{data: `false`, shouldError: true},
		{data: `null`, shouldError: true},
	}

	for i, fixture := range fixtures {
		var j jsonArray
		err := json.Unmarshal([]byte(fixture.data), &j)

		if fixture.shouldError {
			assert.Error(err, "index %d", i)
		} else {
			assert.NoError(err, "index %d", i)
		}
	}
}
