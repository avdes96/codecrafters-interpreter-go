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
		parenthesise(sb, e.operator.Lexeme, e.right)
	case *Literal:
		sb.WriteString(e.String())
	case *Grouping:
		parenthesise(sb, "grouping", e.expression)
	}
}

func parenthesise(sb *strings.Builder, name string, exprs ...Expr) {
	sb.WriteString("(")
	sb.WriteString(name)
	for _, expr := range exprs {
		construct(expr, sb)
	}
	sb.WriteString(")")
}
