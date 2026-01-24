package parser

import (
	"github.com/codecrafters-io/interpreter-starter-go/app/ast"
	"github.com/codecrafters-io/interpreter-starter-go/app/errs"
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

type ParseError struct{}

func NewParseError(t *token.Token, message string) *ParseError {
	errs.ErrorAtToken(t, message)
	return &ParseError{}
}

func (p *Parser) ParseExpression() ast.Expr {
	expr, parseErr := p.expression()
	if parseErr != nil {
		return nil
	}
	return expr
}

func (p *Parser) ParseProgram() []ast.Stmt {
	statements := []ast.Stmt{}
	for !p.isAtEnd() {
		s, parseErr := p.declaration()
		if parseErr != nil {
			return nil
		}
		statements = append(statements, s)
	}
	return statements
}

func (p *Parser) declaration() (ast.Stmt, *ParseError) {
	if p.match(token.VAR) {
		varDec, parseErr := p.varDeclaration()
		if parseErr != nil {
			p.synchronise()
			return nil, parseErr
		}
		return varDec, nil
	}
	statement, parseErr := p.statement()
	if parseErr != nil {
		p.synchronise()
		return nil, parseErr
	}
	return statement, nil
}

func (p *Parser) varDeclaration() (ast.Stmt, *ParseError) {
	name, parseErr := p.consume(token.IDENTIFIER, "Expect variable name.")
	if parseErr != nil {
		return nil, parseErr
	}
	var initialiser ast.Expr = nil
	if p.match(token.EQUAL) {
		initialiser, parseErr = p.expression()
		if parseErr != nil {
			return nil, parseErr
		}
	}
	p.consume(token.SEMICOLON, "Expect ';' after variable declaration.")
	return ast.NewVarStmt(name, initialiser), nil
}

func (p *Parser) statement() (ast.Stmt, *ParseError) {
	if p.match(token.PRINT) {
		return p.printStatement()
	}
	return p.expressionStatement()
}

func (p *Parser) printStatement() (ast.Stmt, *ParseError) {
	expression, parseErr := p.expression()
	if parseErr != nil {
		return nil, parseErr
	}
	p.consume(token.SEMICOLON, "Expect ';' after value.")
	return ast.NewPrintStmt(expression), nil
}

func (p *Parser) expressionStatement() (ast.Stmt, *ParseError) {
	expression, parseErr := p.expression()
	if parseErr != nil {
		return nil, parseErr
	}
	p.consume(token.SEMICOLON, "Expect ';' after value.")
	return ast.NewExpressionStmt(expression), nil
}

func (p *Parser) expression() (ast.Expr, *ParseError) {
	return p.equality()
}

func (p *Parser) equality() (ast.Expr, *ParseError) {
	expr, parseErr := p.comparsion()
	if parseErr != nil {
		return nil, parseErr
	}
	for p.match(token.BANG_EQUAL, token.EQUAL_EQUAL) {
		operator := p.previous()
		right, parseErr := p.term()
		if parseErr != nil {
			return nil, parseErr
		}
		expr = ast.NewBinary(expr, operator, right)
	}
	return expr, nil
}
func (p *Parser) comparsion() (ast.Expr, *ParseError) {
	expr, parseErr := p.term()
	if parseErr != nil {
		return nil, parseErr
	}
	for p.match(token.GREATER, token.GREATER_EQUAL, token.LESS, token.LESS_EQUAL) {
		operator := p.previous()
		right, parseErr := p.term()
		if parseErr != nil {
			return nil, parseErr
		}
		expr = ast.NewBinary(expr, operator, right)
	}
	return expr, nil
}

func (p *Parser) term() (ast.Expr, *ParseError) {
	expr, parseErr := p.factor()
	if parseErr != nil {
		return nil, parseErr
	}
	for p.match(token.MINUS, token.PLUS) {
		operator := p.previous()
		right, parseErr := p.factor()
		if parseErr != nil {
			return nil, parseErr
		}
		expr = ast.NewBinary(expr, operator, right)
	}
	return expr, nil
}

func (p *Parser) factor() (ast.Expr, *ParseError) {
	expr, parseErr := p.unary()
	if parseErr != nil {
		return nil, parseErr
	}
	for p.match(token.SLASH, token.STAR) {
		operator := p.previous()
		right, parseErr := p.unary()
		if parseErr != nil {
			return nil, parseErr
		}
		expr = ast.NewBinary(expr, operator, right)
	}
	return expr, nil
}

func (p *Parser) unary() (ast.Expr, *ParseError) {
	if p.match(token.BANG, token.MINUS) {
		operator := p.previous()
		right, parseErr := p.unary()
		if parseErr != nil {
			return nil, parseErr
		}
		return ast.NewUnary(operator, right), nil
	}
	return p.primary()
}
func (p *Parser) primary() (ast.Expr, *ParseError) {
	if p.match(token.NUMBER, token.STRING) {
		return ast.NewLiteral(p.previous().Literal), nil
	} else if p.match(token.TRUE) {
		return ast.NewLiteral(true), nil
	} else if p.match(token.FALSE) {
		return ast.NewLiteral(false), nil
	} else if p.match(token.NIL) {
		return ast.NewLiteral(nil), nil
	} else if p.match(token.LEFT_PAREN) {
		expr, parseErr := p.expression()
		if parseErr != nil {
			return nil, parseErr
		}

		_, parseErr = p.consume(token.RIGHT_PAREN, "Expect ')' after expression.")
		if parseErr != nil {
			return nil, parseErr
		}
		return ast.NewGrouping(expr), nil
	} else if p.match(token.IDENTIFIER) {
		return ast.NewVariable(p.previous()), nil
	}

	return nil, NewParseError(p.peek(), "Expect expression.")
}

var statementBoundaryTokens = map[token.TokenType]struct{}{
	token.CLASS:  {},
	token.FUN:    {},
	token.VAR:    {},
	token.FOR:    {},
	token.IF:     {},
	token.WHILE:  {},
	token.PRINT:  {},
	token.RETURN: {},
}

func (p *Parser) synchronise() {
	p.advance()
	for !p.isAtEnd() {
		if p.previous().Type == token.SEMICOLON {
			return
		}
		if _, ok := statementBoundaryTokens[p.peek().Type]; ok {
			return
		}
		p.advance()
	}
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
	return p.peek().Type == tokenType
}

func (p *Parser) advance() *token.Token {
	p.current++
	return p.previous()
}

func (p *Parser) peek() *token.Token {
	if p.current >= len(p.tokens) {
		return nil
	}
	return p.tokens[p.current]
}

func (p *Parser) isAtEnd() bool {
	return p.peek().Type == token.EOF
}

func (p *Parser) previous() *token.Token {
	return p.tokens[p.current-1]
}

func (p *Parser) consume(tokenType token.TokenType, message string) (*token.Token, *ParseError) {
	if p.check(tokenType) {
		return p.advance(), nil
	}
	return nil, NewParseError(p.peek(), message)
}
