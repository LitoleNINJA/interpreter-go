package main

import "fmt"

const (
	LITERAL  = "literal"
	UNARY    = "unary"
	BINARY   = "binary"
	GROUPING = "grouping"
)

type Expr struct {
	left     *Expr
	operator Token
	right    *Expr
	expType  string
}

func (expr *Expr) print() {
	if expr.left == nil && expr.right == nil {
		fmt.Println(expr.operator.lexeme)
		return
	}

	fmt.Print("(")
	fmt.Print(expr.operator.lexeme)
	if expr.left.expType == LITERAL {
		fmt.Print(" ", expr.left.operator.literal)
	} else {
		expr.left.print()
	}

	if expr.right.expType == LITERAL {
		fmt.Print(" ", expr.right.operator.literal)
	} else {
		expr.right.print()
	}

	fmt.Print(")\n")
}

func expression(parser *Parser) Expr {
	return equality(parser)
}

func equality(parser *Parser) Expr {
	expr := comparison(parser)

	for parser.match(BANG_EQUAL, EQUAL_EQUAL) {
		operator := parser.previous()
		right := comparison(parser)
		return Expr{
			left:     &expr,
			operator: operator,
			right:    &right,
		}
	}

	return expr
}

func comparison(parser *Parser) Expr {
	expr := term(parser)

	for parser.match(GREATER, GREATER_EQUAL, LESS, LESS_EQUAL) {
		operator := parser.previous()
		right := term(parser)
		return Expr{
			left:     &expr,
			operator: operator,
			right:    &right,
		}
	}

	return expr
}

func term(parser *Parser) Expr {
	expr := factor(parser)

	for parser.match(PLUS, MINUS) {
		operator := parser.previous()
		right := factor(parser)
		return Expr{
			left:     &expr,
			operator: operator,
			right:    &right,
		}
	}

	return expr
}

func factor(parser *Parser) Expr {
	expr := unary(parser)

	for parser.match(SLASH, STAR) {
		operator := parser.previous()
		right := unary(parser)
		return Expr{
			left:     &expr,
			operator: operator,
			right:    &right,
		}
	}

	return expr
}

func unary(parser *Parser) Expr {
	if parser.match(BANG, MINUS) {
		operator := parser.previous()
		right := unary(parser)
		expr := Expr{
			operator: operator,
			right:    &right,
		}

		return expr
	}

	return primary(parser)
}

func primary(parser *Parser) Expr {
	if parser.match(FALSE) {
		return Literal(Token{
			TokenType: FALSE,
			lexeme:    "false",
		})
	} else if parser.match(TRUE) {
		return Literal(Token{
			TokenType: TRUE,
			lexeme:    "true",
		})
	} else if parser.match(NIL) {
		return Literal(Token{
			TokenType: NIL,
			lexeme:    "nil",
		})
	} else if parser.match(NUMBER, STRING) {
		return Literal(parser.previous())
	} else if parser.match(LEFT_PAREN) {
		expr := expression(parser)
		consume(parser, RIGHT_PAREN, "Expect ')' after expression.")
		return expr
	}

	fmt.Println("Should not reach here !")
	return Expr{}
}

func Literal(token Token) Expr {
	return Expr{
		operator: token,
		expType:  LITERAL,
	}
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

	expr := parser.parse()

	return expr, nil
}
