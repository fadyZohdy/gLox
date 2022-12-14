package ast

import "github.com/fadyZohdy/gLox/pkg/scanner"

type Expr interface {
	accept(visitor Visitor) any
}

type Binary struct {
	left     Expr
	operator scanner.Token
	right    Expr
}

func (expr *Binary) accept(visitor Visitor) any {
	return visitor.VisitBinaryExpr(expr)
}

type Grouping struct {
	expression Expr
}

func (expr *Grouping) accept(visitor Visitor) any {
	return visitor.VisitGroupingExpr(expr)
}

type Literal struct {
	value interface{}
}

func (expr *Literal) accept(visitor Visitor) any {
	return visitor.VisitLiteralExpr(expr)
}

type Unary struct {
	operator scanner.Token
	right    Expr
}

func (expr *Unary) accept(visitor Visitor) any {
	return visitor.VisitUnaryExpr(expr)
}

type Ternary struct {
	condition   Expr
	trueBranch  Expr
	falseBranch Expr
}

func (expr *Ternary) accept(visitor Visitor) any {
	return visitor.VisitTernaryExpr(expr)
}

type Logical struct {
	operator    scanner.Token
	left, right Expr
}

func (expr *Logical) accept(v Visitor) any {
	return v.VisitLogicalExpr(expr)
}

type Variable struct {
	name scanner.Token
}

func (expr *Variable) accept(visitor Visitor) any {
	return visitor.VisitVariableExpr(expr)
}

type Assign struct {
	name  scanner.Token
	value Expr
}

func (expr *Assign) accept(visitor Visitor) any {
	return visitor.VisitAssignExpr(expr)
}

type Call struct {
	callee    Expr
	paren     scanner.Token
	arguments []Expr
}

func (expr *Call) accept(visitor Visitor) any {
	return visitor.VisitCallExpr(expr)
}

type Get struct {
	instance Expr
	name     scanner.Token
}

func (expr *Get) accept(visitor Visitor) any {
	return visitor.VisitGetExpr(expr)
}

type Set struct {
	object Expr
	name   scanner.Token
	value  Expr
}

func (expr *Set) accept(visitor Visitor) any {
	return visitor.VisitSetExpr(expr)
}

type This struct {
	keyword scanner.Token
}

func (expr *This) accept(visitor Visitor) any {
	return visitor.VisitThisExpr(expr)
}
