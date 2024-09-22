package main

import "fmt"

func getPrintContents(fileContents []byte) []byte {
	return fileContents[6:]
}

func run(fileContents []byte) error {
	fileContents = getPrintContents(fileContents)
	expr, err := evaluate(fileContents)
	if err != nil {
		return err
	}

	fmt.Println(expr)
	return nil
}
