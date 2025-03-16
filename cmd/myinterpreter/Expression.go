package main

import (
	"fmt"
	"strconv"
)

type Expr interface {
	String() string
	Evaluate() (Value, error)
}

type Callable interface {
	Call(args []Value) (Value, error)
}

type Literal struct {
	value string
	t     string
}

type Variable struct {
	name Token
}

type Unary struct {
	operator Token
	right    Expr
}

type Binary struct {
	left     Expr
	operator Token
	right    Expr
}

type Call struct {
	callee Expr
	args   []Expr
}

type Func struct {
	name Token
	args []Token
	body []Statement
}

type Grouping struct {
	expression Expr
}

type Assignment struct {
	name  Token
	value Expr
}

func (l Literal) String() string {
	return l.value
}

func (l *Literal) Evaluate() (Value, error) {
	if l.t == "bool" {
		return strconv.ParseBool(l.value)
	} else if l.t == "identifier" {
		// check if its a native func
		if fn, ok := native_functions[l.value]; ok {
			return fn, nil
		}

		if val, ok := currentScope.getScopeValue(l.value); !ok {
			exitCode = 70
			return nil, fmt.Errorf("Undefined variable '%s'.", l.value)
		} else {
			return val, nil
			// if strings.ContainsAny(val, "+-*/") {
			// 	value, err := evaluate([]byte(val))
			// 	if err != nil {
			// 		return nil, err
			// 	}

			// 	val = fmt.Sprint(value)
			// }
			// // fmt.Printf("value : %s, type : %s\n", val, valType)

			// valType := getStringType(val)
			// if valType == "bool" {
			// 	return strconv.ParseBool(val)
			// } else if valType == "number" {
			// 	return strconv.ParseFloat(val, 64)
			// } else {
			// 	return val, nil
			// }
		}

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

func (v *Variable) String() string {
	return v.name.lexeme
}

func (v *Variable) Evaluate() (Value, error) {
	if val, ok := currentScope.getScopeValue(v.name.lexeme); ok {
		return val, nil
	} else if val, ok = native_functions[v.name.lexeme]; ok {
		return val, nil
	} else {
		return nil, fmt.Errorf("Undefined variable '%s'.", v.name.lexeme)
	}
}

func (u Unary) String() string {
	return fmt.Sprintf("(%s %s)", u.operator.lexeme, u.right)
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

func (b Binary) String() string {
	return fmt.Sprintf("(%s %s %s)", b.operator.lexeme, b.left, b.right)
}

func (b *Binary) Evaluate() (Value, error) {
	if b.operator.TokenType == OR {
		return evaluateOrExpr(b)
	} else if b.operator.TokenType == AND {
		return evaluateAndExpr(b)
	}

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
			exitCode = 70
			return nil, fmt.Errorf("operands must be numbers.\n[line 1]")
		}
		right, ok = rightVal.(float64)
		if !ok {
			exitCode = 70
			return nil, fmt.Errorf("operands must be numbers.\n[line 1]")
		}
		return left - right, nil
	case STAR:
		var left, right float64
		left, ok := leftVal.(float64)
		if !ok {
			exitCode = 70
			return nil, fmt.Errorf("operands must be numbers.\n[line 1]")
		}
		right, ok = rightVal.(float64)
		if !ok {
			exitCode = 70
			return nil, fmt.Errorf("operands must be numbers.\n[line 1]")
		}
		return left * right, nil
	case SLASH:
		var left, right float64
		left, ok := leftVal.(float64)
		if !ok {
			exitCode = 70
			return nil, fmt.Errorf("operands must be numbers.\n[line 1]")
		}
		right, ok = rightVal.(float64)
		if !ok {
			exitCode = 70
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
	default:
		return nil, fmt.Errorf("unknown operator : %s", b.operator.literal)
	}
}

func (c Call) String() string {
	return fmt.Sprintf("fn(%s %s)", c.callee, c.args)
}

func (c Call) Evaluate() (Value, error) {
	calleeVal, err := c.callee.Evaluate()
	if err != nil {
		return nil, err
	}

	argsVals := make([]Value, len(c.args))
	for i, arg := range c.args {
		val, err := arg.Evaluate()
		if err != nil {
			return nil, err
		}
		argsVals[i] = val
	}

	function, ok := calleeVal.(Callable)
	if !ok {
		return nil, fmt.Errorf("can only call functions")
	}

	return function.Call(argsVals)
}

func (f Func) String() string {
	return fmt.Sprintf("(fn %s %s)", f.name.lexeme, f.args)
}

func (f Func) Evaluate() (Value, error) {
	fn := &FunctionExpr{
		expr: f, 
		scope: NewScope(currentScope),
	}

	if f.name.lexeme != "" {
		fn.scope.setScopeValue(f.name.lexeme, fn)
	}

	return fn, nil
}

func (g Grouping) String() string {
	return fmt.Sprintf("(group %s)", g.expression)
}

func (g *Grouping) Evaluate() (Value, error) {
	return g.expression.Evaluate()
}

func (a Assignment) String() string {
	return fmt.Sprintf("(%s = %s)", a.name.lexeme, a.value)
}

func (a *Assignment) Evaluate() (Value, error) {
	val, err := a.value.Evaluate()
	if err != nil {
		return nil, err
	}

	if success := currentScope.assignScopeValue(a.name.lexeme, val); !success {
		currentScope.setScopeValue(a.name.lexeme, val)
	}

	return val, nil
}

func evaluateOrExpr(b *Binary) (Value, error) {
	leftVal, err := b.left.Evaluate()
	if err != nil {
		return nil, err
	}

	if isTruthy(leftVal) {
		return leftVal, nil
	} 

	rightVal, err := b.right.Evaluate()
	if err != nil {
		return nil, err
	}

	if isTruthy(rightVal) {
		return rightVal, nil
	}

	return false, nil
}

func evaluateAndExpr(b *Binary) (Value, error) {
	leftVal, err := b.left.Evaluate()
	if err != nil {
		return nil, err
	}
	if !isTruthy(leftVal) {
		return leftVal, nil
	}

	rightVal, err := b.right.Evaluate()
	if err != nil {
		return nil, err
	}
	
	
	return rightVal, nil
}