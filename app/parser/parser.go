package parser

import (
	"github.com/codecrafters-io/interpreter-starter-go/app/ast"
	"github.com/codecrafters-io/interpreter-starter-go/app/token"
)

type Parser struct {
	tokens  []token.Token
	current int
}

func NewParser(tokens []token.Token) *Parser {
	return &Parser{
		tokens:  tokens,
		current: 0,
	}
}

func (p *Parser) Parse() ast.Expr {
	if p.match(token.TRUE) {
		return ast.NewLiteral(true)
	} else if p.match(token.FALSE) {
		return ast.NewLiteral(false)
	}
	return ast.NewLiteral(nil)
}

func (p *Parser) match(tokenTypes ...token.TokenType) bool {
	for _, t := range tokenTypes {
		if p.check(t) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) check(tokenType token.TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek() == tokenType
}

func (p *Parser) advance() token.Token {
	tmp := p.tokens[p.current]
	p.current++
	return tmp
}

func (p *Parser) peek() token.TokenType {
	if p.current >= len(p.tokens) {
		return -1
	}
	return p.tokens[p.current].Type
}

func (p *Parser) isAtEnd() bool {
	return p.peek() == token.EOF
}
