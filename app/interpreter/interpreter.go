package interpreter

import (
	"fmt"
	"os"

	"github.com/codecrafters-io/interpreter-starter-go/app/ast"
)

type Interpreter struct{}

func NewInterpreter() *Interpreter {
	return &Interpreter{}
}

func (i *Interpreter) Interpret(expression ast.Expr) any {
	return i.evaluate(expression)
}

func (i *Interpreter) evaluate(expression ast.Expr) any {
	switch e := expression.(type) {
	case *ast.Literal:
		return e.Value
	case *ast.Grouping:
		return i.evaluate(e.Expression)
	}
	fmt.Fprintf(os.Stderr, "Unexpected type: %T\n", expression)
	os.Exit(1)
	return nil
}
