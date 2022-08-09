package ast

import (
	"fmt"

	"github.com/fadyZohdy/gLox/pkg/scanner"
)

type Expr interface {
	accept(visitor Visitor)
}

type Visitor interface {
	VisitTernaryExpr(expr *Ternary)
	VisitBinaryExpr(expr *Binary)
	VisitGroupingExpr(expr *Grouping)
	VisitLiteralExpr(expr *Literal)
	VisitUnaryExpr(expr *Unary)
}

type Ternary struct {
	condition, trueBranch, falseBranch Expr
}

func (expr *Ternary) String() string {
	return fmt.Sprintf("(?: %v ? %v : %v)", expr.condition, expr.trueBranch, expr.falseBranch)
}

func (expr *Ternary) accept(visitor Visitor) {
	visitor.VisitTernaryExpr(expr)
}

type Binary struct {
	left     Expr
	operator scanner.Token
	right    Expr
}

func (expr *Binary) accept(visitor Visitor) {
	visitor.VisitBinaryExpr(expr)
}

func (expr *Binary) String() string {
	return fmt.Sprintf("(%v %v %v)", expr.left, expr.operator.Lexeme, expr.right)
}

type Grouping struct {
	expression Expr
}

func (expr *Grouping) accept(visitor Visitor) {
	visitor.VisitGroupingExpr(expr)
}

func (expr *Grouping) String() string {
	return fmt.Sprintf("(%v)", expr.expression)
}

type Unary struct {
	operator scanner.Token
	right    Expr
}

func (expr *Unary) accept(visitor Visitor) {
	visitor.VisitUnaryExpr(expr)
}

func (expr *Unary) String() string {
	return fmt.Sprintf("(%v %v)", expr.operator.Lexeme, expr.right)
}

type Literal struct {
	value interface{}
}

func (expr *Literal) accept(visitor Visitor) {
	visitor.VisitLiteralExpr(expr)
}

func (expr *Literal) String() string {
	return fmt.Sprintf("%v", expr.value)
}
