package main

func run(fileContents []byte) error {
	init_native_functions()

	parser := &Parser{
		tokens:  tokenizeFile(fileContents),
		current: 0,
	}

	stmts, err := parser.parse()
	if err != nil {
		return err
	}

	for _, stmt := range stmts {
		_, err := stmt.Execute()
		if err != nil {
			return err
		}
	}

	return nil
}
