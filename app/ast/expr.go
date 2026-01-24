package ast

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/codecrafters-io/interpreter-starter-go/app/token"
)

type Expr interface {
	exprNode()
}

type Binary struct {
	Left, Right Expr
	Operator    *token.Token
}

func NewBinary(left Expr, operator *token.Token, right Expr) *Binary {
	return &Binary{
		Left:     left,
		Right:    right,
		Operator: operator,
	}
}

func (b *Binary) exprNode() {}

type Unary struct {
	Operator *token.Token
	Right    Expr
}

func NewUnary(operator *token.Token, right Expr) *Unary {
	return &Unary{
		Operator: operator,
		Right:    right,
	}
}

func (u *Unary) exprNode() {}

type Literal struct {
	Value any
}

func NewLiteral(value any) *Literal {
	return &Literal{Value: value}
}

func (l *Literal) String() string {
	var str string
	switch v := l.Value.(type) {
	case nil:
		str = "nil"
	case bool:
		str = strconv.FormatBool(v)
	case float64:
		str = strconv.FormatFloat(v, 'f', -1, 64)
		if !strings.ContainsRune(str, '.') {
			str += ".0"
		}
	case string:
		str = v
	default:
		fmt.Fprintf(os.Stderr, "Unsupported literal type %T\n", v)
		os.Exit(1)
	}
	return str
}

func (l *Literal) exprNode() {}

type Grouping struct {
	Expression Expr
}

func (g *Grouping) exprNode() {}

func NewGrouping(expression Expr) *Grouping {
	return &Grouping{Expression: expression}
}

type Variable struct {
	Name *token.Token
}

func (v *Variable) exprNode() {}

func NewVariable(name *token.Token) *Variable {
	return &Variable{Name: name}
}
