package ast

import (
	"fmt"

	"github.com/fadyZohdy/gLox/pkg/scanner"
)

type Parser struct {
	tokens         []scanner.Token
	current        int
	error_reporter func(int, string, string)
	// used to track how many loops we are parsing currently to report error if user is breaking outside loop
	loops int
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

	if p.match(scanner.FUN) {
		return p.function("function", false)
	}

	if p.match(scanner.VAR) {
		return p.varDeclaration()
	}
	return p.statement()
}

func (p *Parser) function(kind string, argument bool) Stmt {
	var name scanner.Token
	// only anonymous functions allowed to not have function name.
	// anonymous functions can be assigned to a variable or passed directly as function argument
	if p.check(scanner.IDENTIFIER) {
		name = p.consume(scanner.IDENTIFIER, fmt.Sprintf("expect %s name", kind))
	} else if !argument {
		p.error_reporter(p.peek().Line, "", fmt.Sprintf("expect %s name", kind))
	}

	p.consume(scanner.LEFT_PAREN, fmt.Sprintf("expect '(' after %s name", kind))

	var params []scanner.Token

	if !p.check(scanner.RIGHT_PAREN) {
		if len(params) > 255 {
			p.error_reporter(p.peek().Line, "", "can't have more than 255 parameters")
		}
		params = append(params, p.consume(scanner.IDENTIFIER, "expect parameter name"))
		for p.match(scanner.COMMA) {
			params = append(params, p.consume(scanner.IDENTIFIER, "expect parameter name"))
		}
	}
	p.consume(scanner.RIGHT_PAREN, fmt.Sprintf("expect ')' after %s parameters", kind))

	p.consume(scanner.LEFT_BRACE, fmt.Sprintf("expect '{' before %s body", kind))
	body := p.block()

	return &Function{name: name, params: params, body: body}
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
	if p.match(scanner.BREAK) {
		return p.breakStatement()
	}
	if p.match(scanner.FOR) {
		return p.forStatement()
	}
	if p.match(scanner.IF) {
		return p.ifStatement()
	}
	if p.match(scanner.PRINT) {
		return p.printStatement()
	}
	if p.match(scanner.RETURN) {
		return p.returnStatement()
	}
	if p.match(scanner.WHILE) {
		return p.whileStatement()
	}

	if p.match(scanner.LEFT_BRACE) {
		return &Block{p.block()}
	}

	return p.expressionStatement()
}

func (p *Parser) breakStatement() Stmt {
	// user is trying to break outside a loop
	if p.loops == 0 {
		p.error_reporter(p.peek().Line, "", "'break' outside loop")
		return nil
	}
	p.consume(scanner.SEMICOLON, "expect ';' after break")
	return &Break{}
}

func (p *Parser) forStatement() Stmt {
	p.consume(scanner.LEFT_PAREN, "expect '(' after for")
	var initializer Stmt
	if p.match(scanner.SEMICOLON) {
		initializer = nil
	} else if p.match(scanner.VAR) {
		initializer = p.varDeclaration()
	} else {
		initializer = p.expressionStatement()
	}

	var condition Stmt
	if !p.check(scanner.SEMICOLON) {
		condition = p.expression()
	}
	p.consume(scanner.SEMICOLON, "expect ';' after loop condition")

	var increment Expr
	if !p.check(scanner.RIGHT_PAREN) {
		increment = p.expression()
	}
	p.consume(scanner.RIGHT_PAREN, "expect ')' after for clause")

	p.loops += 1
	defer func() { p.loops -= 1 }()
	body := p.statement()

	if increment != nil {
		body = &Block{[]Stmt{body, increment}}
	}

	if condition == nil {
		condition = &Literal{true}
	}
	body = &While{condition: condition, body: body}

	if initializer != nil {
		body = &Block{[]Stmt{initializer, body}}
	}

	return body
}

func (p *Parser) ifStatement() Stmt {
	p.consume(scanner.LEFT_PAREN, "expect '(' after if")
	condition := p.expression()
	p.consume(scanner.RIGHT_PAREN, "expect ')' after if condition")

	trueBranch := p.statement()
	var falseBranch Stmt = nil
	if p.match(scanner.ELSE) {
		falseBranch = p.statement()
	}

	return &If{condition: condition, trueBranch: trueBranch, falseBranch: falseBranch}
}

func (p *Parser) printStatement() Stmt {
	expr := p.expression()
	p.consume(scanner.SEMICOLON, "expect ';' after expression")
	return &Print{expression: expr}
}

func (p *Parser) returnStatement() Stmt {
	keyword := p.previous()
	var expr Expr
	if !p.check(scanner.SEMICOLON) {
		expr = p.expression()
	}
	p.consume(scanner.SEMICOLON, "expect ';' after expression")
	return &Return{keyword: keyword, value: expr}
}

func (p *Parser) whileStatement() Stmt {
	p.consume(scanner.LEFT_PAREN, "expect '(' after while")
	condition := p.expression()
	p.consume(scanner.RIGHT_PAREN, "expect ')' after while condition")

	p.loops += 1
	defer func() { p.loops -= 1 }()
	body := p.statement()

	return &While{condition: condition, body: body}
}

func (p *Parser) block() (stmts []Stmt) {
	for !p.isAtEnd() && !p.check(scanner.RIGHT_BRACE) {
		stmts = append(stmts, p.declaration())
	}
	p.consume(scanner.RIGHT_BRACE, "expect '}' at end of block")
	return
}

func (p *Parser) expressionStatement() Stmt {
	expr := p.expression()
	p.consume(scanner.SEMICOLON, "expect ';' after expression")
	return &Expression{expression: expr}
}

func (p *Parser) expression(args ...bool) Expr {
	return p.assignment(args...)
}

func (p *Parser) assignment(args ...bool) Expr {
	var expr Expr
	// this is a hack to distinguish between block expressions and functions arguments
	// since both of them use comma as a separator
	if len(args) > 0 && args[0] {
		expr = p.ternary()
	} else {
		expr = p.comma()
	}
	if p.match(scanner.EQUAL) {
		equals := p.previous()
		right := p.assignment(false)
		if variable, ok := expr.(*Variable); ok {
			return &Assign{variable.name, right}
		}
		p.error_reporter(equals.Line, "", "invalid assignment target")
	}

	// TODO: investigate cleaner way for incrementing/decrementing
	if p.match(scanner.INCREMENT, scanner.DECREMENT) {
		operator := p.previous()
		if variable, ok := expr.(*Variable); ok {
			return &Assign{variable.name, &Unary{operator, &Variable{variable.name}}}
		}
		p.error_reporter(operator.Line, "", "invalid assignment target")
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
	expr := p.or()
	if p.match(scanner.QUESTION_MARK) {
		trueBranch := p.expression()
		p.consume(scanner.COLON, "Expect ':' after true branch of ternary expression")
		falseBranch := p.ternary()
		expr = &Ternary{expr, trueBranch, falseBranch}
	}
	return expr
}

func (p *Parser) or() Expr {
	expr := p.and()
	if p.match(scanner.OR) {
		operator := p.previous()
		right := p.and()
		expr = &Logical{left: expr, operator: operator, right: right}
	}
	return expr
}

func (p *Parser) and() Expr {
	expr := p.equality()
	if p.match(scanner.AND) {
		operator := p.previous()
		right := p.equality()
		expr = &Logical{left: expr, operator: operator, right: right}
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
	return p.call()
}

func (p *Parser) call() Expr {
	expr := p.primary()

	for {
		if p.match(scanner.LEFT_PAREN) {
			expr = p.finishCall(expr)
		} else {
			break
		}
	}
	return expr
}

func (p *Parser) finishCall(callee Expr) Expr {
	var arguments []Expr
	// function with no arguments
	if p.match(scanner.RIGHT_PAREN) {
		return &Call{callee: callee, arguments: arguments, paren: p.previous()}
	}

	arguments = append(arguments, p.argument())
	for p.match(scanner.COMMA) {
		if len(arguments) > 255 {
			p.error_reporter(p.peek().Line, "", "can't have more than 255 arguments")
		}
		arguments = append(arguments, p.argument())
	}
	p.consume(scanner.RIGHT_PAREN, "expect ')' after arguments")
	return &Call{callee: callee, arguments: arguments, paren: p.previous()}
}

func (p *Parser) argument() Stmt {
	if p.match(scanner.FUN) {
		return p.function("function", true)
	}
	return p.expression(true)
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
