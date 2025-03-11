package main

import (
	"fmt"
	"strconv"
)

type Value any

func (l *Literal) Evaluate() (Value, error) {
	if l.t == "bool" {
		return strconv.ParseBool(l.value)
	} else if l.t != "number" {
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
			return nil, fmt.Errorf("operand must be a number.\n[line 1]")
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
		return add(leftVal, rightVal)
	case MINUS:
		var left, right float64
		left, ok := leftVal.(float64)
		if !ok {
			return nil, fmt.Errorf("operands must be numbers.\n[line 1]")
		}
		right, ok = rightVal.(float64)
		if !ok {
			return nil, fmt.Errorf("operands must be numbers.\n[line 1]")
		}
		return left - right, nil
	case STAR:
		var left, right float64
		left, ok := leftVal.(float64)
		if !ok {
			return nil, fmt.Errorf("operands must be numbers.\n[line 1]")
		}
		right, ok = rightVal.(float64)
		if !ok {
			return nil, fmt.Errorf("operands must be numbers.\n[line 1]")
		}
		return left * right, nil
	case SLASH:
		var left, right float64
		left, ok := leftVal.(float64)
		if !ok {
			return nil, fmt.Errorf("operands must be numbers.\n[line 1]")
		}
		right, ok = rightVal.(float64)
		if !ok {
			return nil, fmt.Errorf("operands must be numbers.\n[line 1]")
		}
		return left / right, nil
	case GREATER:
		err := checkBothNumber(leftVal, rightVal)
		if err != nil {
			return nil, err
		}
		return leftVal.(float64) > rightVal.(float64), nil
	case GREATER_EQUAL:
		err := checkBothNumber(leftVal, rightVal)
		if err != nil {
			return nil, err
		}
		return leftVal.(float64) >= rightVal.(float64), nil
	case LESS:
		err := checkBothNumber(leftVal, rightVal)
		if err != nil {
			return nil, err
		}
		return leftVal.(float64) < rightVal.(float64), nil
	case LESS_EQUAL:
		err := checkBothNumber(leftVal, rightVal)
		if err != nil {
			return nil, err
		}
		return leftVal.(float64) <= rightVal.(float64), nil
	case EQUAL_EQUAL:
		return checkEqual(leftVal, rightVal), nil
	case BANG_EQUAL:
		return !checkEqual(leftVal, rightVal), nil
	case OR:
		if isTruthy(leftVal) {
			return leftVal, nil
		} else if isTruthy(rightVal) {
			return rightVal, nil
		} else {
			return false, nil
		}
	default:
		return nil, fmt.Errorf("unknown operator : %s", b.operator.literal)
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

		return true
	case float64:
		return val != 0
	default:
		fmt.Println("Unknown type !")
		return false
	}
}

func add(left Value, right Value) (Value, error) {
	switch left.(type) {
	case float64:
		leftVal := left.(float64)
		if rightVal, ok := right.(float64); !ok {
			return nil, fmt.Errorf("operands must be numbers.\n[line 1]")
		} else {
			return leftVal + rightVal, nil
		}
	case string:
		leftVal := left.(string)
		if rightVal, ok := right.(string); !ok {
			return nil, fmt.Errorf("operands must be numbers.\n[line 1]")
		} else {
			return leftVal + rightVal, nil
		}
	default:
		return nil, fmt.Errorf("operands must be numbers.\n[line 1]")
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
	case bool:
		if right, ok := rightVal.(bool); !ok {
			return false
		} else {
			return left == right
		}
	default:
		fmt.Println("Type mismatch !")
		return false
	}
}

func checkBothNumber(leftVal Value, rightVal Value) error {
	switch leftVal.(type) {
	case float64:
		if _, ok := rightVal.(float64); !ok {
			return fmt.Errorf("operands must be numbers.\n[line 1]")
		} else {
			return nil
		}
	default:
		return fmt.Errorf("operands must be numbers.\n[line 1]")
	}
}
