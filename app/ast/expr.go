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

func (b *Binary) exprNode() {}

type Unary struct {
	operator token.Token
	right    Expr
}

func (u *Unary) exprNode() {}

type Literal struct {
	value any
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
