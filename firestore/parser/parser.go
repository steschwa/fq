package parser

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/steschwa/fst/firestore"
)

var (
	errInvalidOperator = errors.New("invalid operator")
	errNoTokens        = errors.New("no tokens")
)

var (
	whereRe = regexp.MustCompile(`^([a-zA-Z._]+) (==|!=|>|<|>=|<=|in|array-contains-any) (.*)$`)
)

func Parse(source string) (firestore.Where, error) {
	matches := whereRe.FindStringSubmatch(source)
	if len(matches) != 4 {
		return firestore.Where{}, fmt.Errorf("regex matching where")
	}

	key := matches[1]

	op, err := parseOperator(matches[2])
	if err != nil {
		return firestore.Where{}, fmt.Errorf("parsing operator: %v", err)
	}

	rawValue := matches[3]
	value, err := parseValue(rawValue)
	if err != nil {
		return firestore.Where{}, fmt.Errorf("parsing value: %v", err)
	}

	return firestore.Where{
		Key:      firestore.KeyPath(key),
		Operator: op,
		Value:    value,
	}, nil
}

func parseOperator(op string) (firestore.Operator, error) {
	switch op {
	case "==":
		return firestore.Eq, nil
	case "!=":
		return firestore.Neq, nil
	case ">":
		return firestore.Gt, nil
	case "<":
		return firestore.Lt, nil
	case ">=":
		return firestore.Gte, nil
	case "<=":
		return firestore.Lte, nil
	case "in":
		return firestore.In, nil
	case "array-contains-any":
		return firestore.ArrayContainsAny, nil
	default:
		return firestore.Operator(0), errInvalidOperator
	}
}

func parseValue(value string) (firestore.Value, error) {
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

func parseValueToken(token token) (firestore.Value, error) {
	switch token.kind {
	case tokenString:
		return firestore.NewStringValue(token.value), nil
	case tokenTrue:
		return firestore.NewBoolValue(true), nil
	case tokenFalse:
		return firestore.NewBoolValue(false), nil
	case tokenNull:
		return firestore.NullValue{}, nil
	case tokenNumber:
		if strings.Contains(token.value, ".") {
			v, err := strconv.ParseFloat(token.value, 64)
			if err != nil {
				return nil, fmt.Errorf("parsing float: %v", err)
			}

			return firestore.NewFloatValue(v), nil
		}

		v, err := strconv.Atoi(token.value)
		if err != nil {
			return nil, fmt.Errorf("parsing int: %v", err)
		}

		return firestore.NewIntValue(v), nil
	}

	return nil, fmt.Errorf("invalid token: %s (%s)", token.kind.String(), token.value)
}

func parseListValueTokens(tokens []token) (firestore.ArrayValue, error) {
	arrayValue := firestore.ArrayValue{}
	for _, token := range tokens {
		v, err := parseValueToken(token)
		if err != nil {
			return firestore.ArrayValue{}, fmt.Errorf("parsing list token: %v", err)
		}

		arrayValue.Add(v)
	}

	return arrayValue, nil
}
