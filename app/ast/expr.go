package ast

import (
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
	switch v := l.value.(type) {
	case nil:
		return "null"
	case float64:
		s := strconv.FormatFloat(v, 'f', -1, 64)
		if !strings.ContainsRune(s, '.') {
			s += ".0"
		}
		return s
	case string:
		return v
	}
	return ""
}

func (l *Literal) exprNode() {}

type Grouping struct {
	expression Expr
}

func (g *Grouping) exprNode() {}

func NewGrouping(expression *Expr) *Grouping {
	return &Grouping{expression: *expression}
}
