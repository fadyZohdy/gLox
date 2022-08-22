package ast

import (
	"fmt"
)

type AstPrinter struct{}

func (p *AstPrinter) Print(expr Expr) string {
	if expr == nil {
		return ""
	}
	res := expr.accept(p)
	if res_s, ok := res.(string); ok {
		return res_s
	}
	return ""
}

func (p *AstPrinter) VisitVarStmt(stmt *Var) any {
	return p.parenthesize(fmt.Sprintf("var %s", stmt.name.Lexeme), stmt.initializer)
}

func (p *AstPrinter) VisitVariableExpr(expr *Variable) any {
	return expr.name.Lexeme
}

func (p *AstPrinter) VisitExpressionStmt(stmt *Expression) any {
	return stmt.expression.accept(p)
}

func (p *AstPrinter) VisitIfStmt(stmt *If) any {
	return fmt.Sprintf("if (%s) %v else %v", stmt.condition.accept(p), stmt.trueBranch.accept(p), stmt.falseBranch.accept(p))
}

func (p *AstPrinter) VisitBlockStmt(block *Block) any {
	res := "{"
	for _, stmt := range block.statements {
		if str, ok := stmt.accept(p).(string); ok {
			res += str
			res += " "
		}
	}
	res += "}"
	return res
}

func (p *AstPrinter) VisitPrintStmt(stmt *Print) any {
	return p.parenthesize("print", stmt.expression)
}

func (p *AstPrinter) VisitWhileStmt(stmt *While) any {
	return fmt.Sprintf("while (%s) {%s}", stmt.condition.accept(p), stmt.body.accept(p))
}

func (p *AstPrinter) VisitBreakStatement(stmt *Break) any {
	return "break"
}

func (p *AstPrinter) VisitFunctionStmt(stmt *Function) any {
	res := fmt.Sprintf("fun %s(", stmt.name.Lexeme)
	for _, p := range stmt.params {
		res += p.Lexeme + ","
	}
	res += ") { "
	for _, s := range stmt.body {
		res += fmt.Sprintf("%v", s.accept(p))
	}
	res += " }"
	return res
}

func (p *AstPrinter) VisitReturnStmt(stmt *Return) any {
	return fmt.Sprintf("return %v", stmt.value.accept(p))
}

func (p *AstPrinter) VisitCallExpr(expr *Call) any {
	if callee, ok := expr.callee.accept(p).(string); ok {
		return p.parenthesize(callee, expr.arguments...)
	}
	return ""
}

func (p *AstPrinter) VisitAssignExpr(expr *Assign) any {
	return p.parenthesize(expr.name.Lexeme+" = ", expr.value)
}

func (p *AstPrinter) VisitLogicalExpr(expr *Logical) any {
	return p.parenthesize(expr.operator.Lexeme, expr.left, expr.right)
}

func (p *AstPrinter) VisitTernaryExpr(expr *Ternary) any {
	return p.parenthesize("?", expr.condition, expr.trueBranch, expr.falseBranch)
}

func (p *AstPrinter) VisitBinaryExpr(expr *Binary) any {
	return p.parenthesize(expr.operator.Lexeme,
		expr.left, expr.right)
}

func (p *AstPrinter) VisitGroupingExpr(expr *Grouping) any {
	return p.parenthesize("group", expr.expression)
}

func (p *AstPrinter) VisitUnaryExpr(expr *Unary) any {
	return p.parenthesize(expr.operator.Lexeme, expr.right)
}

func (p *AstPrinter) VisitLiteralExpr(expr *Literal) any {
	if expr.value != nil {
		return fmt.Sprintf("%v", expr.value)
	}
	return "nil"
}

func (p *AstPrinter) parenthesize(name string, exprs ...Expr) any {
	s := ""
	s += "(" + name
	for _, expr := range exprs {
		s += " "
		if expr != nil {
			ss := expr.accept(p)
			if sss, ok := ss.(string); ok {
				s += sss
			}
		} else {
			s += "nil"
		}
	}
	s += ")"
	return s
}
