package ast

import "github.com/fadyZohdy/gLox/pkg/scanner"

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

type Grouping struct {
	expression Expr
}

func (expr *Grouping) accept(visitor Visitor) {
	visitor.VisitGroupingExpr(expr)
}

type Literal struct {
	value interface{}
}

func (expr *Literal) accept(visitor Visitor) {
	visitor.VisitLiteralExpr(expr)
}

type Unary struct {
	operator scanner.Token
	right    Expr
}

func (expr *Unary) accept(visitor Visitor) {
	visitor.VisitUnaryExpr(expr)
}
