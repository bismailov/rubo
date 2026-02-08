package parser

import (
	"fmt"
	"strconv"
)

const (
	_ int = iota
	LOWEST
	EQUALS      // ==
	LESSGREATER // > or <
	SUM         // +
	PRODUCT     // *
	PREFIX      // -X or !X
	CALL        // myFunction(X)
)

var precedences = map[TokenType]int{
	EQ:       EQUALS,
	NOT_EQ:   EQUALS,
	LT:       LESSGREATER,
	GT:       LESSGREATER,
	PLUS:     SUM,
	MINUS:    SUM,
	SLASH:    PRODUCT,
	ASTERISK: PRODUCT,
	LPAREN:   CALL,
}

type (
	prefixParseFn func() Expression
	infixParseFn  func(Expression) Expression
)

type Parser struct {
	l      *Lexer
	errors []string

	curToken  Token
	peekToken Token

	prefixParseFns map[TokenType]prefixParseFn
	infixParseFns  map[TokenType]infixParseFn
}

func NewParser(l *Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	// Read two tokens, so curToken and peekToken are both set
	p.nextToken()
	p.nextToken()

	p.prefixParseFns = make(map[TokenType]prefixParseFn)
	p.registerPrefix(IDENT, p.parseIdentifier)
	p.registerPrefix(INT, p.parseIntegerLiteral)
	p.registerPrefix(DEF, p.parseFunctionLiteral)
	p.registerPrefix(MINUS, p.parsePrefixExpression)
	p.registerPrefix(LPAREN, p.parseGroupedExpression)
	p.registerPrefix(STRING, p.parseStringLiteral)

	p.infixParseFns = make(map[TokenType]infixParseFn)
	p.registerInfix(PLUS, p.parseInfixExpression)
	p.registerInfix(MINUS, p.parseInfixExpression)
	p.registerInfix(SLASH, p.parseInfixExpression)
	p.registerInfix(ASTERISK, p.parseInfixExpression)
	p.registerInfix(EQ, p.parseInfixExpression)
	p.registerInfix(NOT_EQ, p.parseInfixExpression)
	p.registerInfix(LT, p.parseInfixExpression)
	p.registerInfix(GT, p.parseInfixExpression)
	p.registerInfix(LPAREN, p.parseCallExpression)

	return p
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) ParseProgram() *Program {
	program := &Program{}
	program.Statements = []Statement{}

	for p.curToken.Type != EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}

	return program
}

func (p *Parser) parseStatement() Statement {
	// Skip extra newlines at start of statement
	for p.curToken.Type == NEWLINE {
		p.nextToken()
	}

	if p.curToken.Type == EOF {
		return nil
	}

	// For now, we only have expression statements and defs (which are expressions in Rubo)
	return p.parseExpressionStatement()
}

func (p *Parser) parseExpressionStatement() *ExpressionStatement {
	stmt := &ExpressionStatement{Token: p.curToken}

	stmt.Expression = p.parseExpression(LOWEST)

	// In Ruby, newlines or EOF terminate statements
	if p.peekToken.Type == NEWLINE {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseExpression(precedence int) Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}
	leftExp := prefix()

	for p.peekToken.Type != NEWLINE && p.peekToken.Type != EOF && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}

		p.nextToken()

		leftExp = infix(leftExp)
	}

	return leftExp
}

func (p *Parser) parseIdentifier() Expression {
	return &Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseIntegerLiteral() Expression {
	lit := &IntegerLiteral{Token: p.curToken}

	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	lit.Value = value

	return lit
}

func (p *Parser) parseFunctionLiteral() Expression {
	lit := &FunctionLiteral{Token: p.curToken}

	if !p.expectPeek(IDENT) {
		return nil
	}

	lit.Name = &Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectPeek(LPAREN) {
		return nil
	}

	lit.Parameters = p.parseFunctionParameters()

	// Ruby's def usually followed by body, then end
	// Minimal whitespace handling: skip newlines
	for p.peekToken.Type == NEWLINE {
		p.nextToken()
	}

	lit.Body = p.parseBlockStatement()

	return lit
}

func (p *Parser) parseFunctionParameters() []*Identifier {
	identifiers := []*Identifier{}

	if p.peekToken.Type == RPAREN {
		p.nextToken()
		return identifiers
	}

	p.nextToken()

	ident := &Identifier{Token: p.curToken, Value: p.curToken.Literal}
	identifiers = append(identifiers, ident)

	for p.peekToken.Type == COMMA {
		p.nextToken()
		p.nextToken()
		ident := &Identifier{Token: p.curToken, Value: p.curToken.Literal}
		identifiers = append(identifiers, ident)
	}

	if !p.expectPeek(RPAREN) {
		return nil
	}

	return identifiers
}

func (p *Parser) parseBlockStatement() *BlockStatement {
	block := &BlockStatement{Token: p.curToken}
	block.Statements = []Statement{}

	p.nextToken()

	for p.curToken.Type != END && p.curToken.Type != EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}

	return block
}

func (p *Parser) parseInfixExpression(left Expression) Expression {
	expression := &InfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}

	precedence := p.curPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)

	return expression
}

func (p *Parser) curTokenIs(t TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) peekError(t TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead",
		t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) registerPrefix(tokenType TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) noPrefixParseFnError(t TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.errors = append(p.errors, msg)
}
func (p *Parser) parseCallExpression(function Expression) Expression {
	exp := &CallExpression{Token: p.curToken, Function: function}
	exp.Arguments = p.parseExpressionList(RPAREN)
	return exp
}

func (p *Parser) parseExpressionList(end TokenType) []Expression {
	list := []Expression{}

	if p.peekTokenIs(end) {
		p.nextToken()
		return list
	}

	p.nextToken()
	list = append(list, p.parseExpression(LOWEST))

	for p.peekTokenIs(COMMA) {
		p.nextToken()
		p.nextToken()
		list = append(list, p.parseExpression(LOWEST))
	}

	if !p.expectPeek(end) {
		return nil
	}

	return list
}
func (p *Parser) parsePrefixExpression() Expression {
	expression := &PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}

	p.nextToken()

	expression.Right = p.parseExpression(PREFIX)

	return expression
}

func (p *Parser) parseGroupedExpression() Expression {
	p.nextToken()

	exp := p.parseExpression(LOWEST)

	if !p.expectPeek(RPAREN) {
		return nil
	}

	return exp
}

func (p *Parser) parseStringLiteral() Expression {
	return &StringLiteral{Token: p.curToken, Value: p.curToken.Literal}
}
