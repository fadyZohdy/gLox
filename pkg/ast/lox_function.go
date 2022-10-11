package ast

type LoxFunction struct {
	declaration *Function
	closure     *Environment
}

func (f *LoxFunction) arity() int {
	return len(f.declaration.params)
}

func (f *LoxFunction) isConstructor() bool {
	return f.declaration.name.Lexeme == "init"
}

func (f *LoxFunction) call(interpreter *Interpreter, arguments []any) any {
	env := NewEnvironment(f.closure)
	for i, param := range f.declaration.params {
		env.define(param.Lexeme, arguments[i])
	}
	result := interpreter.executeBlock(f.declaration.body, env)

	// special handling for calling constructor(init) on a class innstance
	if f.isConstructor() {
		return f.closure.getAt(0, "this")
	}

	return result
}

func (l LoxFunction) String() string {
	if l.declaration.isAnon() {
		return "<fn>"
	}
	return "<fn " + l.declaration.name.Lexeme + ">"
}

func (f LoxFunction) bind(instance *LoxInstance) *LoxFunction {
	env := NewEnvironment(f.closure)
	env.define("this", instance)
	return &LoxFunction{f.declaration, env}
}
