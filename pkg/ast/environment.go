package ast

import (
	"github.com/fadyZohdy/gLox/pkg/scanner"
)

type Environment struct {
	values    map[string]any
	enclosing *Environment
}

func NewEnvironment(enclosing *Environment) *Environment {
	return &Environment{values: make(map[string]any), enclosing: enclosing}
}

func (env *Environment) define(name string, value any) {
	env.values[name] = value
}

func (env *Environment) assign(name scanner.Token, value any) {
	_, ok := env.values[name.Lexeme]
	if ok {
		env.values[name.Lexeme] = value
		return
	}
	if env.enclosing != nil {
		env.enclosing.assign(name, value)
		return
	}
	panic(&RuntimeError{Message: "assign: undefined variable '" + name.Lexeme + "'", Token: name})
}

func (env *Environment) assignAt(depth int, name scanner.Token, value any) {
	e := env.ancestor(depth)
	e.values[name.Lexeme] = value
}

func (env *Environment) get(name scanner.Token) any {
	if value, ok := env.values[name.Lexeme]; ok {
		return value
	}
	if env.enclosing != nil {
		return env.enclosing.get(name)
	}
	panic(&RuntimeError{Message: "undefined variable '" + name.Lexeme + "'", Token: name})
}

func (env *Environment) getAt(depth int, name string) any {
	return env.ancestor(depth).values[name]
}

func (env *Environment) ancestor(depth int) *Environment {
	e := env
	for i := 0; i < depth; i++ {
		e = env.enclosing
	}
	return e
}
