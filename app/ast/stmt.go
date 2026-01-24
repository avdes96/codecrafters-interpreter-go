package ast

import "github.com/codecrafters-io/interpreter-starter-go/app/token"

type Stmt interface {
	stmtNode()
}

type PrintStmt struct {
	Expression Expr
}

func NewPrintStmt(expression Expr) *PrintStmt {
	return &PrintStmt{Expression: expression}
}

func (p *PrintStmt) stmtNode() {}

type ExpressionStmt struct {
	Expression Expr
}

func NewExpressionStmt(expression Expr) *ExpressionStmt {
	return &ExpressionStmt{Expression: expression}
}

func (e *ExpressionStmt) stmtNode() {}

type VarStmt struct {
	Name        *token.Token
	Initialiser Expr
}

func NewVarStmt(name *token.Token, initialiser Expr) *VarStmt {
	return &VarStmt{
		Name:        name,
		Initialiser: initialiser,
	}
}

func (e *VarStmt) stmtNode() {}
