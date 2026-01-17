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
	return ast.NewLiteral(nil)
}
