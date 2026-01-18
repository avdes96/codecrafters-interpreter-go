package parser

import (
	"github.com/codecrafters-io/interpreter-starter-go/app/ast"
	"github.com/codecrafters-io/interpreter-starter-go/app/token"
)

type Parser struct {
	tokens  []*token.Token
	current int
}

func NewParser(tokens []*token.Token) *Parser {
	return &Parser{
		tokens:  tokens,
		current: 0,
	}
}

func (p *Parser) Parse() ast.Expr {
	return p.expression()
}

func (p *Parser) expression() ast.Expr {
	return p.equality()
}

func (p *Parser) equality() ast.Expr {
	expr := p.comparsion()
	for p.match(token.BANG_EQUAL, token.EQUAL_EQUAL) {
		operator := p.previous()
		right := p.term()
		expr = ast.NewBinary(expr, operator, right)
	}
	return expr
}
func (p *Parser) comparsion() ast.Expr {
	expr := p.term()
	for p.match(token.GREATER, token.GREATER_EQUAL, token.LESS, token.LESS_EQUAL) {
		operator := p.previous()
		right := p.term()
		expr = ast.NewBinary(expr, operator, right)
	}
	return expr
}

func (p *Parser) term() ast.Expr {
	expr := p.factor()
	for p.match(token.MINUS, token.PLUS) {
		operator := p.previous()
		right := p.factor()
		expr = ast.NewBinary(expr, operator, right)
	}
	return expr
}

func (p *Parser) factor() ast.Expr {
	expr := p.unary()
	for p.match(token.SLASH, token.STAR) {
		operator := p.previous()
		right := p.unary()
		expr = ast.NewBinary(expr, operator, right)
	}
	return expr
}

func (p *Parser) unary() ast.Expr {
	if p.match(token.BANG, token.MINUS) {
		operator := p.previous()
		right := p.unary()
		return ast.NewUnary(operator, right)
	}
	return p.primary()
}
func (p *Parser) primary() ast.Expr {
	if p.match(token.NUMBER, token.STRING) {
		return ast.NewLiteral(p.previous().Literal)
	} else if p.match(token.TRUE) {
		return ast.NewLiteral(true)
	} else if p.match(token.FALSE) {
		return ast.NewLiteral(false)
	} else if p.match(token.LEFT_PAREN) {
		expr := p.Parse()
		p.advance() // Assume correct right paren for now
		return ast.NewGrouping(expr)
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

func (p *Parser) advance() *token.Token {
	p.current++
	return p.previous()
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

func (p *Parser) previous() *token.Token {
	return p.tokens[p.current-1]
}
