package parser

import (
	"strings"
	"unicode"
)

type (
	valueLexer struct {
		reader *strings.Reader
	}

	token struct {
		kind  tokenKind
		value string
	}

	tokenKind int
)

const (
	tokenEOF tokenKind = iota + 1
	tokenIllegal
	tokenWhitespace
	tokenComma

	tokenSquareBracketOpen
	tokenSquareBracketClose

	tokenIdent

	tokenString
	tokenNumber
	tokenTrue
	tokenFalse
	tokenNull
)

func newValueLexer(value string) *valueLexer {
	return &valueLexer{
		reader: strings.NewReader(value),
	}
}

func (l *valueLexer) read() rune {
	r, _, err := l.reader.ReadRune()
	if err != nil {
		return rune(0)
	}

	return r
}

func (l *valueLexer) unread() {
	l.reader.UnreadRune()
}

func (l *valueLexer) lex() token {
	r := l.read()

	switch r {
	case rune(0):
		return token{kind: tokenEOF}
	case ',':
		return token{kind: tokenComma, value: ","}
	case '[':
		return token{kind: tokenSquareBracketOpen, value: "["}
	case ']':
		return token{kind: tokenSquareBracketClose, value: "]"}
	case '"', '\'':
		l.unread()
		value := l.lexString()
		return token{kind: tokenString, value: value}
	case '-':
		l.unread()
		return l.lexNumber()
	}

	if unicode.IsSpace(r) {
		l.unread()
		l.lexWhitespace()
		return token{kind: tokenWhitespace}
	} else if unicode.IsLetter(r) {
		l.unread()
		value := l.lexIdent()

		switch value {
		case "true":
			return token{kind: tokenTrue, value: "true"}
		case "false":
			return token{kind: tokenFalse, value: "false"}
		case "null":
			return token{kind: tokenNull, value: "null"}
		}

		return token{kind: tokenIdent, value: value}
	} else if unicode.IsDigit(r) {
		l.unread()
		return l.lexNumber()
	}

	return token{kind: tokenIllegal}
}

func (l *valueLexer) lexWhitespace() {
	for {
		r := l.read()

		if unicode.IsSpace(r) {
			continue
		}

		l.unread()
		break
	}
}

func (l *valueLexer) lexIdent() string {
	value := ""

	for {
		r := l.read()

		if unicode.IsLetter(r) {
			value += string(r)
			continue
		}

		l.unread()
		break
	}

	return value
}

func (l *valueLexer) lexString() string {
	borderChar := l.read()
	value := ""

	for {
		r := l.read()

		if r == rune(0) || r == borderChar {
			break
		}

		value += string(r)
	}

	return value
}

func (l *valueLexer) lexNumber() token {
	value := ""

	for {
		r := l.read()

		if unicode.IsDigit(r) || r == '-' {
			value += string(r)
			continue
		}

		l.unread()
		break
	}

	if value == "-" {
		return token{kind: tokenIllegal}
	}

	if r := l.read(); r != '.' {
		return token{kind: tokenNumber, value: value}
	}

	value += "."

	for {
		r := l.read()

		if unicode.IsDigit(r) {
			value += string(r)
			continue
		}

		l.unread()
		break
	}

	// things like 5. without decimal numbers should be illegal
	if strings.HasSuffix(value, ".") {
		return token{kind: tokenIllegal}
	}

	return token{kind: tokenNumber, value: value}
}

func (k tokenKind) String() string {
	switch k {
	case tokenEOF:
		return "EOF"
	case tokenIllegal:
		return "Illegal"
	case tokenWhitespace:
		return "Whitespace"
	case tokenComma:
		return "Comma"
	case tokenSquareBracketOpen:
		return "SquareBracketOpen"
	case tokenSquareBracketClose:
		return "SquareBracketClose"
	case tokenIdent:
		return "Ident"
	case tokenString:
		return "String"
	case tokenNumber:
		return "Number"
	case tokenTrue:
		return "True"
	case tokenFalse:
		return "False"
	case tokenNull:
		return "Null"
	default:
		return ""
	}
}
