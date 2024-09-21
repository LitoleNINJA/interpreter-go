package main

type Value interface{}

func (l *Literal) Evaluate() Value {
	return l.value
}

func (u *Unary) Evaluate() Value {
	return u.operator.lexeme
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
