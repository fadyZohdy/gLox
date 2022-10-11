package ast

import (
	"github.com/fadyZohdy/gLox/pkg/lib/stack"
	"github.com/fadyZohdy/gLox/pkg/scanner"
)

type Resolver struct {
	interpreter      *Interpreter
	scopes           *stack.Stack[map[string]bool]
	error_reporter   func(int, string, string)
	visitedFunctions *stack.Stack[FunctionType]
	inClass          bool
}

func NewResolver(interpreter *Interpreter, error_reporter func(int, string, string)) *Resolver {
	return &Resolver{
		interpreter, stack.New[map[string]bool](), error_reporter, stack.New[FunctionType](), false,
	}
}

func (r *Resolver) VisitClassStmt(class *Class) any {
	r.inClass = true
	r.declare(class.name)
	r.define(class.name)

	r.beginScope()
	scope := r.scopes.Peek()
	(*scope)["this"] = true

	for _, method := range class.methods {
		if method.name.Lexeme == "init" {
			r.resolveFunction(method, CONSTRUCTOR)
		} else {
			r.resolveFunction(method, METHOD)
		}
	}
	r.endScope()
	r.inClass = false
	return nil
}

func (r *Resolver) VisitBlockStmt(block *Block) any {
	r.beginScope()
	r.Resolve(&block.statements)
	r.endScope()
	return nil
}

/*
*
We split binding into two steps, declaring then defining, in order to handle funny edge cases like this:
var a = "outer";

	{
	  var a = a;
	}

*
*/
func (r *Resolver) VisitVarStmt(varStmt *Var) any {
	r.declare(varStmt.name)
	if varStmt.initializer != nil {
		r.resolveExpr(varStmt.initializer)
	}
	r.define(varStmt.name)
	return nil
}

func (r *Resolver) VisitVariableExpr(expr *Variable) any {
	if !r.scopes.IsEmpty() {
		scope := r.scopes.Peek()
		if defined, ok := (*scope)[expr.name.Lexeme]; ok && !defined {
			r.error_reporter(expr.name.Line, expr.name.Lexeme, "can't read local variable in its own initializer")
		}
	}
	r.resolveLocal(expr, expr.name)
	return nil
}

func (r *Resolver) VisitAssignExpr(expr *Assign) any {
	r.resolveExpr(expr.value)
	r.resolveLocal(expr, expr.name)
	return nil
}

func (r *Resolver) VisitFunctionStmt(stmt *Function) any {
	r.declare(stmt.name)
	r.define(stmt.name)

	r.resolveFunction(stmt, FUNCTION)
	return nil
}

func (r *Resolver) VisitExpressionStmt(stmt *Expression) any {
	r.resolveExpr(stmt.expression)

	return nil
}

func (r *Resolver) VisitIfStmt(stmt *If) any {
	r.resolveExpr(stmt.condition)
	r.resolveStmt(stmt.trueBranch)
	if stmt.falseBranch != nil {
		r.resolveStmt(stmt.falseBranch)
	}
	return nil
}

func (r *Resolver) VisitPrintStmt(stmt *Print) any {
	r.resolveExpr(stmt.expression)

	return nil
}

func (r *Resolver) VisitWhileStmt(stmt *While) any {
	r.resolveStmt(stmt.condition)
	r.resolveStmt(stmt.body)

	return nil
}

func (r *Resolver) VisitBreakStatement(stmt *Break) any {
	return nil
}

func (r *Resolver) VisitReturnStmt(stmt *Return) any {
	if r.visitedFunctions.Len() == 0 {
		r.error_reporter(stmt.keyword.Line, "", "return outside function body")
	}
	if stmt.value != nil {
		if *r.visitedFunctions.Peek() == CONSTRUCTOR {
			r.error_reporter(stmt.keyword.Line, "", "can't return a value from an initializer")
		}
		r.resolveExpr(stmt.value)
	}

	return nil
}

func (r *Resolver) VisitBinaryExpr(expr *Binary) any {
	r.resolveExpr(expr.left)
	r.resolveExpr(expr.right)

	return nil
}

func (r *Resolver) VisitGroupingExpr(expr *Grouping) any {
	r.resolveExpr(expr.expression)

	return nil
}

func (r *Resolver) VisitLiteralExpr(expr *Literal) any {
	return nil
}

func (r *Resolver) VisitUnaryExpr(expr *Unary) any {
	r.resolveExpr(expr.right)

	return nil
}

func (r *Resolver) VisitTernaryExpr(expr *Ternary) any {
	r.resolveExpr(expr.condition)
	r.resolveExpr(expr.trueBranch)
	r.resolveExpr(expr.falseBranch)

	return nil
}

func (r *Resolver) VisitLogicalExpr(expr *Logical) any {
	r.resolveExpr(expr.left)
	r.resolveExpr(expr.right)

	return nil
}

func (r *Resolver) VisitCallExpr(expr *Call) any {
	r.resolveExpr(expr.callee)

	for _, arg := range expr.arguments {
		r.resolveExpr(arg)
	}

	return nil
}

func (r *Resolver) VisitGetExpr(expr *Get) any {
	r.resolveExpr(expr.instance)
	return nil
}

func (r *Resolver) VisitSetExpr(expr *Set) any {
	r.resolveExpr(expr.object)
	r.resolveExpr(expr.value)
	return nil
}

func (r *Resolver) VisitThisExpr(expr *This) any {
	if !r.inClass {
		r.error_reporter(expr.keyword.Line, "", "can't use 'this' outside of a class")
	}
	r.resolveLocal(expr, expr.keyword)
	return nil
}

func (r *Resolver) Resolve(stmts *[]Stmt) {
	for _, stmt := range *stmts {
		r.resolveStmt(stmt)
	}
}

func (r *Resolver) resolveStmt(stmt Stmt) {
	stmt.accept(r)
}

func (r *Resolver) resolveExpr(expr Expr) {
	expr.accept(r)
}

func (r *Resolver) resolveLocal(expr Expr, name scanner.Token) {
	for i := r.scopes.Len() - 1; i >= 0; i-- {
		scope := (*r.scopes)[i]
		if _, ok := scope[name.Lexeme]; ok {
			r.interpreter.resolveLocal(expr, r.scopes.Len()-1-i)
			return
		}
	}
}

func (r *Resolver) resolveFunction(stmt *Function, functionType FunctionType) {
	r.visitedFunctions.Push(functionType)

	r.beginScope()

	for _, param := range stmt.params {
		r.declare(param)
		r.define(param)
	}
	r.Resolve(&stmt.body)

	r.endScope()

	r.visitedFunctions.Pop()
}

func (r *Resolver) beginScope() {
	r.scopes.Push(make(map[string]bool))
}

func (r *Resolver) endScope() {
	r.scopes.Pop()
}

func (r *Resolver) declare(name scanner.Token) {
	if r.scopes.IsEmpty() {
		return
	}

	scope := r.scopes.Peek()

	(*scope)[name.Lexeme] = false
}

func (r *Resolver) define(name scanner.Token) {
	if r.scopes.IsEmpty() {
		return
	}

	scope := r.scopes.Peek()

	(*scope)[name.Lexeme] = true
}
