package ast

type LoxFunction struct {
	declaration *Function
	closure     *Environment
}

func (f *LoxFunction) arity() int {
	return len(f.declaration.params)
}

func (f *LoxFunction) call(interpreter *Interpreter, arguments []any) any {
	env := NewEnvironment(f.closure)
	for i, param := range f.declaration.params {
		env.define(param.Lexeme, arguments[i])
	}
	result := interpreter.executeBlock(f.declaration.body, env)
	return result
}

func (l LoxFunction) String() string {
	return "<fn " + l.declaration.name.Lexeme + ">"
}
