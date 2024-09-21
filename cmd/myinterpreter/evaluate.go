package main

import (
	"fmt"
	"strconv"
)

type Value interface{}

func (l *Literal) Evaluate() Value {
	if f, err := strconv.ParseFloat(l.value, 64); err != nil {
		return l.value
	} else {
		return f
	}
}

func (u *Unary) Evaluate() Value {
	val := u.right.Evaluate()

	switch u.operator.TokenType {
	case MINUS:
		val = -val.(float64)
	case BANG:
		val = !isTruthy(val)
	}

	return val
}

func (b *Binary) Evaluate() Value {
	return b.operator.lexeme
}

func (g *Grouping) Evaluate() Value {
	return g.expression.Evaluate()
}

func evaluate(fileContents []byte) (Value, error) {
	expr, err := parseFile(fileContents)
	if err != nil {
		return nil, err
	}

	return expr.Evaluate(), nil
}

func isTruthy(val Value) bool {
	switch val := val.(type) {
	case bool:
		return val
	case string:
		if b, err := strconv.ParseBool(val); err == nil {
			return b
		}
		if f, err := strconv.ParseFloat(val, 64); err == nil {
			return f != 0
		}
		if val == "nil" {
			return false
		}

		fmt.Println("Should not reach here !")
		return false
	case float64:
		return val != 0
	default:
		fmt.Println("Unknown type !")
		return false
	}
}
