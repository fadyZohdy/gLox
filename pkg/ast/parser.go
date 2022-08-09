package ast

import (
	"errors"

	scanner "github.com/fadyZohdy/gLox/pkg/scanner"
)

type Parser struct {
	tokens         []scanner.Token
	current        int
	error_reporter func(int, string, string)
}

func NewParser(tokens []scanner.Token, error_reporter func(int, string, string)) *Parser {
	return &Parser{tokens: tokens, error_reporter: error_reporter}
}

func (p *Parser) Parse() (exprs []Expr, err error) {
	defer func() {
		e := recover()
		if e == nil {
			// no panic error
			return
		}
		if e, ok := e.(error); ok {
			// err panic occurred
			err = e
		}
	}()

	for !p.isAtEnd() {
		expr := p.expression()
		exprs = append(exprs, expr)
	}

	return
}

func (p *Parser) expression() Expr {
	return p.ternary()
}

func (p *Parser) ternary() Expr {
	expr := p.equality()
	if p.match(scanner.QUESTION_MARK) {
		trueBranch := p.expression()
		p.consume(scanner.COLON, "Expect ':' after true branch of ternary expression")
		falseBranch := p.ternary()
		expr = &Ternary{expr, trueBranch, falseBranch}
	}
	return expr
}

func (p *Parser) equality() Expr {
	expr := p.comparison()
	for p.match(scanner.EQUAL, scanner.EQUAL_EQUAL) {
		operator := p.previous()
		right := p.comparison()
		expr = &Binary{expr, operator, right}
	}
	return expr
}

func (p *Parser) comparison() Expr {
	expr := p.term()
	for p.match(scanner.GREATER, scanner.GREATER_EQUAL, scanner.LESS, scanner.LESS_EQUAL) {
		operator := p.previous()
		right := p.term()
		expr = &Binary{expr, operator, right}
	}
	return expr
}

func (p *Parser) term() Expr {
	expr := p.factor()
	for p.match(scanner.PLUS, scanner.MINUS) {
		operator := p.previous()
		right := p.factor()
		expr = &Binary{expr, operator, right}
	}
	return expr
}

func (p *Parser) factor() Expr {
	expr := p.unary()
	for p.match(scanner.STAR, scanner.SLASH, scanner.MODULO) {
		operator := p.previous()
		right := p.unary()
		expr = &Binary{expr, operator, right}
	}
	return expr
}

func (p *Parser) unary() Expr {
	if p.match(scanner.MINUS, scanner.BANG) {
		operator := p.previous()
		right := p.unary()
		return &Unary{operator, right}
	}
	return p.primary()
}

func (p *Parser) primary() Expr {
	if p.match(scanner.FALSE) {
		return &Literal{false}
	}
	if p.match(scanner.TRUE) {
		return &Literal{true}
	}
	if p.match(scanner.NIL) {
		return &Literal{nil}
	}

	if p.match(scanner.Number, scanner.STRING) {
		return &Literal{p.previous().Literal}
	}

	if p.match(scanner.LEFT_PAREN) {
		expr := p.expression()
		p.consume(scanner.RIGHT_PAREN, "Expect ') after expression.")
		return &Grouping{expr}
	}

	p.error(p.peek(), "Expect expression.")
	panic(errors.New("parse error"))
}

func (p *Parser) match(types ...scanner.TokenType) bool {
	for _, tokenType := range types {
		if p.check(tokenType) {
			p.advance()
			return true

		}
	}
	return false
}

func (p *Parser) consume(tokenType scanner.TokenType, error_msg string) scanner.Token {
	if p.check(tokenType) {
		return p.advance()
	}
	p.error(p.peek(), error_msg)
	panic(errors.New("parse error"))
}

func (p *Parser) check(tokenType scanner.TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().Type == tokenType
}

func (p *Parser) isAtEnd() bool {
	return p.peek().Type == scanner.EOF
}

func (p *Parser) advance() scanner.Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p *Parser) previous() scanner.Token {
	return p.tokens[p.current-1]
}

func (p *Parser) peek() scanner.Token {
	return p.tokens[p.current]
}

func (p *Parser) error(token scanner.Token, message string) {
	if token.Type == scanner.EOF {
		p.error_reporter(token.Line, " at end", message)
	} else {
		p.error_reporter(token.Line, " at '"+token.Lexeme+"'", message)
	}
}
