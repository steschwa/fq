package parser

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type (
	KeyPath string

	Operator int

	Value interface {
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
	ArrayValue struct {
		values []Value
	}
	NullValue struct{}

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
	errInvalidOperator = errors.New("invalid operator")

	errNoTokens = errors.New("no tokens")
)

var (
	_ Value = StringValue{}
	_ Value = IntValue{}
	_ Value = FloatValue{}
	_ Value = BoolValue{}
	_ Value = ArrayValue{}
	_ Value = NullValue{}
)

var (
	whereRe = regexp.MustCompile(`^([a-zA-Z._]+) (==|!=|>|<|>=|<=|in|array-contains-any) (.*)$`)
)

func Parse(source string) (Where, error) {
	matches := whereRe.FindStringSubmatch(source)
	if len(matches) != 4 {
		return Where{}, fmt.Errorf("regex matching where")
	}

	key := matches[1]

	op, err := parseOperator(matches[2])
	if err != nil {
		return Where{}, fmt.Errorf("parsing operator: %v", err)
	}

	rawValue := matches[3]
	value, err := parseValue(rawValue)
	if err != nil {
		return Where{}, fmt.Errorf("parsing value: %v", err)
	}

	return Where{
		Key:      KeyPath(key),
		Operator: op,
		Value:    value,
	}, nil
}

func parseOperator(op string) (Operator, error) {
	switch op {
	case "==":
		return Eq, nil
	case "!=":
		return Neq, nil
	case ">":
		return Gt, nil
	case "<":
		return Lt, nil
	case ">=":
		return Gte, nil
	case "<=":
		return Lte, nil
	case "in":
		return In, nil
	case "array-contains-any":
		return ArrayContainsAny, nil
	default:
		return Operator(0), errInvalidOperator
	}
}

func parseValue(value string) (Value, error) {
	lexer := newValueLexer(value)
	var tokens []token
	for {
		token := lexer.lex()
		if token.kind == tokenEOF {
			break
		} else if token.kind == tokenIllegal {
			return nil, fmt.Errorf("illegal token")
		} else if token.kind == tokenWhitespace {
			continue
		} else if token.kind == tokenComma {
			continue
		}

		tokens = append(tokens, token)
	}

	if len(tokens) == 0 {
		return nil, errNoTokens
	}

	if len(tokens) == 1 {
		v, err := parseValueToken(tokens[0])
		if err != nil {
			return nil, fmt.Errorf("parsing token: %v", err)
		}

		return v, nil
	}

	var (
		firstToken = tokens[0]
		lastToken  = tokens[len(tokens)-1]
	)
	if firstToken.kind == tokenSquareBracketOpen && lastToken.kind == tokenSquareBracketClose {
		actualTokens := tokens[1 : len(tokens)-1]
		v, err := parseListValueTokens(actualTokens)
		if err != nil {
			return nil, fmt.Errorf("parsing token list: %v", err)
		}

		return v, nil
	}

	return nil, fmt.Errorf("invalid value")
}

func parseValueToken(token token) (Value, error) {
	switch token.kind {
	case tokenString:
		return StringValue{value: token.value}, nil
	case tokenTrue:
		return BoolValue{value: true}, nil
	case tokenFalse:
		return BoolValue{value: false}, nil
	case tokenNull:
		return NullValue{}, nil
	case tokenNumber:
		if strings.Contains(token.value, ".") {
			v, err := strconv.ParseFloat(token.value, 64)
			if err != nil {
				return nil, fmt.Errorf("parsing float: %v", err)
			}

			return FloatValue{value: v}, nil
		}

		v, err := strconv.Atoi(token.value)
		if err != nil {
			return nil, fmt.Errorf("parsing int: %v", err)
		}

		return IntValue{value: v}, nil
	}

	return nil, fmt.Errorf("invalid token: %s (%s)", token.kind.String(), token.value)
}

func parseListValueTokens(tokens []token) (ArrayValue, error) {
	arrayValue := ArrayValue{}
	for _, token := range tokens {
		v, err := parseValueToken(token)
		if err != nil {
			return ArrayValue{}, fmt.Errorf("parsing list token: %v", err)
		}

		arrayValue.values = append(arrayValue.values, v)
	}

	return arrayValue, nil
}

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

func (v StringValue) Value() any {
	return v.value
}

func (v IntValue) Value() any {
	return v.value
}

func (v FloatValue) Value() any {
	return v.value
}

func (v BoolValue) Value() any {
	return v.value
}

func (v ArrayValue) Value() any {
	list := make([]any, len(v.values))
	for i, value := range v.values {
		list[i] = value.Value()
	}

	return list
}

func (v NullValue) Value() any {
	return nil
}
