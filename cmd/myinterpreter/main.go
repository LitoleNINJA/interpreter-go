package main

import (
	"fmt"
	"os"
)

type Token struct {
	TokenType string
	lexeme    string
	literal   struct{}
}

func (token *Token) setToken(tokenType string, text string) {
	token.TokenType = tokenType
	token.lexeme = text
	token.literal = struct{}{}
}

func (token *Token) printToken() {
	fmt.Printf("%s %s null\n", token.TokenType, token.lexeme)
}

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "Usage: ./your_program.sh tokenize <filename>")
		os.Exit(1)
	}

	command := os.Args[1]

	if command != "tokenize" {
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		os.Exit(1)
	}

	filename := os.Args[2]
	fileContents, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	tokenizeFile(fileContents)
}

func tokenizeFile(fileContents []byte) {
	fileContentString := string(fileContents)

	tokens := []Token{}
	for i := 0; i < len(fileContentString); i++ {
		newToken := addToken(string(fileContentString[i]))
		tokens = append(tokens, newToken)
	}

	for _, token := range tokens {
		token.printToken()
	}
	fmt.Println("EOF  null")
}

func addToken(ch string) Token {
	token := Token{}

	switch ch {
	case "(":
		token.setToken("LEFT_PAREN", "(")
	case ")":
		token.setToken("RIGHT_PAREN", ")")
	default:
		fmt.Println("Unknown char : ", ch)
	}

	return token
}
