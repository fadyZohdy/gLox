package ast

import (
	"fmt"
	"math"

	"github.com/fadyZohdy/gLox/pkg/scanner"
)

type Interpreter struct {
	errorReporter    func(error *RuntimeError)
	env              *Environment
	breakEncountered bool
	returnValue      any
	locals           map[Expr]int
}

func NewInterpreter(errorReporter func(error *RuntimeError)) *Interpreter {
	i := &Interpreter{errorReporter, NewEnvironment(nil), false, nil, make(map[Expr]int)}
	i.env.define("clock", &Clock{})
	return i
}

func (i *Interpreter) Interpret(stmts *[]Stmt) (err error) {
	defer func() {
		e := recover()
		if e, ok := e.(*RuntimeError); ok {
			// err panic occurred
			err = e

			i.errorReporter(e)
		}
	}()

	for _, stmt := range *stmts {
		i.execute(stmt)
	}

	return
}

func (i *Interpreter) execute(stmt Stmt) {
	if i.breakEncountered {
		return
	}
	if i.returnValue != nil {
		return
	}
	stmt.accept(i)
}

func (i *Interpreter) Evaluate(expr Expr) any {
	if expr == nil {
		return nil
	}
	return expr.accept(i)
}

func (i *Interpreter) resolveLocal(expr Expr, depth int) {
	i.locals[expr] = depth
}

func (i *Interpreter) VisitLiteralExpr(expr *Literal) any {
	return expr.value
}

func (i *Interpreter) VisitGroupingExpr(expr *Grouping) any {
	return i.Evaluate(expr.expression)
}

func (i *Interpreter) VisitUnaryExpr(expr *Unary) any {
	// TODO: investigate cleaner way for incrementing/decrementing
	if variable, ok := expr.right.(*Variable); ok {
		value := i.env.get(variable.name)
		if value == nil {
			panic(&RuntimeError{fmt.Sprintf("%s is declared but not initialized", variable.name.Lexeme), variable.name})
		}
		if i, ok := value.(float64); ok {
			if expr.operator.Type == scanner.INCREMENT {
				return i + 1
			} else {
				return i - 1
			}
		} else {
			panic(&RuntimeError{fmt.Sprintf("%s is not a number", variable.name.Lexeme), variable.name})
		}
	}
	right := i.Evaluate(expr.right)
	switch expr.operator.Type {
	case scanner.MINUS:
		return -1 * checkNumber(right, expr.operator)
	case scanner.BANG:
		return !isTruthy(right)
	}
	return nil
}

func (i *Interpreter) VisitBinaryExpr(expr *Binary) any {
	left := i.Evaluate(expr.left)
	right := i.Evaluate(expr.right)

	switch expr.operator.Type {
	case scanner.EQUAL_EQUAL:
		return isEqual(left, right)
	case scanner.BANG_EQUAL:
		return !isEqual(left, right)
	case scanner.GREATER:
		return checkNumber(left, expr.operator) > checkNumber(right, expr.operator)
	case scanner.GREATER_EQUAL:
		return checkNumber(left, expr.operator) >= checkNumber(right, expr.operator)
	case scanner.LESS:
		return checkNumber(left, expr.operator) < checkNumber(right, expr.operator)
	case scanner.LESS_EQUAL:
		return checkNumber(left, expr.operator) <= checkNumber(right, expr.operator)
	case scanner.MINUS:
		return checkNumber(left, expr.operator) - checkNumber(right, expr.operator)
	case scanner.STAR:
		return checkNumber(left, expr.operator) * checkNumber(right, expr.operator)
	case scanner.SLASH:
		right := checkNumber(right, expr.operator)
		if right == 0 {
			panicWithToken(DivisionByZeroError, expr.operator)
		}
		return checkNumber(left, expr.operator) / right
	case scanner.MODULO:
		right := checkNumber(right, expr.operator)
		if right == 0 {
			panicWithToken(DivisionByZeroError, expr.operator)
		}
		return math.Mod(checkNumber(left, expr.operator), right)
	case scanner.PLUS:
		if isNumber(left) && isNumber(right) {
			return checkNumber(left, expr.operator) + checkNumber(right, expr.operator)
		} else if isString(left) && isString(right) {
			return checkString(left, expr.operator) + checkString(right, expr.operator)
		} else {
			if isNumber(left) && isString(right) {
				return fmt.Sprintf("%v", checkNumber(left, expr.operator)) + checkString(right, expr.operator)
			}
			if isString(left) && isNumber(right) {
				return checkString(left, expr.operator) + fmt.Sprintf("%v", checkNumber(right, expr.operator))
			}
			panicWithToken(OnlyStringOrNumberError, expr.operator)
		}
	case scanner.COMMA:
		return right
	}
	panicWithToken(UnknownOperatorError, expr.operator)
	return nil
}

func (i *Interpreter) VisitTernaryExpr(expr *Ternary) any {
	condition := i.Evaluate(expr.condition)
	if flag, ok := condition.(bool); ok {
		if flag {
			return i.Evaluate(expr.trueBranch)
		} else {
			return i.Evaluate(expr.falseBranch)
		}
	} else {
		panic(&RuntimeError{Message: "ternary condition value is not a boolean"})
	}
}

func (i *Interpreter) VisitLogicalExpr(expr *Logical) any {
	left := i.Evaluate(expr.left)
	if expr.operator.Type == scanner.OR {
		if isTruthy(left) {
			return left
		}
	} else {
		if !isTruthy(left) {
			return left
		}
	}
	return i.Evaluate(expr.right)
}

func (i *Interpreter) VisitVariableExpr(expr *Variable) any {
	x := i.lookUpVariable(expr.name, expr)
	return x
	// value := i.env.get(expr.name)
	// if value == nil {
	// 	panic(&RuntimeError{fmt.Sprintf("%s is declared but not initialized", expr.name.Lexeme), expr.name})
	// }
	// return value
}

func (i *Interpreter) lookUpVariable(name scanner.Token, expr Expr) any {
	if depth, ok := i.locals[expr]; ok {
		return i.env.getAt(depth, name.Lexeme)
	} else {
		return i.env.get(name)
	}

}

func (i *Interpreter) VisitVarStmt(stmt *Var) any {
	var value any
	if stmt.initializer != nil {
		value = i.Evaluate(stmt.initializer)
	}
	i.env.define(stmt.name.Lexeme, value)
	return nil
}

func (i *Interpreter) VisitExpressionStmt(stmt *Expression) any {
	return i.Evaluate(stmt.expression)
}

func (i *Interpreter) VisitIfStmt(stmt *If) any {
	flag := isTruthy(i.Evaluate(stmt.condition))
	if flag {
		i.execute(stmt.trueBranch)
	} else if stmt.falseBranch != nil {
		i.execute(stmt.falseBranch)
	}
	return nil
}

func (i *Interpreter) VisitPrintStmt(stmt *Print) any {
	value := i.Evaluate(stmt.expression)
	fmt.Println(value)
	return nil
}

func (i *Interpreter) VisitWhileStmt(stmt *While) any {
	for isTruthy(i.Evaluate(stmt.condition)) {
		i.execute(stmt.body)
		if i.breakEncountered {
			i.breakEncountered = false
			return nil
		}
	}
	return nil
}

func (i *Interpreter) VisitAssignExpr(expr *Assign) any {
	value := i.Evaluate(expr.value)
	if depth, ok := i.locals[expr]; ok {
		i.env.assignAt(depth, expr.name, value)
	} else {
		i.env.assign(expr.name, value)
	}
	return value
}

func (i *Interpreter) VisitCallExpr(expr *Call) any {
	callee := i.Evaluate(expr.callee)

	arguments := []any{}

	for _, arg := range expr.arguments {
		if f, ok := arg.(*Function); ok {
			arguments = append(arguments, &LoxFunction{declaration: f, closure: i.env})
		} else {
			arguments = append(arguments, i.Evaluate(arg))
		}
	}

	if function, ok := callee.(LoxCallable); ok {
		if len(arguments) != function.arity() {
			panic(&RuntimeError{fmt.Sprintf("expected %d arguments but got %d", function.arity(), len(arguments)), expr.paren})
		}
		return function.call(i, arguments)
	} else {
		panic(&RuntimeError{"can only call functions or classes", expr.paren})
	}
}

func (i *Interpreter) VisitReturnStmt(stmt *Return) any {
	var value any
	if stmt.value != nil {
		value = i.Evaluate(stmt.value)
	}
	i.returnValue = value
	return nil
}

func (i *Interpreter) VisitBlockStmt(block *Block) any {
	if len(block.statements) > 0 {
		env := NewEnvironment(i.env)
		i.executeBlock(block.statements, env)
	}
	return nil
}

func (i *Interpreter) VisitBreakStatement(stmt *Break) any {
	i.breakEncountered = true
	return nil
}

func (i *Interpreter) VisitFunctionStmt(stmt *Function) any {
	f := &LoxFunction{declaration: stmt, closure: i.env}
	i.env.define(stmt.name.Lexeme, f)
	return nil
}

func (i *Interpreter) executeBlock(stmts []Stmt, env *Environment) any {
	prevEnv := i.env

	defer func() {
		i.env = prevEnv
	}()

	i.env = env

	for _, stmt := range stmts {
		i.execute(stmt)
		if i.returnValue != nil {
			v := i.returnValue
			i.returnValue = nil
			return v
		}
	}
	return nil
}

func panicWithToken(e *RuntimeError, token scanner.Token) {
	e.Token = token
	panic(e)
}

func checkNumber(value any, token scanner.Token) (f float64) {
	if value_f, ok := value.(float64); ok {
		f = value_f
	} else {
		panicWithToken(NotNumberError, token)
	}
	return
}

func isNumber(value any) bool {
	_, ok := value.(float64)
	return ok
}

func checkString(value any, token scanner.Token) (s string) {
	if value_s, ok := value.(string); ok {
		s = value_s
	} else {
		panicWithToken(NotStringError, token)
	}
	return
}

func isString(value any) bool {
	_, ok := value.(string)
	return ok
}

func isEqual(left, right any) bool {
	if left == nil && right == nil {
		return true
	}
	if left == nil {
		return false
	}
	return left == right
}

// Lox follows Ruby’s simple rule: false and nil are falsey, and everything else is truthy.
func isTruthy(value any) bool {
	if value == nil {
		return false
	}
	if value_b, ok := value.(bool); ok {
		return value_b
	}
	return true
}
