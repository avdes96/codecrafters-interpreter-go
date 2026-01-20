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
	case *ast.Binary:
		return i.evaluateBinary(e)
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

func (i *Interpreter) evaluateBinary(binary *ast.Binary) any {
	left := i.evaluate(binary.Left)
	right := i.evaluate(binary.Right)
	switch binary.Operator.Type {
	case token.MINUS:
		return left.(float64) - right.(float64)
	case token.STAR:
		return left.(float64) * right.(float64)
	case token.SLASH:
		return left.(float64) / right.(float64)
	case token.PLUS:
		if isFloat64(left) && isFloat64(right) {
			return left.(float64) + right.(float64)
		}
		if isString(left) && isString(right) {
			return left.(string) + right.(string)
		}
	case token.GREATER:
		return left.(float64) > right.(float64)
	case token.GREATER_EQUAL:
		return left.(float64) >= right.(float64)
	case token.LESS:
		return left.(float64) < right.(float64)
	case token.LESS_EQUAL:
		return left.(float64) <= right.(float64)
	case token.EQUAL_EQUAL:
		return isEqual(left, right)
	case token.BANG_EQUAL:
		return !isEqual(left, right)
	}
	return nil
}

func isFloat64(val any) bool {
	_, ok := val.(float64)
	return ok
}

func isString(val any) bool {
	_, ok := val.(string)
	return ok
}

func isBool(val any) bool {
	_, ok := val.(bool)
	return ok
}

func isEqual(a any, b any) bool {
	if a == nil && b == nil {
		return true
	}
	if isString(a) && isString(b) {
		return a.(string) == b.(string)
	}
	if isBool(a) && isBool(b) {
		return a.(bool) == b.(bool)
	}
	if isFloat64(a) && isFloat64(b) {
		return a.(float64) == b.(float64)
	}
	return false
}
