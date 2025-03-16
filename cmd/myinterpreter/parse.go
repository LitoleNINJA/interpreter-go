package main

import (
	"fmt"
)

func expression(parser *Parser) (Expr, error) {
	return assignment(parser)
}

func assignment(parser *Parser) (Expr, error) {
	expr, err := equality(parser)

	for parser.match(EQUAL) {
		value, err := expression(parser)
		if err != nil {
			return nil, err
		}

		if varExpr, ok := expr.(*Variable); ok {
			return &Assignment{
				name:  varExpr.name,
				value: value,
			}, nil
		} else {
			return nil, fmt.Errorf("Invalid assignment target")
		}
	}

	return expr, err
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

	for parser.match(GREATER, GREATER_EQUAL, LESS, LESS_EQUAL, OR, AND) {
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

	return call(parser)
}

func call(parser *Parser) (Expr, error) {
	expr, err := primary(parser)
	if err != nil {
		return nil, err
	}

	for parser.match(LEFT_PAREN) {
		expr, err = finishCall(parser, expr)
		if err != nil {
			return nil, err
		}
	}

	return expr, nil
}

// finishCall parses the arguments of a function call
func finishCall(parser *Parser, callee Expr) (Expr, error) {
	args := make([]Expr, 0)

	for !parser.check(RIGHT_PAREN) {
		expr, err := expression(parser)
		if err != nil {
			return nil, err
		}

		args = append(args, expr)
		if !parser.match(COMMA) {
			break
		}
	}

	consume(parser, RIGHT_PAREN, "Expect ')' after arguments.")

	return &Call{
		callee: callee,
		args:   args,
	}, nil
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
		if err != nil {
			return &Grouping{}, err
		}

		consume(parser, RIGHT_PAREN, "Expect ')' after expression.")
		return &Grouping{
			expression: expr,
		}, err
	} else if parser.match(IDENTIFIER) {
		return &Variable{
			name: parser.previous(),
		}, nil
	} else if parser.match(FUN) {
		return functionalExpression(parser)
	}

	exitCode = 65
	return nil, fmt.Errorf("Error at ')': Expect expression")
}

func functionalExpression(parser *Parser) (Expr, error) {
	var name Token 
	if !parser.check(LEFT_PAREN) {
		name = consume(parser, IDENTIFIER, "Expect function name")
	}

	args, err := argsList(parser)
	if err != nil {
		return nil, err
	}

	consume(parser, LEFT_BRACE, "Expect '{' before function body")
	body, err := parser.blockStatement()

	return &Func{
		name: name,
		args: args,
		body: body,
	}, err
}

func argsList(parser *Parser) ([]Token, error) {
	consume(parser, LEFT_PAREN, "Expect '(' after function name")

	args := make([]Token, 0)
	for !parser.check(RIGHT_PAREN) {
		if len(args) > max_args_count {
			return nil, fmt.Errorf("cannot have more than 128 arguments")
		}

		args = append(args, consume(parser, IDENTIFIER, "Expect parameter name"))
		if !parser.match(COMMA) {
			break
		}
	}

	consume(parser, RIGHT_PAREN, "Expect ')' after function parameters")
	return args, nil
}

func parseFile(fileContent []byte) (Expr, error) {
	parser := &Parser{
		tokens:  tokenizeFile(fileContent),
		current: 0,
	}

	expr, err := expression(parser)
	return expr, err
}
