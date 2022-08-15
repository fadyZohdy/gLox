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
		expected map[string]any
	}{
		{"block scopes",
			`
		var a = "global a";
		var b = "global b";
		var c = "global c";
		{
		  var a = "outer a";
		  var b = "outer b";
		  {
			var a = "inner a";
			b = "inner b";
		  }
		}
		`, map[string]any{"a": "global a", "b": "global b", "c": "global c"}},
		{"for with break", `
		for (var i = 0; i <= 15; i++) {
			for (var j = 0; i <= 15; j ++) {
				if(j >= 3) {
					break;
					// unreachable code
					j++
				} 
			}
		}
		`, map[string]any{"i": float64(15), "j": float64(3)}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := scanner.NewScanner(tt.input, func(int, string) {})
			tokens := s.ScanTokens()
			parser := NewParser(tokens, func(l int, w, m string) { log.Println("[line ", l, "] Error", w, ": ", m) })
			stmts, _ := parser.Parse()
			interpreter := NewInterpreter(func(err *RuntimeError) { log.Println(err.Message, "[line ", err.Token.Line, "]") })
			interpreter.Interpret(stmts)
			environment := interpreter.env.values
			for k, v := range environment {
				if v != tt.expected[k] {
					t.Errorf("Interpreter.interpret(%v). got = %s, want %s", tt.input, environment, tt.expected)
				}
			}
		})
	}

}
