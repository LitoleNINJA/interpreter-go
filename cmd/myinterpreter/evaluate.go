package main

import (
	"fmt"
	"strconv"
	"strings"
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
	leftVal := b.left.Evaluate()
	rightVal := b.right.Evaluate()

	switch b.operator.TokenType {
	case PLUS:
		return add(leftVal, rightVal)
	case MINUS:
		return leftVal.(float64) - rightVal.(float64)
	case STAR:
		return leftVal.(float64) * rightVal.(float64)
	case SLASH:
		return leftVal.(float64) / rightVal.(float64)
	case GREATER:
		return leftVal.(float64) > rightVal.(float64)
	case GREATER_EQUAL:
		return leftVal.(float64) >= rightVal.(float64)
	case LESS:
		return leftVal.(float64) < rightVal.(float64)
	case LESS_EQUAL:
		return leftVal.(float64) <= rightVal.(float64)
	default:
		fmt.Println("Unknown operator!")
		return nil
	}
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

func add(left Value, right Value) Value {
	switch left.(type) {
	case float64:
		return left.(float64) + right.(float64)
	case string:
		leftVal := left.(string)
		var rightVal string
		switch right := right.(type) {
		case string:
			rightVal = right
		case float64:
			rightVal = fmt.Sprintf("%f", right)
			rightVal = strings.Trim(rightVal, "0")
			rightVal = strings.Trim(rightVal, ".")
		}
		return leftVal + rightVal
	default:
		return nil
	}
}
