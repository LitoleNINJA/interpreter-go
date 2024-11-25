package main

import (
	"fmt"
	"os"
	"strings"
)

type Expr interface {
	String() string
	Evaluate() (Value, error)
}

type Literal struct {
	value string
	t     string
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
	return fmt.Sprintf("(%s %s %s)", b.operator.lexeme, b.left, b.right)
}

type Grouping struct {
	expression Expr
}

func (g Grouping) String() string {
	return fmt.Sprintf("(group %s)", g.expression)
}

func expression(parser *Parser) (Expr, error) {
	return equality(parser)
}

func equality(parser *Parser) (Expr, error) {
	expr, err := comparison(parser)

	for parser.match(BANG_EQUAL, EQUAL_EQUAL) {
		operator := parser.previous()
		right, err := comparison(parser)
		if err != nil {
			return &Binary{}, err
		}
		expr = &Binary{
			left:     expr,
			operator: operator,
			right:    right,
		}
	}

	return expr, err
}

func comparison(parser *Parser) (Expr, error) {
	expr, err := term(parser)

	for parser.match(GREATER, GREATER_EQUAL, LESS, LESS_EQUAL) {
		operator := parser.previous()
		right, err := term(parser)
		if err != nil {
			return &Binary{}, err
		}
		expr = &Binary{
			left:     expr,
			operator: operator,
			right:    right,
		}
	}

	return expr, err
}

func term(parser *Parser) (Expr, error) {
	expr, err := factor(parser)

	for parser.match(MINUS, PLUS) {
		operator := parser.previous()
		right, err := factor(parser)
		if err != nil {
			return &Binary{}, err
		}
		expr = &Binary{
			left:     expr,
			operator: operator,
			right:    right,
		}
	}

	return expr, err
}

func factor(parser *Parser) (Expr, error) {
	expr, err := unary(parser)

	for parser.match(STAR, SLASH) {
		operator := parser.previous()
		right, err := unary(parser)
		if err != nil {
			return &Binary{}, err
		}
		expr = &Binary{
			left:     expr,
			operator: operator,
			right:    right,
		}
	}

	return expr, err
}

func unary(parser *Parser) (Expr, error) {
	if parser.match(BANG, MINUS) {
		operator := parser.previous()
		right, err := unary(parser)
		if err != nil {
			return &Binary{}, err
		}
		return &Unary{
			operator: operator,
			right:    right,
		}, err
	}

	return primary(parser)
}

func primary(parser *Parser) (Expr, error) {
	if parser.match(FALSE) {
		return &Literal{
			value: "false",
			t:     "bool",
		}, nil
	} else if parser.match(TRUE) {
		return &Literal{
			value: "true",
			t:     "bool",
		}, nil
	} else if parser.match(NIL) {
		return &Literal{
			value: "nil",
			t:     "nil",
		}, nil
	} else if parser.match(STRING) {
		return &Literal{
			value: parser.previous().literal,
			t:     "string",
		}, nil
	} else if parser.match(NUMBER) {
		return &Literal{
			value: parser.previous().literal,
			t:     "number",
		}, nil
	} else if parser.match(LEFT_PAREN) {
		expr, err := expression(parser)
		consume(parser, RIGHT_PAREN, "Expect ')' after expression.")
		return &Grouping{
			expression: expr,
		}, err
	} else if parser.match(IDENTIFIER) {
		if val, ok := values[parser.previous().lexeme]; !ok {
			exitCode = 70
			return &Literal{}, fmt.Errorf("Undefined variable '%s'.", parser.previous().lexeme)
		} else {
			if strings.ContainsAny(val, "+-*/") {
				value, err := evaluate([]byte(val))
				if err != nil {
					return &Literal{}, err
				}

				val = fmt.Sprint(value)
			}
			valType := getStringType(val)
			// fmt.Printf("value : %s, type : %s\n", val, valType)
			return &Literal{
				value: val,
				t:     valType,
			}, nil
		}
	}

	return &Grouping{}, fmt.Errorf("Error at ')': Expect expression")
}

func consume(parser *Parser, tokenType string, msg string) {
	if !parser.match(tokenType) {
		fmt.Fprintf(os.Stderr, "ERROR : %s\n", msg)
		os.Exit(65)
	}
}

func parseFile(fileContent []byte) (Expr, error) {
	parser := &Parser{
		tokens:  tokenizeFile(fileContent),
		current: 0,
	}

	// fmt.Println("Tokens : ", parser.tokens)
	expr, err := parser.parse()

	return expr, err
}
