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
