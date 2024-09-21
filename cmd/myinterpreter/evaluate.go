package main

import (
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
	// fmt.Println(val)
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
