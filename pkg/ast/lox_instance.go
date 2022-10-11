package ast

import (
	"fmt"

	"github.com/fadyZohdy/gLox/pkg/scanner"
)

type LoxInstance struct {
	class  *LoxClass
	fields map[string]any
}

func (instance *LoxInstance) String() string {
	return fmt.Sprintf("%s instance", instance.class.name)
}

func NewLoxInstance(class *LoxClass) *LoxInstance {
	return &LoxInstance{class: class, fields: make(map[string]any)}
}

func (instance *LoxInstance) get(name scanner.Token) any {
	if v, ok := instance.fields[name.Lexeme]; ok {
		return v
	}
	method := instance.class.findMethod(name.Lexeme)
	if method != nil {
		return method.bind(instance)
	}
	panic(&RuntimeError{fmt.Sprintf("undefined property %s", name.Lexeme), name})
}

func (instance *LoxInstance) set(name scanner.Token, value any) {
	instance.fields[name.Lexeme] = value
}
