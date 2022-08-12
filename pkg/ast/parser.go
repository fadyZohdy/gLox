package ast

import (
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

func (p *Parser) Parse() (stmts []Stmt, err error) {
	defer func() (err error) {
		e := recover()
		if e, ok := e.(error); ok {
			// err panic occurred
			err = e
		}
		return
	}()

	stmts = make([]Stmt, 0, len(p.tokens))

	for !p.isAtEnd() {
		stmt := p.declaration()
		stmts = append(stmts, stmt)
	}

	return
}

func (p *Parser) declaration() Stmt {
	defer func() {
		e := recover()
		if e != nil {
			if _, ok := e.(ParseError); ok {
				p.synchronize()
				return
			} else {
				panic(e)
			}
		}
	}()

	if p.match(scanner.VAR) {
		return p.varDeclaration()
	}
	return p.statement()
}

func (p *Parser) varDeclaration() Stmt {
	ident := p.consume(scanner.IDENTIFIER, "expect variable name")
	var initializer Expr
	if p.match(scanner.EQUAL) {
		initializer = p.expression()
	}
	p.consume(scanner.SEMICOLON, "expect ';' after variable declaration")
	return &Var{ident, initializer}
}

func (p *Parser) statement() Stmt {
	if p.match(scanner.PRINT) {
		return p.printStatement()
	}

	if p.match(scanner.LEFT_BRACE) {
		return &Block{p.block()}
	}

	return p.expressionStatement()
}

func (p *Parser) block() (stmts []Stmt) {
	for !p.isAtEnd() && !p.check(scanner.RIGHT_BRACE) {
		stmts = append(stmts, p.declaration())
	}
	p.consume(scanner.RIGHT_BRACE, "expect '}' at end of block")
	return
}

func (p *Parser) printStatement() Stmt {
	expr := p.expression()
	p.consume(scanner.SEMICOLON, "expect ';' after expression")
	return &Print{expression: expr}
}

func (p *Parser) expressionStatement() Stmt {
	expr := p.expression()
	p.consume(scanner.SEMICOLON, "expect ';' after expression")
	return &Expression{expression: expr}
}

func (p *Parser) expression() Expr {
	return p.assignment()
}

func (p *Parser) assignment() Expr {
	expr := p.comma()
	if p.match(scanner.EQUAL) {
		equals := p.previous()
		right := p.assignment()
		if variable, ok := expr.(*Variable); ok {
			return &Assign{variable.name, right}
		}
		p.error_reporter(equals.Line, "", "invalid assignment target")
	}
	return expr
}

func (p *Parser) comma() Expr {
	expr := p.ternary()

	for p.match(scanner.COMMA) {
		operator := p.previous()
		right := p.ternary()
		expr = &Binary{expr, operator, right}
	}

	return expr
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
	for p.match(scanner.BANG_EQUAL, scanner.EQUAL_EQUAL) {
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

	if p.match(scanner.IDENTIFIER) {
		return &Variable{p.previous()}
	}

	if p.match(scanner.LEFT_PAREN) {
		expr := p.expression()
		p.consume(scanner.RIGHT_PAREN, "Expect ') after expression.")
		return &Grouping{expr}
	}

	if p.match(scanner.BANG_EQUAL, scanner.EQUAL_EQUAL) {
		operator := p.previous()
		p.comparison()
		p.error(operator, "missing left hand operand")
		return nil
	}

	if p.match(scanner.GREATER, scanner.GREATER_EQUAL, scanner.LESS, scanner.LESS_EQUAL) {
		operator := p.previous()
		p.term()
		p.error(operator, "missing left hand operand")
		return nil
	}

	if p.match(scanner.PLUS) {
		operator := p.previous()
		p.factor()
		p.error(operator, "missing left hand operand")
		return nil
	}

	if p.match(scanner.STAR, scanner.SLASH) {
		operator := p.previous()
		p.unary()
		p.error(operator, "missing left hand operand")
		return nil
	}

	p.error(p.peek(), "Expect expression.")
	panic(ParseErrorObj)
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
	panic(ParseErrorObj)
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

func (p *Parser) synchronize() {
	p.advance()

	for !p.isAtEnd() {
		if p.previous().Type == scanner.SEMICOLON {
			return
		}

		switch p.peek().Type {
		case scanner.CLASS:
			return
		case scanner.FUN:
			return
		case scanner.VAR:
			return
		case scanner.FOR:
			return
		case scanner.IF:
			return
		case scanner.WHILE:
			return
		case scanner.PRINT:
			return
		case scanner.RETURN:
			return
		}

		p.advance()
	}
}
