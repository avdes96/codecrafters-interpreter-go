package interpreter

import (
	"fmt"
	"os"

	"github.com/codecrafters-io/interpreter-starter-go/app/ast"
	"github.com/codecrafters-io/interpreter-starter-go/app/errs"
	"github.com/codecrafters-io/interpreter-starter-go/app/token"
)

type Interpreter struct{}

func NewInterpreter() *Interpreter {
	return &Interpreter{}
}

type RuntimeError struct{}

func NewRuntimeError(t *token.Token, message string) *RuntimeError {
	errs.RuntimeError(t.Line, message)
	return &RuntimeError{}
}

var numberOnlyOperators = map[token.TokenType]struct{}{
	token.MINUS:         {},
	token.STAR:          {},
	token.SLASH:         {},
	token.GREATER:       {},
	token.GREATER_EQUAL: {},
	token.LESS:          {},
	token.LESS_EQUAL:    {},
}

func (i *Interpreter) Interpret(expression ast.Expr) any {
	eval, runtimeErr := i.evaluate(expression)
	if runtimeErr != nil {
		return nil
	}
	return eval
}

func (i *Interpreter) evaluate(expression ast.Expr) (any, *RuntimeError) {
	switch e := expression.(type) {
	case *ast.Binary:
		return i.evaluateBinary(e)
	case *ast.Unary:
		return i.evaluateUnary(e)
	case *ast.Literal:
		return e.Value, nil
	case *ast.Grouping:
		return i.evaluate(e.Expression)
	}
	fmt.Fprintf(os.Stderr, "Unexpected type: %T\n", expression)
	os.Exit(1)
	return nil, nil
}

func (i *Interpreter) evaluateUnary(unary *ast.Unary) (any, *RuntimeError) {
	right, runtimeErr := i.evaluate(unary.Right)
	if runtimeErr != nil {
		return nil, runtimeErr
	}
	switch unary.Operator.Type {
	case token.MINUS:
		if !isFloat64(right) {
			return nil, NewRuntimeError(unary.Operator, "Operand must be a number.")
		}
		return -right.(float64), nil
	case token.BANG:
		return !isTruthy(right), nil
	}
	return nil, nil
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

func (i *Interpreter) evaluateBinary(binary *ast.Binary) (any, *RuntimeError) {
	left, runtimeErr := i.evaluate(binary.Left)
	if runtimeErr != nil {
		return nil, runtimeErr
	}
	right, runtimeErr := i.evaluate(binary.Right)
	if runtimeErr != nil {
		return nil, runtimeErr
	}
	if _, ok := numberOnlyOperators[binary.Operator.Type]; ok && !bothNumberOperands(left, right) {
		return nil, NewRuntimeError(binary.Operator, "Operands must be numbers.")
	}
	switch binary.Operator.Type {
	case token.MINUS:
		return left.(float64) - right.(float64), nil
	case token.STAR:
		return left.(float64) * right.(float64), nil
	case token.SLASH:
		return left.(float64) / right.(float64), nil
	case token.PLUS:
		if isFloat64(left) && isFloat64(right) {
			return left.(float64) + right.(float64), nil
		}
		if isString(left) && isString(right) {
			return left.(string) + right.(string), nil
		}
		return nil, NewRuntimeError(binary.Operator, "Operands must be two numbers or two strings.")
	case token.GREATER:
		return left.(float64) > right.(float64), nil
	case token.GREATER_EQUAL:
		return left.(float64) >= right.(float64), nil
	case token.LESS:
		return left.(float64) < right.(float64), nil
	case token.LESS_EQUAL:
		return left.(float64) <= right.(float64), nil
	case token.EQUAL_EQUAL:
		return isEqual(left, right), nil
	case token.BANG_EQUAL:
		return !isEqual(left, right), nil
	}
	return nil, nil
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

func bothNumberOperands(a any, b any) bool {
	return isFloat64(a) && isFloat64(b)
}
