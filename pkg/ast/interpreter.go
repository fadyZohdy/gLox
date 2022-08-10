package ast

import (
	"fmt"

	"github.com/fadyZohdy/gLox/pkg/scanner"
)

type Interpreter struct {
	error_reporter func(error RuntimeError)
}

func NewInterpreter(error_reporter func(error RuntimeError)) *Interpreter {
	return &Interpreter{error_reporter}
}

func (i *Interpreter) Interpret(expr Expr) (value any, err error) {
	defer func() {
		e := recover()
		if e == nil {
			// no panic error
			return
		}
		if e, ok := e.(RuntimeError); ok {
			// err panic occurred
			err = e
			i.error_reporter(e)
		}
	}()

	value = i.evaluate(expr)
	return
}

func (i *Interpreter) VisitLiteralExpr(expr *Literal) any {
	return expr.value
}

func (i *Interpreter) VisitGroupingExpr(expr *Grouping) any {
	return i.evaluate(expr.expression)
}

func (i *Interpreter) VisitUnaryExpr(expr *Unary) any {
	right := i.evaluate(expr.right)
	switch expr.operator.Type {
	case scanner.MINUS:
		return -1 * checkNumber(right, expr.operator)
	case scanner.BANG:
		return !isTruthy(right)
	}
	return nil
}

func (i *Interpreter) VisitBinaryExpr(expr *Binary) any {
	left := i.evaluate(expr.left)
	right := i.evaluate(expr.right)

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
			panic(RuntimeError{Message: "Division by zero.", Token: expr.operator})
		}
		return checkNumber(left, expr.operator) / right
	case scanner.PLUS:
		if isNumber(left) && isNumber(right) {
			return checkNumber(left, expr.operator) + checkNumber(right, expr.operator)
		} else if isString(left) && isString(right) {
			return checkString(left) + checkString(right)
		} else {
			if isNumber(left) && isString(right) {
				return fmt.Sprintf("%v", checkNumber(left, expr.operator)) + checkString(right)
			}
			if isString(left) && isNumber(right) {
				return checkString(left) + fmt.Sprintf("%v", checkNumber(right, expr.operator))
			}
			panic(RuntimeError{Message: "operands can be numbers or strings", Token: expr.operator})
		}
	case scanner.COMMA:
		return right
	}
	panic(RuntimeError{Message: "Unknown operator.", Token: expr.operator})
}

func (i *Interpreter) VisitTernaryExpr(expr *Ternary) any {
	condition := i.evaluate(expr.condition)
	if flag, ok := condition.(bool); ok {
		if flag {
			return i.evaluate(expr.trueBranch)
		} else {
			return i.evaluate(expr.falseBranch)
		}
	} else {
		panic(RuntimeError{Message: "Ternary condition value is not a boolean."})
	}
}

func checkNumber(value any, token scanner.Token) float64 {
	if value_f, ok := value.(float64); ok {
		return value_f
	} else {
		panic(RuntimeError{Message: "Operand must be a number.", Token: scanner.Token{}})
	}
}

func isNumber(value any) bool {
	_, ok := value.(float64)
	return ok
}

func checkString(value any) string {
	if value_s, ok := value.(string); ok {
		return value_s
	} else {
		panic(RuntimeError{Message: "Operand must be a string.", Token: scanner.Token{}})
	}
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

func (i *Interpreter) evaluate(expr Expr) any {
	if expr == nil {
		return nil
	}
	return expr.accept(i)
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
