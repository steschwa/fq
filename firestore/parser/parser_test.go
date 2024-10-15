package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	assert := assert.New(t)

	fixtures := []struct {
		source string

		key      string
		operator string
		value    Value
	}{
		{source: `foo == "bar"`, key: "foo", operator: "==", value: StringValue{value: "bar"}},
		{source: `name != 'hi'`, key: "name", operator: "!=", value: StringValue{value: "hi"}},
		{source: `active == true`, key: "active", operator: "==", value: BoolValue{value: true}},
		{source: `deleted == false`, key: "deleted", operator: "==", value: BoolValue{value: false}},
		{source: `age > 30`, key: "age", operator: ">", value: IntValue{value: 30}},
		{source: `price <= 2.5`, key: "price", operator: "<=", value: FloatValue{value: 2.5}},
		{source: `date == null`, key: "date", operator: "==", value: NullValue{}},
	}

	for _, fixture := range fixtures {
		w, err := Parse(fixture.source)

		assert.NoError(err)

		assert.Equal(fixture.key, string(w.Key))
		assert.Equal(fixture.operator, w.Operator.String())
		assert.IsType(fixture.value, w.Value)
		assert.Equal(fixture.value.Value(), w.Value.Value())
	}
}

func TestParseArray(t *testing.T) {
	assert := assert.New(t)

	v, err := parseValue(`["foo", 10, 5.75, true, false, null]`)

	assert.NoError(err)
	assert.IsType(ArrayValue{}, v)

	arr := v.(ArrayValue)

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

func TestParseOperator(t *testing.T) {
	assert := assert.New(t)

	o, err := parseOperator("==")
	assert.NoError(err)
	assert.Equal(Eq, o)

	o, err = parseOperator("!=")
	assert.NoError(err)
	assert.Equal(Neq, o)

	o, err = parseOperator(">")
	assert.NoError(err)
	assert.Equal(Gt, o)

	o, err = parseOperator("<")
	assert.NoError(err)
	assert.Equal(Lt, o)

	o, err = parseOperator(">=")
	assert.NoError(err)
	assert.Equal(Gte, o)

	o, err = parseOperator("<=")
	assert.NoError(err)
	assert.Equal(Lte, o)

	o, err = parseOperator("in")
	assert.NoError(err)
	assert.Equal(In, o)

	o, err = parseOperator("array-contains-any")
	assert.NoError(err)
	assert.Equal(o, ArrayContainsAny)

	o, err = parseOperator("")
	assert.Error(err)
	assert.ErrorIs(err, errInvalidOperator)
}
