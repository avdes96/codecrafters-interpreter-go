package scanner

import (
	"fmt"

	"github.com/codecrafters-io/interpreter-starter-go/app/errs"
	"github.com/codecrafters-io/interpreter-starter-go/app/token"
)

type Scanner struct {
	source  string
	tokens  []token.Token
	start   int
	current int
}

func NewScanner(source string) *Scanner {
	return &Scanner{
		source: source,
		tokens: []token.Token{},
	}
}

func (s *Scanner) ScanTokens() []token.Token {
	for s.current < len(s.source) {
		s.start = s.current
		s.scanToken()
	}
	s.tokens = append(s.tokens, token.NewToken(token.EOF, "", nil))
	return s.tokens
}

func (s *Scanner) scanToken() {
	switch c := s.advance(); c {
	case '(':
		s.addTokenOfType(token.LEFT_PAREN)
	case ')':
		s.addTokenOfType(token.RIGHT_PAREN)
	case '{':
		s.addTokenOfType(token.LEFT_BRACE)
	case '}':
		s.addTokenOfType(token.RIGHT_BRACE)
	case ',':
		s.addTokenOfType(token.COMMA)
	case '.':
		s.addTokenOfType(token.DOT)
	case '-':
		s.addTokenOfType(token.MINUS)
	case '+':
		s.addTokenOfType(token.PLUS)
	case '*':
		s.addTokenOfType(token.STAR)
	case ';':
		s.addTokenOfType(token.SEMICOLON)
	default:
		errs.Error(fmt.Sprintf("Unexpected character: %c", c))
	}
}

func (s *Scanner) advance() rune {
	c := s.source[s.current]
	s.current++
	return rune(c)
}

func (s *Scanner) addTokenOfType(tokenType token.TokenType) {
	s.addToken(tokenType, nil)
}

func (s *Scanner) addToken(tokenType token.TokenType, literal any) {
	lexeme := s.source[s.start:s.current]
	s.tokens = append(s.tokens, token.NewToken(tokenType, lexeme, literal))
}
