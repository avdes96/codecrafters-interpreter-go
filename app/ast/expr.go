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
	left, right Expr
	operator    token.Token
}

func NewBinary(left, right Expr, operator token.Token) *Binary {
	return &Binary{
		left:     left,
		right:    right,
		operator: operator,
	}
}

func (b *Binary) exprNode() {}

type Unary struct {
	operator token.Token
	right    Expr
}

func NewUnary(right Expr, operator token.Token) *Unary {
	return &Unary{
		right:    right,
		operator: operator,
	}
}

func (u *Unary) exprNode() {}

type Literal struct {
	value any
}

func NewLiteral(value any) *Literal {
	return &Literal{value: value}
}

func (l *Literal) String() string {
	var str string
	switch v := l.value.(type) {
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
	expression Expr
}

func (g *Grouping) exprNode() {}

func NewGrouping(expression *Expr) *Grouping {
	return &Grouping{expression: *expression}
}
