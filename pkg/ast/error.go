package ast

import "github.com/fadyZohdy/gLox/pkg/scanner"

type RuntimeError struct {
	Message string
	Token   scanner.Token
}

func (err RuntimeError) Error() string {
	return err.Message
}
