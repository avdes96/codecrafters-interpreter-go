package interpreter

import (
	"fmt"
	"os"

	"github.com/codecrafters-io/interpreter-starter-go/app/ast"
	"github.com/codecrafters-io/interpreter-starter-go/app/token"
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
	case *ast.Unary:
		return i.evaluateUnary(e)
	case *ast.Literal:
		return e.Value
	case *ast.Grouping:
		return i.evaluate(e.Expression)
	}
	fmt.Fprintf(os.Stderr, "Unexpected type: %T\n", expression)
	os.Exit(1)
	return nil
}

func (i *Interpreter) evaluateUnary(unary *ast.Unary) any {
	right := i.evaluate(unary.Right)
	switch unary.Operator.Type {
	case token.MINUS:
		return -1 * right.(float64)
	case token.BANG:
		return !isTruthy(right)
	}
	return nil
}

func isTruthy(val any) bool {
	if val == nil {
		return false
	}
	if v, ok := val.(bool); ok {
		return v
	}
	return true
}
