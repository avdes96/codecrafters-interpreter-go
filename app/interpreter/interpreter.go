package interpreter

import (
	"fmt"
	"os"

	"github.com/codecrafters-io/interpreter-starter-go/app/ast"
	"github.com/codecrafters-io/interpreter-starter-go/app/errs"
	"github.com/codecrafters-io/interpreter-starter-go/app/token"
	"github.com/codecrafters-io/interpreter-starter-go/app/utils"
)

type Interpreter struct {
	env *environment
}

func NewInterpreter() *Interpreter {
	return &Interpreter{
		env: NewBaseEnviroment(),
	}
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

func (i *Interpreter) Interpret(statements []ast.Stmt) {
	for _, statement := range statements {
		runtimeErr := i.execute(statement)
		if runtimeErr != nil {
			return
		}
	}
}

func (i *Interpreter) execute(statement ast.Stmt) *RuntimeError {
	switch s := statement.(type) {
	case *ast.PrintStmt:
		return i.executePrintStatement(s)
	case *ast.ExpressionStmt:
		return i.executeExpressionStatement(s)
	case *ast.VarStmt:
		return i.executeVarStatement(s)
	case *ast.BlockStmt:
		return i.executeBlockStatement(s)
	case *ast.IfStmt:
		return i.executeIfStatement(s)
	}
	return nil
}

func (i *Interpreter) executePrintStatement(statement *ast.PrintStmt) *RuntimeError {
	eval, runtimeErr := i.evaluate(statement.Expression)
	if runtimeErr != nil {
		return runtimeErr
	}
	fmt.Println(utils.Stringify(eval))
	return nil
}

func (i *Interpreter) executeExpressionStatement(statement *ast.ExpressionStmt) *RuntimeError {
	_, runtimeErr := i.evaluate(statement.Expression)
	return runtimeErr
}

func (i *Interpreter) executeVarStatement(statement *ast.VarStmt) *RuntimeError {
	var value any = nil
	var runtimeErr *RuntimeError = nil
	if statement.Initialiser != nil {
		value, runtimeErr = i.evaluate(statement.Initialiser)
		if runtimeErr != nil {
			return runtimeErr
		}
	}
	i.env.define(statement.Name.Lexeme, value)
	return nil
}

func (i *Interpreter) executeBlockStatement(statement *ast.BlockStmt) *RuntimeError {
	prev := i.env
	defer func() { i.env = prev }()
	i.env = NewSubEnvironment(prev)
	for _, stmt := range statement.Stmts {
		if runtimeErr := i.execute(stmt); runtimeErr != nil {
			return runtimeErr
		}
	}
	return nil
}

func (i *Interpreter) executeIfStatement(statement *ast.IfStmt) *RuntimeError {
	val, runtimeErr := i.evaluate(statement.Condition)
	if runtimeErr != nil {
		return runtimeErr
	}
	if isTruthy(val) {
		return i.execute(statement.ThenBranch)
	}
	if statement.ElseBranch != nil {
		return i.execute(statement.ElseBranch)
	}
	return nil
}

func (i *Interpreter) Evaluate(expression ast.Expr) any {
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
	case *ast.Variable:
		return i.env.get(e.Name)
	case *ast.Assign:
		value, runtimeErr := i.evaluate(e.Expression)
		if runtimeErr != nil {
			return nil, runtimeErr
		}
		runtimeErr = i.env.assign(e.Name, value)
		if runtimeErr != nil {
			return nil, runtimeErr
		}
		return value, nil
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
