package parser

import (
	"testing"

	"github.com/steschwa/fst/firestore"
	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	assert := assert.New(t)

	fixtures := []struct {
		source string

		key      string
		operator string
		value    firestore.Value
	}{
		{source: `foo == "bar"`, key: "foo", operator: "==", value: firestore.NewStringValue("bar")},
		{source: `name != 'hi'`, key: "name", operator: "!=", value: firestore.NewStringValue("hi")},
		{source: `active == true`, key: "active", operator: "==", value: firestore.NewBoolValue(true)},
		{source: `deleted == false`, key: "deleted", operator: "==", value: firestore.NewBoolValue(false)},
		{source: `age > 30`, key: "age", operator: ">", value: firestore.NewIntValue(30)},
		{source: `price <= 2.5`, key: "price", operator: "<=", value: firestore.NewFloatValue(2.5)},
		{source: `date == null`, key: "date", operator: "==", value: firestore.NewNullValue()},
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
	assert.IsType(firestore.ArrayValue{}, v)

	arr := v.(firestore.ArrayValue)

	assert.IsType(firestore.StringValue{}, arr.Values[0])
	assert.IsType(firestore.IntValue{}, arr.Values[1])
	assert.IsType(firestore.FloatValue{}, arr.Values[2])
	assert.IsType(firestore.BoolValue{}, arr.Values[3])
	assert.IsType(firestore.BoolValue{}, arr.Values[4])
	assert.IsType(firestore.NullValue{}, arr.Values[5])

	assert.Equal("foo", arr.Values[0].Value())
	assert.Equal(10, arr.Values[1].Value())
	assert.Equal(5.75, arr.Values[2].Value())
	assert.Equal(true, arr.Values[3].Value())
	assert.Equal(false, arr.Values[4].Value())
	assert.Equal(nil, arr.Values[5].Value())
}

func TestParseOperator(t *testing.T) {
	assert := assert.New(t)

	o, err := parseOperator("==")
	assert.NoError(err)
	assert.Equal(firestore.Eq, o)

	o, err = parseOperator("!=")
	assert.NoError(err)
	assert.Equal(firestore.Neq, o)

	o, err = parseOperator(">")
	assert.NoError(err)
	assert.Equal(firestore.Gt, o)

	o, err = parseOperator("<")
	assert.NoError(err)
	assert.Equal(firestore.Lt, o)

	o, err = parseOperator(">=")
	assert.NoError(err)
	assert.Equal(firestore.Gte, o)

	o, err = parseOperator("<=")
	assert.NoError(err)
	assert.Equal(firestore.Lte, o)

	o, err = parseOperator("in")
	assert.NoError(err)
	assert.Equal(firestore.In, o)

	o, err = parseOperator("array-contains-any")
	assert.NoError(err)
	assert.Equal(o, firestore.ArrayContainsAny)

	o, err = parseOperator("")
	assert.Error(err)
	assert.ErrorIs(err, errInvalidOperator)
}
