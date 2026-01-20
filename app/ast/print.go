package ast

import (
	"fmt"
	"strings"
)

func Print(expr Expr) {
	sb := &strings.Builder{}
	construct(expr, sb)
	fmt.Println(sb.String())
}

func construct(expr Expr, sb *strings.Builder) {
	switch e := expr.(type) {
	case *Binary:
		parenthesise(sb, e.operator.Lexeme, e.left, e.right)
	case *Unary:
		parenthesise(sb, e.Operator.Lexeme, e.Right)
	case *Literal:
		sb.WriteString(e.String())
	case *Grouping:
		parenthesise(sb, "group", e.Expression)
	}
}

func parenthesise(sb *strings.Builder, name string, exprs ...Expr) {
	sb.WriteString("(")
	sb.WriteString(name)
	for _, expr := range exprs {
		sb.WriteString(" ")
		construct(expr, sb)
	}
	sb.WriteString(")")
}
