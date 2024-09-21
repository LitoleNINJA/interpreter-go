package main

import (
	"fmt"
	"strconv"
	"strings"
)

type Value interface{}

func (l *Literal) Evaluate() (Value, error) {
	if l.t != "number" {
		return l.value, nil
	} else {
		val, err := strconv.ParseFloat(l.value, 64)
		if err != nil {
			return nil, err
		}
		return val, nil
	}
}

func (u *Unary) Evaluate() (Value, error) {
	val, err := u.right.Evaluate()
	if err != nil {
		return nil, err
	}

	switch u.operator.TokenType {
	case MINUS:
		if valFloat, ok := val.(float64); ok {
			val = -valFloat
		} else {
			return nil, fmt.Errorf("Operand must be a number.\n[line 1]")
		}
	case BANG:
		val = !isTruthy(val)
	}

	return val, nil
}

func (b *Binary) Evaluate() (Value, error) {
	leftVal, err := b.left.Evaluate()
	if err != nil {
		return nil, err
	}
	rightVal, err := b.right.Evaluate()
	if err != nil {
		return nil, err
	}

	switch b.operator.TokenType {
	case PLUS:
		return add(leftVal, rightVal), nil
	case MINUS:
		return leftVal.(float64) - rightVal.(float64), nil
	case STAR:
		return leftVal.(float64) * rightVal.(float64), nil
	case SLASH:
		return leftVal.(float64) / rightVal.(float64), nil
	case GREATER:
		return leftVal.(float64) > rightVal.(float64), nil
	case GREATER_EQUAL:
		return leftVal.(float64) >= rightVal.(float64), nil
	case LESS:
		return leftVal.(float64) < rightVal.(float64), nil
	case LESS_EQUAL:
		return leftVal.(float64) <= rightVal.(float64), nil
	case EQUAL_EQUAL:
		return checkEqual(leftVal, rightVal), nil
	case BANG_EQUAL:
		return !checkEqual(leftVal, rightVal), nil
	default:
		return nil, fmt.Errorf("Unknown Operator !")
	}
}

func (g *Grouping) Evaluate() (Value, error) {
	return g.expression.Evaluate()
}

func evaluate(fileContents []byte) (Value, error) {
	expr, err := parseFile(fileContents)
	if err != nil {
		return nil, err
	}

	return expr.Evaluate()
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

func checkEqual(leftVal Value, rightVal Value) bool {
	switch left := leftVal.(type) {
	case float64:
		if right, ok := rightVal.(float64); !ok {
			return false
		} else {
			return left == right
		}
	case string:
		if right, ok := rightVal.(string); !ok {
			return false
		} else {
			return left == right
		}
	default:
		fmt.Println("Type mismatch !")
		return false
	}
}
