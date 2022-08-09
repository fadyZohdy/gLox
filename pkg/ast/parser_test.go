package ast

import (
	"testing"

	"github.com/fadyZohdy/gLox/pkg/scanner"
)

func TestParser(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"-123 * (45.67)", "(* (- 123) (group 45.67))"},
		{"1 + 2 * 3 / 4 + 15.1 - \"test\"", "(- (+ (+ 1 (/ (* 2 3) 4)) 15.1) test)"},
		{"1 == 2 > 3 / 4 + 5 % 6", "(== 1 (> 2 (+ (/ 3 4) (% 5 6))))"},
		{"!true == false", "(== (! true) false)"},
		{"4 > 3 ? 1 : 2", "(? (> 4 3) 1 2)"},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			s := scanner.NewScanner(tt.input, func(int, string) {})
			tokens := s.ScanTokens()
			parser := NewParser(tokens, func(int, string, string) {})
			exprs, _ := parser.Parse()
			printer := &AstPrinter{}
			printer.Print(exprs[0])
			if printer.Repr != tt.expected {
				t.Errorf("AstPrinter.Print(%v). got = %s, want %s", tt.input, printer.Repr, tt.expected)
			}
		})
	}

}
