package ast

import (
	"log"
	"testing"

	"github.com/fadyZohdy/gLox/pkg/scanner"
)

func TestInterpreter(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected any
		err      error
	}{
		{"test addition", "3 + 4", float64(7), nil},
		{"test string concatenation", "\"test\" + \"test\"", "testtest", nil},
		{"test number + string concatenation", "3 + \"test\"", "3test", nil},
		{"test string + number concatenation", "\"test\" + 3", "test3", nil},
		{"test precedence", "5 - 4 / 2", float64(3), nil},
		{"test uniary minus", "-1 * 3", float64(-3), nil},
		{"test comma operator", "4 * 5, 3 + 2", float64(5), nil},
		{"test comma operator with erronous unary operator", "+2, 3 + 2", float64(5), nil},
		{"test comma with ternary", "false ? 1 : 4 <= 2 ? 2 : 3, 5 + 2, 6 - 3", float64(3), nil},
		{"test ternary operator", "3 >= 2 ? 1 : 2", float64(1), nil},
		{"test ternary operator string comparison", "\"hello\" != \"hello\" ? 1 : 2", float64(2), nil},
		{"compare string with number should give false", "\"hello\" == 3", false, nil},
		{"compare nil with nil", "nil == nil", true, nil},
		{"compare value with nil", "1 == nil", false, nil},
		{"test negation + compare nil with nil", "!(nil == nil)", false, nil},
		{"test division by zero", "1 / 0", nil,
			RuntimeError{Message: "Division by zero.", Token: scanner.Token{Type: scanner.SLASH, Literal: "/", Line: 1}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := scanner.NewScanner(tt.input, func(int, string) {})
			tokens := s.ScanTokens()
			parser := NewParser(tokens, func(l int, w, m string) { log.Println("[line ", l, "] Error", w, ": ", m) })
			exprs, _ := parser.Parse()
			interpreter := NewInterpreter(func(err RuntimeError) { log.Println(err.Message, "[line ", err.Token.Line, "]") })
			res, err := interpreter.Interpret(exprs[0])
			if res != tt.expected {
				t.Errorf("Interpreter.interpret(%v). got = %s, want %s", tt.input, res, tt.expected)
			}
			if _, ok := err.(RuntimeError); !ok && tt.err != nil {
				t.Errorf("Interpreter.interpret(%v). got = %v, want %v", tt.input, err, tt.err)
			}
		})
	}

}
