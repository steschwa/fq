package firestore

import (
	"fmt"
	"strings"
)

type (
	KeyPath string

	Operator int

	Value interface {
		String() string
		Value() any
	}

	StringValue struct {
		value string
	}
	IntValue struct {
		value int
	}
	FloatValue struct {
		value float64
	}
	BoolValue struct {
		value bool
	}
	NullValue  struct{}
	ArrayValue struct {
		Values []Value
	}

	Where struct {
		Key      KeyPath
		Operator Operator
		Value    Value
	}
)

const (
	Eq               Operator = iota + 1 // ==
	Neq                                  // !=
	Gt                                   // >
	Lt                                   // <
	Gte                                  // >=
	Lte                                  // <=
	In                                   // in
	ArrayContainsAny                     // array-contains-any
)

var (
	_ Value = StringValue{}
	_ Value = IntValue{}
	_ Value = FloatValue{}
	_ Value = BoolValue{}
	_ Value = ArrayValue{}
	_ Value = NullValue{}
)

func (p KeyPath) Segments() []string {
	return strings.Split(string(p), ".")
}

func (o Operator) String() string {
	switch o {
	case Eq:
		return "=="
	case Neq:
		return "!="
	case Gt:
		return ">"
	case Lt:
		return "<"
	case Gte:
		return ">="
	case Lte:
		return "<="
	case In:
		return "in"
	case ArrayContainsAny:
		return "array-contains-any"
	default:
		return ""
	}
}

func NewStringValue(value string) StringValue {
	return StringValue{value: value}
}
func (v StringValue) Value() any {
	return v.value
}
func (v StringValue) String() string {
	return fmt.Sprintf(`"%s"`, v.value)
}

func NewIntValue(value int) IntValue {
	return IntValue{value: value}
}
func (v IntValue) Value() any {
	return v.value
}
func (v IntValue) String() string {
	return fmt.Sprint(v.value)
}

func NewFloatValue(value float64) FloatValue {
	return FloatValue{value: value}
}
func (v FloatValue) Value() any {
	return v.value
}
func (v FloatValue) String() string {
	return fmt.Sprint(v.value)
}

func NewBoolValue(value bool) BoolValue {
	return BoolValue{value: value}
}
func (v BoolValue) Value() any {
	return v.value
}
func (v BoolValue) String() string {
	return fmt.Sprint(v.value)
}

func NewNullValue() NullValue {
	return NullValue{}
}
func (v NullValue) Value() any {
	return nil
}
func (v NullValue) String() string {
	return "null"
}

func NewArrayValue() ArrayValue {
	return ArrayValue{}
}
func (v *ArrayValue) Add(value Value) {
	v.Values = append(v.Values, value)
}
func (v ArrayValue) Value() any {
	list := make([]any, len(v.Values))
	for i, value := range v.Values {
		list[i] = value.Value()
	}

	return list
}
func (v ArrayValue) String() string {
	members := make([]string, len(v.Values))
	for i, v := range v.Values {
		members[i] = v.String()
	}

	formattedMembers := strings.Join(members, ", ")

	return fmt.Sprintf("[%s]", formattedMembers)
}

func (w Where) String() string {
	return fmt.Sprintf("%s %s %s", w.Key, w.Operator.String(), w.Value.String())
}
