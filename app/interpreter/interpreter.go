package interpreter

import "github.com/codecrafters-io/interpreter-starter-go/app/ast"

type Interpreter struct{}

func NewInterpreter() *Interpreter {
	return &Interpreter{}
}

func (i *Interpreter) Interpret(expression ast.Expr) {

}
