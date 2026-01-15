package scanner

import (
	"fmt"
	"strconv"

	"github.com/codecrafters-io/interpreter-starter-go/app/errs"
	"github.com/codecrafters-io/interpreter-starter-go/app/token"
)

type Scanner struct {
	source  string
	tokens  []token.Token
	start   int
	current int
	line    int
}

func NewScanner(source string) *Scanner {
	return &Scanner{
		source:  source,
		tokens:  []token.Token{},
		start:   0,
		current: 0,
		line:    1,
	}
}

func (s *Scanner) ScanTokens() []token.Token {
	for !s.atEnd() {
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
	case '=':
		if s.match('=') {
			s.addTokenOfType(token.EQUAL_EQUAL)
			break
		}
		s.addTokenOfType(token.EQUAL)
	case '!':
		if s.match('=') {
			s.addTokenOfType(token.BANG_EQUAL)
			break
		}
		s.addTokenOfType(token.BANG)
	case '<':
		if s.match('=') {
			s.addTokenOfType(token.LESS_EQUAL)
			break
		}
		s.addTokenOfType(token.LESS)
	case '>':
		if s.match('=') {
			s.addTokenOfType(token.GREATER_EQUAL)
			break
		}
		s.addTokenOfType(token.GREATER)
	case '/':
		if s.match('/') {
			s.consumeLine()
			break
		}
		s.addTokenOfType(token.SLASH)
	case ' ', '\t', '\r':
		// Ignore the whitespace
	case '\n':
		s.line++
	case '"':
		s.addString()
	default:
		if isDigit(c) {
			s.consumeNumber()
			break
		}
		errs.Error(s.line, fmt.Sprintf("Unexpected character: %c", c))
	}
}

func (s *Scanner) atEnd() bool {
	return s.current >= len(s.source)
}

func (s *Scanner) advance() rune {
	c := s.source[s.current]
	s.current++
	return rune(c)
}

func (s *Scanner) match(expected rune) bool {
	if s.atEnd() || rune(s.source[s.current]) != expected {
		return false
	}
	s.current++
	return true
}

func (s *Scanner) addTokenOfType(tokenType token.TokenType) {
	s.addToken(tokenType, nil)
}

func (s *Scanner) addToken(tokenType token.TokenType, literal any) {
	lexeme := s.source[s.start:s.current]
	s.tokens = append(s.tokens, token.NewToken(tokenType, lexeme, literal))
}

func (s *Scanner) consumeLine() {
	for !s.atEnd() && s.peek() != '\n' {
		s.current++
	}
}

func (s *Scanner) peek() rune {
	if s.atEnd() {
		return 0
	}
	return rune(s.source[s.current])
}

func (s *Scanner) peekNext() rune {
	if s.current+1 >= len(s.source) {
		return 0
	}
	return rune(s.source[s.current+1])
}

func (s *Scanner) addString() {
	for !s.atEnd() && s.peek() != '"' {
		s.current++
	}
	if s.atEnd() {
		errs.Error(s.line, "Unterminated string.")
		return
	}
	s.advance()
	s.addToken(token.STRING, s.source[s.start+1:s.current-1])
}

func isDigit(c rune) bool {
	return c >= '0' && c <= '9'
}

func (s *Scanner) consumeNumber() {
	for isDigit(s.peek()) {
		s.advance()
	}
	if s.peek() == '.' && isDigit(s.peekNext()) {
		s.advance()
		for isDigit(s.peek()) {
			s.advance()
		}
	}
	numAsString := s.source[s.start:s.current]
	literal, err := strconv.ParseFloat(numAsString, 64)
	if err != nil {
		errs.Error(s.line, fmt.Sprintf("Error parsing %s to number", numAsString))
		return
	}
	s.addToken(token.NUMBER, literal)
}
