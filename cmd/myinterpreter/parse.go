package main

import "fmt"

type Expr interface{}

type Literal struct {
	value string
}

func (l Literal) String() string {
	return l.value
}

type Unary struct {
	operator Token
	right    Expr
}

func (u Unary) String() string {
	return fmt.Sprintf("(%s %s)", u.operator.lexeme, u.right)
}

type Binary struct {
	left     Expr
	operator Token
	right    Expr
}

func (b Binary) String() string {
	var operator string
	switch b.operator.TokenType {
	case STAR:
		operator = "*"
	case SLASH:
		operator = "/"
	case PLUS:
		operator = "+"
	case MINUS:
		operator = "-"
	}
	return fmt.Sprintf("(%s %s %s)", operator, b.left, b.right)
}

type Grouping struct {
	expression Expr
}

func (g Grouping) String() string {
	return fmt.Sprintf("(group %s)", g.expression)
}

func expression(parser *Parser) Expr {
	return equality(parser)
}

func equality(parser *Parser) Expr {
	expr := comparison(parser)

	for parser.match(BANG_EQUAL, EQUAL_EQUAL) {
		operator := parser.previous()
		right := comparison(parser)
		return Binary{
			left:     expr,
			operator: operator,
			right:    right,
		}
	}

	return expr
}

func comparison(parser *Parser) Expr {
	expr := term(parser)

	for parser.match(GREATER, GREATER_EQUAL, LESS, LESS_EQUAL) {
		operator := parser.previous()
		right := term(parser)
		return Binary{
			left:     expr,
			operator: operator,
			right:    right,
		}
	}

	return expr
}

func term(parser *Parser) Expr {
	expr := factor(parser)

	for parser.match(PLUS, MINUS) {
		operator := parser.previous()
		right := factor(parser)
		return Binary{
			left:     expr,
			operator: operator,
			right:    right,
		}
	}

	return expr
}

func factor(parser *Parser) Expr {
	expr := unary(parser)

	for parser.match(SLASH, STAR) {
		operator := parser.previous()
		right := unary(parser)
		return Binary{
			left:     expr,
			operator: operator,
			right:    right,
		}
	}

	return expr
}

func unary(parser *Parser) Expr {
	if parser.match(BANG, MINUS) {
		operator := parser.previous()
		right := unary(parser)
		return Unary{
			operator: operator,
			right:    right,
		}
	}

	return primary(parser)
}

func primary(parser *Parser) Expr {
	if parser.match(FALSE) {
		return Literal{
			value: "false",
		}
	} else if parser.match(TRUE) {
		return Literal{
			value: "true",
		}
	} else if parser.match(NIL) {
		return Literal{
			value: "nil",
		}
	} else if parser.match(NUMBER, STRING) {
		return Literal{value: parser.previous().literal}
	} else if parser.match(LEFT_PAREN) {
		expr := expression(parser)
		consume(parser, RIGHT_PAREN, "Expect ')' after expression.")
		return Grouping{
			expression: expr,
		}
	}

	fmt.Println("Should not reach here !")
	return Literal{}
}

func consume(parser *Parser, tokenType string, msg string) {
	if parser.match(tokenType) {
		parser.advance()
	}

	fmt.Errorf("ERROR : %s", msg)
}

func parseFile(fileContent []byte) (Expr, error) {
	parser := &Parser{
		tokens:  tokenizeFile(fileContent),
		current: 0,
	}

	// fmt.Println(parser.tokens)
	expr := parser.parse()

	return expr, nil
}
