package token

import "fmt"

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
	EOF
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
	switch t.Literal.(type) {
	case nil:
		literal = "null"
	case int:
		literal = "NUMBER"
	case string:
		literal = "STRING"
	}
	return fmt.Sprintf("%s %s %s", t.Type, t.Lexeme, literal)
}
