package main

import (
	"fmt"
	"os"
)

var exitCode = 0
var fileContentString string
var line = 1

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "Usage: ./your_program.sh tokenize <filename>")
		os.Exit(1)
	}

	command := os.Args[1]

	filename := os.Args[2]
	fileContents, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	switch command {
	case "tokenize":
		tokens := tokenizeFile(fileContents)

		for _, token := range tokens {
			if token != (Token{}) {
				token.printToken()
			}
		}
	case "parse":
		expr, err := parseFile(fileContents)
		if err != nil {
			fmt.Fprintf(os.Stderr, "[line %d] %v\n", line, err)
			os.Exit(65)
		}

		fmt.Println(expr)
	case "evaluate":
		val, err := evaluate(fileContents)
		if err != nil {
			fmt.Fprintf(os.Stderr, "[line %d] %v\n", line, err)
			if exitCode == 0 {
				exitCode = 70
			}
			os.Exit(exitCode)
		}

		fmt.Println(val)
	case "run":
		err := run(fileContents)
		if err != nil {
			fmt.Fprintf(os.Stderr, "[line %d] %v\n", line, err)
			if exitCode == 0 {
				exitCode = 70
			}
			os.Exit(exitCode)
		}
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		os.Exit(1)
	}

	os.Exit(exitCode)
}
