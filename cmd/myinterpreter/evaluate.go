package main

type Value any

func evaluate(fileContents []byte) (Value, error) {
	expr, err := parseFile(fileContents)
	if err != nil {
		exitCode = 65
		return nil, err
	}

	val, err := expr.Evaluate()
	return val, err
}
