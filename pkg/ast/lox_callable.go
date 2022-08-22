package ast

type LoxCallable interface {
	arity() int
	call(interpreter *Interpreter, arguments []any) any
}
