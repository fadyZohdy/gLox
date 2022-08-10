package ast

import (
	"fmt"

	"github.com/fadyZohdy/gLox/pkg/scanner"
)

type Expr interface {
	accept(visitor Visitor) any
}

type Visitor interface {
	VisitTernaryExpr(expr *Ternary) any
	VisitBinaryExpr(expr *Binary) any
	VisitGroupingExpr(expr *Grouping) any
	VisitLiteralExpr(expr *Literal) any
	VisitUnaryExpr(expr *Unary) any
}

type Ternary struct {
	condition, trueBranch, falseBranch Expr
}

func (expr *Ternary) String() string {
	return fmt.Sprintf("(?: %v ? %v : %v)", expr.condition, expr.trueBranch, expr.falseBranch)
}

func (expr *Ternary) accept(visitor Visitor) any {
	return visitor.VisitTernaryExpr(expr)
}

type Binary struct {
	left     Expr
	operator scanner.Token
	right    Expr
}

func (expr *Binary) accept(visitor Visitor) any {
	return visitor.VisitBinaryExpr(expr)
}

func (expr *Binary) String() string {
	return fmt.Sprintf("(%v %v %v)", expr.left, expr.operator.Lexeme, expr.right)
}

type Grouping struct {
	expression Expr
}

func (expr *Grouping) accept(visitor Visitor) any {
	return visitor.VisitGroupingExpr(expr)
}

func (expr *Grouping) String() string {
	return fmt.Sprintf("(%v)", expr.expression)
}

type Unary struct {
	operator scanner.Token
	right    Expr
}

func (expr *Unary) accept(visitor Visitor) any {
	return visitor.VisitUnaryExpr(expr)
}

func (expr *Unary) String() string {
	return fmt.Sprintf("(%v %v)", expr.operator.Lexeme, expr.right)
}

type Literal struct {
	value interface{}
}

func (expr *Literal) accept(visitor Visitor) any {
	return visitor.VisitLiteralExpr(expr)
}

func (expr *Literal) String() string {
	return fmt.Sprintf("%v", expr.value)
}
