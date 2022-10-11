package ast

import (
	"fmt"

	"github.com/fadyZohdy/gLox/pkg/scanner"
)

type Stmt interface {
	accept(visitor Visitor) any
}

type Expression struct {
	expression Expr
}

func (stmt *Expression) accept(visitor Visitor) any {
	return visitor.VisitExpressionStmt(stmt)
}

type If struct {
	condition   Expr
	trueBranch  Stmt
	falseBranch Stmt
}

func (stmt If) String() string {
	return fmt.Sprintf("if (%s) {%s} else {%s}", stmt.condition, stmt.trueBranch, stmt.falseBranch)
}

func (stmt *If) accept(v Visitor) any {
	return v.VisitIfStmt(stmt)
}

type Print struct {
	expression Expr
}

func (p Print) String() string {
	return fmt.Sprintf("print %v", p.expression)
}

func (stmt *Print) accept(visitor Visitor) any {
	return visitor.VisitPrintStmt(stmt)
}

type Var struct {
	name        scanner.Token
	initializer Expr
}

func (v Var) String() string {
	return fmt.Sprintf("var %v = %v", v.name.Lexeme, v.initializer)
}

func (stmt *Var) accept(visitor Visitor) any {
	return visitor.VisitVarStmt(stmt)
}

type Block struct {
	statements []Stmt
}

func (stmt *Block) accept(visitor Visitor) any {
	return visitor.VisitBlockStmt(stmt)
}

type While struct {
	condition Expr
	body      Stmt
}

func (stmt *While) accept(v Visitor) any {
	return v.VisitWhileStmt(stmt)
}

type Break struct{}

func (stmt *Break) accept(v Visitor) any {
	return v.VisitBreakStatement(stmt)
}

type Function struct {
	name         scanner.Token
	params       []scanner.Token
	body         []Stmt
	functionType FunctionType
}

func (stmt Function) accept(v Visitor) any {
	return v.VisitFunctionStmt(&stmt)
}

func (stmt *Function) isAnon() bool {
	return stmt.name.Lexeme == ""
}

type Return struct {
	keyword scanner.Token
	value   Expr
}

func (stmt *Return) accept(v Visitor) any {
	return v.VisitReturnStmt(stmt)
}

type Class struct {
	name    scanner.Token
	methods []*Function
}

func (stmt *Class) accept(v Visitor) any {
	return v.VisitClassStmt(stmt)
}

func (c Class) String() string {
	return fmt.Sprintf("<class %s>", c.name.Lexeme)
}
