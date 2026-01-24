package ast

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
