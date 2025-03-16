package main

import (
	"fmt"
	"os"
	"slices"
)

// Parser represents a parser that processes tokens into abstract syntax tree nodes
type Parser struct {
	tokens  []Token // Collection of tokens to parse
	current int     // Current position in the token list
}

// parse creates a parser instance and starts parsing the token stream
func (parser *Parser) parse() ([]Statement, error) {
	statements := []Statement{}

	for !parser.isEOF() {
		statement, err := parser.declaration()
		if err != nil {
			return nil, err
		}
		statements = append(statements, statement)
	}

	return statements, nil
}

// declaration parses a declaration statement (variable declarations, etc.)
func (parser *Parser) declaration() (Statement, error) {
	if parser.match(VAR) {
		return parser.varDeclaration()
	} 
	if parser.match(FUN) {
		return parser.funcDeclaration()
	}

	return parser.statement()
}

// statement parses various statement types (print, block, if, loops, etc.)
func (parser *Parser) statement() (Statement, error) {
	if parser.match(PRINT) {
		return parser.printStatement()
	}
	if parser.match(LEFT_BRACE) {
		stmts, err := parser.blockStatement()
		return &BlockStatement{
			stmts: stmts,
		}, err
	}
	if parser.match(IF) {
		return parser.ifStatement()
	}
	if parser.match(WHILE) {
		return parser.whileStatement()
	}
	if parser.match(FOR) {
		return parser.forStatement()
	}

	return parser.expressionStatement()
}

// varDeclaration parses variable declarations
func (parser *Parser) varDeclaration() (Statement, error) {
	name := consume(parser, "IDENTIFIER", "Expect variable name")

	var initializer Expr
	if parser.match("EQUAL") {
		initializer, _ = expression(parser)
	}

	consume(parser, "SEMICOLON", "Expect ';' after variable declaration")

	return &VarStatement{
		name: name,
		init: initializer,
	}, nil
}

func (parser *Parser) funcDeclaration() (Statement, error) {
	name := consume(parser, "IDENTIFIER", "Expect function name")
	consume(parser, LEFT_PAREN, "Expect '(' after function name")

	var args []Token
	for !parser.check(RIGHT_PAREN) {
		args = append(args, consume(parser, "IDENTIFIER", "Expect parameter name"))
		if !parser.match(COMMA) {
			break
		}
	}

	consume(parser, RIGHT_PAREN, "Expect ')' after function parameters")
	consume(parser, LEFT_BRACE, "Expect '{' before function body")
	
	body, err := parser.blockStatement()

	return &FuncStatement{
		name: name,
		args: args,
		body: body,
	}, err
}

// printStatement parses print statements
func (parser *Parser) printStatement() (Statement, error) {
	expr, err := expression(parser)
	if err != nil {
		return nil, err
	}

	consume(parser, "SEMICOLON", "Expect ';' at end")

	return &PrintStatement{
		expr: expr,
	}, nil
}

// ifStatement parses if-else conditional statements
func (parser *Parser) ifStatement() (Statement, error) {
	consume(parser, "LEFT_PAREN", "Expect '(' after 'if'")
	condition, err := expression(parser)
	if err != nil {
		return nil, err
	}
	consume(parser, "RIGHT_PAREN", "Expect ')' after if condition")

	thenStmt, _ := parser.statement()
	var elseStmt Statement
	if parser.match("ELSE") {
		elseStmt, _ = parser.statement()
	}

	return &IfStatement{
		condition: condition,
		thenStmt:  thenStmt,
		elseStmt:  elseStmt,
	}, nil
}

// whileStatement parses while loop statements
func (parser *Parser) whileStatement() (Statement, error) {
	consume(parser, LEFT_PAREN, "Expect '(' after 'while'")
	condition, err := expression(parser)
	if err != nil {
		return nil, err
	}
	consume(parser, RIGHT_PAREN, "Expect ')' after while condition")

	body, err := parser.statement()
	if err != nil {
		return nil, err
	}

	return &WhileStatement{
		condition: condition,
		body:      body,
	}, nil
}

// forStatement parses for loop statements by converting them into while loop
func (parser *Parser) forStatement() (Statement, error) {
	consume(parser, LEFT_PAREN, "Expect '(' after 'for'")

	var initialization Statement
	if parser.match(SEMICOLON) {
		initialization = nil
	} else if parser.match(VAR) {
		initialization, _ = parser.varDeclaration()
	} else {
		initialization, _ = parser.expressionStatement()
	}

	var condition, updation Expr
	if !parser.check(SEMICOLON) {
		condition, _ = expression(parser)
	} else {
		condition = &Literal{
			value: "true",
			t:     "bool",
		}
	}
	consume(parser, SEMICOLON, "Expect ';' after loop condition")

	if !parser.check(SEMICOLON) && !parser.check(RIGHT_PAREN) {
		updation, _ = expression(parser)
	}
	consume(parser, RIGHT_PAREN, "Expect ')' after for")

	body, err := parser.statement()
	if err != nil {
		return nil, err
	}

	// If there is an updation statement, add it to the end of the body block
	if updation != nil {
		body = &BlockStatement{
			stmts: []Statement{
				body,
				&ExprStatement{
					expr: updation,
				},
			},
		}
	}

	// Create a while loop with the condition and body
	body = &WhileStatement{
		condition: condition,
		body:      body,
	}

	// If there is an initialization statement, add it to the beginning
	// of the body block
	if initialization != nil {
		body = &BlockStatement{
			stmts: []Statement{
				initialization,
				body,
			},
		}
	}

	return body, nil
}

// expressionStatement parses statements that are expressions
func (parser *Parser) expressionStatement() (Statement, error) {
	expr, err := expression(parser)
	if err != nil {
		return nil, err
	}

	if !parser.match(SEMICOLON) {
		fmt.Fprintf(os.Stderr, "ERROR : Expect ';' after expression\n")
		os.Exit(70)
	}

	return &ExprStatement{
		expr: expr,
	}, nil
}

// blockStatement parses a block of statements enclosed in braces
func (parser *Parser) blockStatement() ([]Statement, error) {
	var stmts []Statement
	for !parser.isEOF() && !parser.check(RIGHT_BRACE) {
		stmt, err := parser.declaration()
		if err != nil {
			return nil, err
		}
		stmts = append(stmts, stmt)
	}
	consume(parser, RIGHT_BRACE, "Expect '}' after block")
	return stmts, nil
}

// match checks if the current token matches any of the given types and advances if true
func (parser *Parser) match(tokenTypes ...string) bool {
	if slices.ContainsFunc(tokenTypes, parser.check) {
		parser.advance()
		return true
	}

	return false
}

// check tests if the current token is of the specified type without advancing
func (parser *Parser) check(tokenType string) bool {
	if parser.current >= len(parser.tokens) {
		return false
	}

	return parser.peek().TokenType == tokenType
}

// peek returns the current token without advancing
func (parser *Parser) peek() Token {
	if parser.current >= len(parser.tokens) {
		return Token{}
	}
	return parser.tokens[parser.current]
}

// advance moves to the next token and returns the previous token
func (parser *Parser) advance() {
	parser.current += 1
}

// previous returns the previously consumed token
func (parser *Parser) previous() Token {
	return parser.tokens[parser.current-1]
}

// isEOF checks if we've reached the end of the token stream
func (parser *Parser) isEOF() bool {
	return parser.peek().TokenType == "EOF"
}

// consume advances if the current token is of the expected type, otherwise reports an error
func consume(parser *Parser, tokenType string, msg string) Token {
	token := parser.peek()
	if !parser.match(tokenType) {
		fmt.Fprintf(os.Stderr, "ERROR : %s\n", msg)
		os.Exit(65)
	}

	return token
}
