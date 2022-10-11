package ast

type LoxClass struct {
	name    string
	methods map[string]*LoxFunction
}

func (klass *LoxClass) String() string {
	return klass.name
}

func (klass *LoxClass) arity() int {
	init := klass.findMethod("init")
	if init != nil {
		return init.arity()
	}
	return 0
}

func (klass *LoxClass) call(interpreter *Interpreter, arguments []any) any {
	instance := NewLoxInstance(klass)
	init := klass.findMethod("init")
	if init != nil {
		init.bind(instance).call(interpreter, arguments)
	}
	return instance
}

func (klass *LoxClass) findMethod(name string) *LoxFunction {
	if m, ok := klass.methods[name]; ok {
		return m
	}
	return nil
}
