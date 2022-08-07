package scanner

import "fmt"

type Token struct {
	Type    TokenType
	Lexeme  string
	Literal interface{}
	Line    int
}

func (t Token) String() string {
	return fmt.Sprintf(string(t.Type), t.Lexeme, t.Literal, t.Line)
}

func (t Token) Equal(other Token) bool {
	return t.Type == other.Type && t.Lexeme == other.Lexeme && t.Literal == other.Literal && t.Line == other.Line
}
