package ast

import (
	"testing"

	"github.com/fadyZohdy/gLox/pkg/scanner"
)

func TestPrinter(t *testing.T) {
	tests := []struct {
		input    Expr
		expected string
	}{
		{&Binary{
			&Unary{
				scanner.Token{Type: scanner.MINUS, Lexeme: "-", Literal: nil, Line: 1},
				&Literal{123}},
			scanner.Token{Type: scanner.STAR, Lexeme: "*", Literal: nil, Line: 1},
			&Grouping{
				&Literal{45.67}}}, "(* (- 123) (group 45.67))"},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			p := &AstPrinter{}
			res := p.Print(tt.input)
			if res != tt.expected {
				t.Errorf("AstPrinter.Print(%v). got = %s, want %s", tt.input, res, tt.expected)
			}
		})
	}

}
