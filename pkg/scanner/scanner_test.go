package scanner

import (
	"testing"
)

func TestMatch(t *testing.T) {
	tests := []struct {
		input    string
		current  int
		match    rune
		expected bool
	}{
		{"", 0, 'a', false},
		{"a", 0, 'a', true},
		{"==", 1, '=', true},
		{"==", 2, '=', false},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			scanner := NewScanner(tt.input, func(i int, s string) {})
			scanner.current = tt.current
			got := scanner.match(tt.match)
			if got != tt.expected {
				t.Errorf(
					"match(%s)",
					string(tt.match),
				)
			}
		})
	}

}

func compare(input string, got, expected []Token, t *testing.T) {
	for i, token := range got {
		if !token.Equal(expected[i]) {
			t.Errorf(
				"ScanTokens(%v). got[%d] = %s, want %s",
				input,
				i,
				token,
				expected[i],
			)
		}
	}

}
func TestScanTokens(t *testing.T) {
	tests := []struct {
		input    string
		expected []Token
	}{
		{"", []Token{{EOF, "", nil, 1}}},
		{" \n {}", []Token{
			{LEFT_BRACE, "{", nil, 2},
			{RIGHT_BRACE, "}", nil, 2},
			{EOF, "", nil, 2}},
		},
		{"// This is a comment\n\n", []Token{{EOF, "", nil, 3}}},
		{"\"hello\" // \"dkndknjd ndjjnkj\" \n (\"world\") ", []Token{
			{STRING, "\"hello\"", "hello", 1},
			{LEFT_PAREN, "(", nil, 2},
			{STRING, "\"world\"", "world", 2},
			{RIGHT_PAREN, ")", nil, 2},
			{EOF, "", nil, 2}},
		},
		{"123.22.", []Token{{Number, "123.22", float64(123.22), 1}, {DOT, ".", nil, 1}, {EOF, "", nil, 1}}},
		{"orange = 0\n while (orange <= 3) {\n if (orange % 2 == 0) {\n print(orange)\n orange = orange  + 1 \n} ", []Token{
			{IDENTIFIER, "orange", nil, 1},
			{EQUAL, "=", nil, 1},
			{Number, "0", float64(0), 1},
			{WHILE, "while", nil, 2},
			{LEFT_PAREN, "(", nil, 2},
			{IDENTIFIER, "orange", nil, 2},
			{LESS_EQUAL, "<=", nil, 2},
			{Number, "3", float64(3), 2},
			{RIGHT_PAREN, ")", nil, 2},
			{LEFT_BRACE, "{", nil, 2},
			{IF, "if", nil, 3},
			{LEFT_PAREN, "(", nil, 3},
			{IDENTIFIER, "orange", nil, 3},
			{MODULO, "%", nil, 3},
			{Number, "2", float64(2), 3},
			{EQUAL_EQUAL, "==", nil, 3},
			{Number, "0", float64(0), 3},
			{RIGHT_PAREN, ")", nil, 3},
			{LEFT_BRACE, "{", nil, 3},
			{PRINT, "print", nil, 4},
			{LEFT_PAREN, "(", nil, 4},
			{IDENTIFIER, "orange", nil, 4},
			{RIGHT_PAREN, ")", nil, 4},
			{IDENTIFIER, "orange", nil, 5},
			{EQUAL, "=", nil, 5},
			{IDENTIFIER, "orange", nil, 5},
			{PLUS, "+", nil, 5},
			{Number, "1", float64(1), 5},
			{RIGHT_BRACE, "}", nil, 6},
			{EOF, "", nil, 6},
		}},
		{"/*dkmdknde\ndjwndjkwn\nkwjhdkjw*/", []Token{{EOF, "", nil, 3}}},

		//malformed multiline comment
		{"/*dkmdknde\ndjwndjkwn\nkwjhdkjw", []Token{{EOF, "", nil, 3}}},
		//malformed multiline comment
		{"/*dkmdknde\ndjwndjkwn\nkwjhdkjw*", []Token{{EOF, "", nil, 3}}},
		// unknown character
		{"3 ^ 4", []Token{{Number, "3", float64(3), 1}, {Number, "4", float64(4), 1}, {EOF, "", nil, 1}}},
		//unterminated string
		{"\"hello 4 * 2", []Token{{EOF, "", nil, 1}}},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			scanner := NewScanner(tt.input, func(i int, s string) {})
			got := scanner.ScanTokens()
			compare(tt.input, got, tt.expected, t)
		})
	}
}
