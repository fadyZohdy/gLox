package scanner

import (
	"strconv"
	"unicode"
)

var keywords = map[string]TokenType{
	"and":    AND,
	"class":  CLASS,
	"else":   ELSE,
	"false":  FALSE,
	"for":    FOR,
	"fun":    FUN,
	"if":     IF,
	"nil":    NIL,
	"or":     OR,
	"print":  PRINT,
	"return": RETURN,
	"super":  SUPER,
	"this":   THIS,
	"true":   TRUE,
	"var":    VAR,
	"while":  WHILE,
}

type Scanner struct {
	source string
	tokens []Token
	// start of current token not file
	start, current, line int
	error_reporter       func(int, string)
}

func NewScanner(source string, error func(int, string)) *Scanner {
	return &Scanner{source: source, error_reporter: error, line: 1}
}

func (s *Scanner) ScanTokens() []Token {
	for !s.isAtEnd() {
		s.start = s.current
		s.scanToken()
	}
	s.tokens = append(s.tokens, Token{Type: EOF, Line: s.line})
	return s.tokens
}

func (s *Scanner) scanToken() {

	c := s.advance()

	switch c {
	case '(':
		s.addToken(LEFT_PAREN)
	case ')':
		s.addToken(RIGHT_PAREN)
	case '{':
		s.addToken(LEFT_BRACE)
	case '}':
		s.addToken(RIGHT_BRACE)
	case ',':
		s.addToken(COMMA)
	case '.':
		s.addToken(DOT)
	case '-':
		s.addToken(MINUS)
	case '+':
		s.addToken(PLUS)
	case ';':
		s.addToken(SEMICOLON)
	case '*':
		s.addToken(STAR)
	case '%':
		s.addToken(MODULO)
	case '!':
		if s.match('=') {
			s.addToken(BANG_EQUAL)
		} else {
			s.addToken(BANG)
		}
	case '=':
		if s.match('=') {
			s.addToken(EQUAL_EQUAL)
		} else {
			s.addToken(EQUAL)
		}
	case '<':
		if s.match('=') {
			s.addToken(LESS_EQUAL)
		} else {
			s.addToken(LESS)
		}
	case '>':
		if s.match('=') {
			s.addToken(GREATER_EQUAL)
		} else {
			s.addToken(GREATER)
		}
	case '/':
		if s.match('/') {
			for !s.isAtEnd() && s.peek() != '\n' {
				s.advance()
			}
		} else if s.match('*') {
			s.scan_multiline_comment()
		} else {
			s.addToken(SLASH)
		}
	case ' ', '\r', '\t':
		break
	case '\n':
		s.line++
	case '"':
		s.scan_string()
	default:
		if unicode.IsDigit(c) {
			s.scan_number()
		} else if unicode.IsLetter(c) {
			s.scan_identifier()
		} else {
			s.error_reporter(s.line, "Unexpected character '"+string(c)+"'")
		}
	}

}

func (s *Scanner) scan_string() {
	for !s.isAtEnd() && s.peek() != '"' {
		if s.peek() == '\n' {
			s.line++
		}
		s.advance()
	}

	// we reached the end and didn't encounter the closing quote
	if s.isAtEnd() {
		s.error_reporter(s.line, "Unterminated string.")
		return
	}

	// the closing "
	s.advance()

	value := s.source[s.start+1 : s.current-1]
	s.addTokenWithLiteral(STRING, value)
}

func (s *Scanner) scan_number() {
	for unicode.IsDigit(s.peek()) {
		s.advance()
	}

	if s.peek() == '.' && unicode.IsDigit(s.peekNext()) {
		s.advance()

		for unicode.IsDigit(s.peek()) {
			s.advance()
		}
	}

	f, err := strconv.ParseFloat(s.source[s.start:s.current], 64)
	if err != nil {
		s.error_reporter(s.line, err.Error())
	}
	s.addTokenWithLiteral(Number, f)
}

func (s *Scanner) scan_identifier() {
	for isAlphaNumeric(s.peek()) {
		s.advance()
	}

	word := s.source[s.start:s.current]

	if token_type, ok := keywords[word]; ok {
		s.addToken(token_type)
	} else {
		s.addToken(IDENTIFIER)
	}
}

func (s *Scanner) scan_multiline_comment() {
	for !s.isAtEnd() && !(s.peek() == '*' && s.peekNext() == '/') {
		if s.peek() == '\n' {
			s.line++
		}
		s.advance()
	}

	if s.isAtEnd() {
		s.error_reporter(s.line, "unterminated multiline comment")
		return
	}

	// consume the trailing '*/'
	s.advance()
	s.advance()
}

func isAlphaNumeric(c rune) bool {
	return c == '_' || unicode.IsDigit(c) || unicode.IsLetter(c)
}

func (s *Scanner) match(expected rune) bool {
	if s.isAtEnd() {
		return false
	}
	if rune(s.source[s.current]) != expected {
		return false
	}
	s.current++
	return true
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}

func (s *Scanner) advance() rune {
	c := rune(s.source[s.current])
	// fmt.Println("advance:", string(c))
	s.current++
	return c
}

func (s *Scanner) peek() rune {
	if s.isAtEnd() {
		return rune(0)
	}
	return rune(s.source[s.current])
}

func (s *Scanner) peekNext() rune {
	if s.current+1 >= len(s.source) {
		return rune(0)
	}
	return rune(s.source[s.current+1])
}

func (s *Scanner) addToken(tokenType TokenType) {
	s.tokens = append(s.tokens, Token{Type: tokenType, Lexeme: s.source[s.start:s.current], Line: s.line})
}

func (s *Scanner) addTokenWithLiteral(tokenType TokenType, literal interface{}) {
	s.tokens = append(s.tokens, Token{Type: tokenType, Lexeme: s.source[s.start:s.current], Literal: literal, Line: s.line})
}
