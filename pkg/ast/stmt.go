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
