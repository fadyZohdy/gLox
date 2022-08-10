package ast

import "github.com/fadyZohdy/gLox/pkg/scanner"

type RuntimeError struct {
	Message string
	Token   scanner.Token
}

func (err RuntimeError) Error() string {
	return err.Message
}

var OnlyStringOrNumberError = &RuntimeError{Message: "operands can be numbers or strings"}

var DivisionByZeroError = &RuntimeError{Message: "division by zero"}

var UnknownOperatorError = &RuntimeError{Message: "unknown operator"}

var NotNumberError = &RuntimeError{Message: "operand is not a number"}

var NotStringError = &RuntimeError{Message: "operand is not a string"}
