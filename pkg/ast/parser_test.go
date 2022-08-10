package ast

import (
	"log"
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
		{"true ? 1 : false ? 2 : 3, 5 + 2, 6 - 3", "(, (, (? true 1 (? false 2 3)) (+ 5 2)) (- 6 3))"},
		{"== 2, 3 + 2", "(, nil (+ 3 2))"},
		{">= 2, +3, *4, false ? /2 : 5", "(, (, (, nil nil) nil) (? false nil 5))"},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			s := scanner.NewScanner(tt.input, func(int, string) {})
			tokens := s.ScanTokens()
			parser := NewParser(tokens, func(l int, w, m string) { log.Println("[line ", l, "] Error", w, ": ", m) })
			exprs, _ := parser.Parse()
			printer := &AstPrinter{}
			res := printer.Print(exprs[0])
			if res != tt.expected {
				t.Errorf("AstPrinter.Print(%v). got = %s, want %s", tt.input, res, tt.expected)
			}
		})
	}

}
