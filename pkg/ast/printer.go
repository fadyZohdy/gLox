package ast

import (
	"fmt"
)

type AstPrinter struct {
	Repr string
}

func (p *AstPrinter) Print(expr Expr) {
	// clean internal state of the printer
	p.Repr = ""
	expr.accept(p)
}

func (p *AstPrinter) VisitTernaryExpr(expr *Ternary) {
	p.parenthesize("?", expr.condition, expr.trueBranch, expr.falseBranch)
}

func (p *AstPrinter) VisitBinaryExpr(expr *Binary) {
	p.parenthesize(expr.operator.Lexeme,
		expr.left, expr.right)
}

func (p *AstPrinter) VisitGroupingExpr(expr *Grouping) {
	p.parenthesize("group", expr.expression)
}

func (p *AstPrinter) VisitUnaryExpr(expr *Unary) {
	p.parenthesize(expr.operator.Lexeme, expr.right)
}

func (p *AstPrinter) VisitLiteralExpr(expr *Literal) {
	if expr.value == nil {
		p.Repr += "nil"
	} else {
		p.Repr += fmt.Sprintf("%v", expr.value)
	}
}

func (p *AstPrinter) parenthesize(name string, exprs ...Expr) {
	p.Repr += "(" + name
	for _, expr := range exprs {
		p.Repr += " "
		expr.accept(p)
	}
	p.Repr += ")"
}
