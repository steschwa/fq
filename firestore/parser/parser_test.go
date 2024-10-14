package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	assert := assert.New(t)

	fixtures := []struct {
		source string

		key        string
		operator   string
		valueType  Value
		valueValue any
	}{
		{source: `foo == "bar"`, key: "foo", operator: "==", valueType: StringValue{}, valueValue: "bar"},
		{source: `name != 'hi'`, key: "name", operator: "!=", valueType: StringValue{}, valueValue: "hi"},
		{source: `active == true`, key: "active", operator: "==", valueType: BoolValue{}, valueValue: true},
		{source: `deleted == false`, key: "deleted", operator: "==", valueType: BoolValue{}, valueValue: false},
		{source: `age > 30`, key: "age", operator: ">", valueType: IntValue{}, valueValue: 30},
		{source: `price <= 2.5`, key: "price", operator: "<=", valueType: FloatValue{}, valueValue: 2.5},
		{source: `date == null`, key: "date", operator: "==", valueType: NullValue{}, valueValue: nil},
	}

	for _, fixture := range fixtures {
		w, err := Parse(fixture.source)

		assert.Nil(err)

		assert.Equal(fixture.key, string(w.Key))
		assert.Equal(fixture.operator, w.Operator.String())
		assert.IsType(fixture.valueType, w.Value)
		assert.Equal(fixture.valueValue, w.Value.Value())
	}
}

func TestParseList(t *testing.T) {
	assert := assert.New(t)

	w, err := Parse(`rules in ["foo", 10, 5.75, true, false, null]`)

	assert.Nil(err)
	assert.Equal("rules", string(w.Key))
	assert.Equal("in", w.Operator.String())
	assert.IsType(ArrayValue{}, w.Value)

	arr := w.Value.(ArrayValue)

	assert.IsType(StringValue{}, arr.values[0])
	assert.IsType(IntValue{}, arr.values[1])
	assert.IsType(FloatValue{}, arr.values[2])
	assert.IsType(BoolValue{}, arr.values[3])
	assert.IsType(BoolValue{}, arr.values[4])
	assert.IsType(NullValue{}, arr.values[5])

	assert.Equal("foo", arr.values[0].Value())
	assert.Equal(10, arr.values[1].Value())
	assert.Equal(5.75, arr.values[2].Value())
	assert.Equal(true, arr.values[3].Value())
	assert.Equal(false, arr.values[4].Value())
	assert.Equal(nil, arr.values[5].Value())
}
