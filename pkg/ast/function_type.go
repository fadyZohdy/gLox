package ast

type FunctionType int

const (
	FUNCTION FunctionType = iota
	METHOD
	STATIC_METHOD
	CONSTRUCTOR
)
