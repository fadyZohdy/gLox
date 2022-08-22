package ast

import (
	"log"
	"testing"

	"github.com/fadyZohdy/gLox/pkg/scanner"
)

func TestParser(t *testing.T) {
	tests := []struct {
		input    string
		expected []string
	}{
		{"-123 * (45.67);", []string{"(* (- 123) (group 45.67))"}},
		{"1 + 2 * 3 / 4 + 15.1 - \"test\";", []string{"(- (+ (+ 1 (/ (* 2 3) 4)) 15.1) test)"}},
		{"1 == 2 > 3 / 4 + 5 % 6;", []string{"(== 1 (> 2 (+ (/ 3 4) (% 5 6))))"}},
		{"!true == false;", []string{"(== (! true) false)"}},
		{"4 > 3 ? 1 : 2;", []string{"(? (> 4 3) 1 2)"}},
		{"true ? 1 : false ? 2 : 3, 5 + 2, 6 - 3;", []string{"(, (, (? true 1 (? false 2 3)) (+ 5 2)) (- 6 3))"}},
		{"== 2, 3 + 2;", []string{"(, nil (+ 3 2))"}},
		{">= 2, +3, *4, false ? /2 : 5;", []string{"(, (, (, nil nil) nil) (? false nil 5))"}},
		{"print 3 + 5;", []string{"(print (+ 3 5))"}},
		{"var x = 3 + 5;", []string{"(var x (+ 3 5))"}},
		{"var y;", []string{"(var y nil)"}},
		{"var x = ;", []string{}},
		{"var x; x = 3;", []string{"(var x nil)", "(x =  3)"}},
		{"1 == 2 and 3 < 2 or 4 > 3;", []string{"(or (and (== 1 2) (< 3 2)) (> 4 3))"}},
		//while
		{"var x = 0; while (x > 5) { print x; x++; }", []string{"(var x 0)", "while ((> x 5)) {{(print x) (x =  (++ x)) }}"}},
		//for
		{"for (var x = 1; x < 5; x ++) print x;", []string{"{(var x 1) while ((< x 5)) {{(print x) (x =  (++ x)) }} }"}},
		{"var x; for (x = 1; x < 5; x ++) print x;", []string{"(var x nil)", "{(x =  1) while ((< x 5)) {{(print x) (x =  (++ x)) }} }"}},
		{"var x = 1; for (; x < 5; x ++) print x;", []string{"(var x 1)", "while ((< x 5)) {{(print x) (x =  (++ x)) }}"}},
		{"for (var x = 1;; x ++) print x;", []string{"{(var x 1) while (true) {{(print x) (x =  (++ x)) }} }"}},
		{"for (var x = 1;;) print x;", []string{"{(var x 1) while (true) {(print x)} }"}},
		{"for (;;) print x;", []string{"while (true) {(print x)}"}},
		{"for (;;) break;", []string{"while (true) {break}"}},
		{"break;", []string{""}},
		{"call(x, y);", []string{"(call x y)"}},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			s := scanner.NewScanner(tt.input, func(int, string) {})
			tokens := s.ScanTokens()
			parser := NewParser(tokens, func(l int, w, m string) { log.Println("[line ", l, "] Error", w, ": ", m) })
			stmts, _ := parser.Parse()
			for i, stmt := range stmts {
				printer := &AstPrinter{}
				res := printer.Print(stmt)
				if res != tt.expected[i] {
					t.Errorf("AstPrinter.Print(%v). got = %s, want %s", tt.input, res, tt.expected[i])
				}
			}
		})
	}

}
