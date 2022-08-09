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
				scanner.Token{scanner.MINUS, "-", nil, 1},
				&Literal{scanner.LoxFloat64(123)}},
			scanner.Token{scanner.STAR, "*", nil, 1},
			&Grouping{
				&Literal{scanner.LoxFloat64(45.67)}}}, "(* (- 123) (group 45.67))"},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			p := &AstPrinter{}
			p.Print(tt.input)
			if p.Repr != tt.expected {
				t.Errorf("AstPrinter.Print(%v). got = %s, want %s", tt.input, p.Repr, tt.expected)
			}
		})
	}

}
