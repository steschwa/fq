package parser

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLex(t *testing.T) {
	assert := assert.New(t)

	fixtures := []struct {
		value    string
		expected token
	}{
		{value: `    `, expected: token{kind: tokenWhitespace}},
		{value: `,`, expected: token{kind: tokenComma, value: ","}},
		{value: `[`, expected: token{kind: tokenSquareBracketOpen, value: "["}},
		{value: `]`, expected: token{kind: tokenSquareBracketClose, value: "]"}},
		{value: `50`, expected: token{kind: tokenNumber, value: "50"}},
		{value: `-50`, expected: token{kind: tokenNumber, value: "-50"}},
		{value: `10.6`, expected: token{kind: tokenNumber, value: "10.6"}},
		{value: `-10.6`, expected: token{kind: tokenNumber, value: "-10.6"}},
		{value: `false`, expected: token{kind: tokenFalse, value: "false"}},
		{value: `true`, expected: token{kind: tokenTrue, value: "true"}},
		{value: `"foo"`, expected: token{kind: tokenString, value: "foo"}},
		{value: `'foo'`, expected: token{kind: tokenString, value: "foo"}},
		{value: `null`, expected: token{kind: tokenNull, value: "null"}},
		{value: `foo`, expected: token{kind: tokenIdent, value: "foo"}},
		{value: `(`, expected: token{kind: tokenIllegal, value: ""}},
		{value: `.`, expected: token{kind: tokenIllegal, value: ""}},
		{value: `_`, expected: token{kind: tokenIllegal, value: ""}},
		{value: `-`, expected: token{kind: tokenIllegal, value: ""}},
	}

	for i, fixture := range fixtures {
		lexer := newValueLexer(fixture.value)
		token := lexer.lex()

		assert.Equal(fixture.expected.kind, token.kind, fmt.Sprintf("fixture: %d", i+1))
		assert.Equal(fixture.expected.value, token.value, fmt.Sprintf("fixture: %d", i+1))
	}
}

func TestLexNumbers(t *testing.T) {
	assert := assert.New(t)

	lexer := newValueLexer(`5`)
	token := lexer.lex()
	assert.Equal(token.kind, tokenNumber)
	assert.Equal(token.value, "5")

	lexer = newValueLexer(`5.0`)
	token = lexer.lex()
	assert.Equal(token.kind, tokenNumber)
	assert.Equal(token.value, "5.0")

	lexer = newValueLexer(`.5`)
	token = lexer.lex()
	assert.Equal(token.kind, tokenIllegal)

	lexer = newValueLexer(`5.`)
	token = lexer.lex()
	assert.Equal(token.kind, tokenIllegal)
}
