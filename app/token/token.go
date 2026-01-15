package token

import (
	"fmt"
	"strconv"
	"strings"
)

//go:generate stringer -type=TokenType
type TokenType int

const (
	LEFT_PAREN TokenType = iota
	RIGHT_PAREN
	LEFT_BRACE
	RIGHT_BRACE
	COMMA
	DOT
	MINUS
	PLUS
	STAR
	SEMICOLON
	EQUAL
	EQUAL_EQUAL
	BANG
	BANG_EQUAL
	LESS
	LESS_EQUAL
	GREATER
	GREATER_EQUAL
	SLASH
	STRING
	EOF
	NUMBER
	IDENTIFIER
)

type Token struct {
	Type    TokenType
	Lexeme  string
	Literal any
}

func NewToken(t TokenType, lexeme string, literal any) Token {
	return Token{
		Type:    t,
		Lexeme:  lexeme,
		Literal: literal,
	}
}

func (t Token) String() string {
	var literal string
	switch lit := t.Literal.(type) {
	case nil:
		literal = "null"
	case float64:
		literal = strconv.FormatFloat(lit, 'f', -1, 64)
		if !strings.ContainsRune(literal, '.') {
			literal += ".0"
		}
	case string:
		literal = string(lit)
	}
	return fmt.Sprintf("%s %s %s", t.Type, t.Lexeme, literal)
}
