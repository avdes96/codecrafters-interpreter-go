package scanner

import "github.com/codecrafters-io/interpreter-starter-go/app/token"

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
	switch s.advance() {
	case '(':
		s.addTokenOfType(token.LEFT_PAREN)
	case ')':
		s.addTokenOfType(token.RIGHT_PAREN)
	case '{':
		s.addTokenOfType(token.LEFT_BRACE)
	case '}':
		s.addTokenOfType(token.RIGHT_BRACE)
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
