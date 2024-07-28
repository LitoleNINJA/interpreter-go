package main

import (
	"fmt"
	"os"
)

const (
	LEFT_PAREN    = "LEFT_PAREN"
	RIGHT_PAREN   = "RIGHT_PAREN"
	LEFT_BRACE    = "LEFT_BRACE"
	RIGHT_BRACE   = "RIGHT_BRACE"
	STAR          = "STAR"
	DOT           = "DOT"
	COMMA         = "COMMA"
	PLUS          = "PLUS"
	MINUS         = "MINUS"
	SEMICOLON     = "SEMICOLON"
	EQUAL         = "EQUAL"
	EQUAL_EQUAL   = "EQUAL_EQUAL"
	BANG          = "BANG"
	BANG_EQUAL    = "BANG_EQUAL"
	LESS          = "LESS"
	LESS_EQUAL    = "LESS_EQUAL"
	GREATER       = "GREATER"
	GREATER_EQUAL = "GREATER_EQUAL"
	SLASH         = "SLASH"
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

var exitCode = 0
var fileContentString string
var line = 1

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
	os.Exit(exitCode)
}

func tokenizeFile(fileContents []byte) {
	fileContentString = string(fileContents)

	tokens := []Token{}
	for i := 0; i < len(fileContentString); i++ {
		newToken := addToken(string(fileContentString[i]), &i)
		tokens = append(tokens, newToken)
	}

	for _, token := range tokens {
		if token != (Token{}) {
			token.printToken()
		}
	}
	fmt.Println("EOF  null")
}

func addToken(ch string, index *int) Token {
	token := Token{}

	switch ch {
	case "(":
		token.setToken(LEFT_PAREN, ch)
	case ")":
		token.setToken(RIGHT_PAREN, ch)
	case "{":
		token.setToken(LEFT_BRACE, ch)
	case "}":
		token.setToken(RIGHT_BRACE, ch)
	case ",":
		token.setToken(COMMA, ch)
	case ".":
		token.setToken(DOT, ch)
	case "*":
		token.setToken(STAR, ch)
	case "+":
		token.setToken(PLUS, ch)
	case "-":
		token.setToken(MINUS, ch)
	case ";":
		token.setToken(SEMICOLON, ch)
	case "=":
		if nextToken(*index) == "=" {
			token.setToken(EQUAL_EQUAL, "==")
			*index++
		} else {
			token.setToken(EQUAL, ch)
		}
	case "!":
		if nextToken(*index) == "=" {
			token.setToken(BANG_EQUAL, "!=")
			*index++
		} else {
			token.setToken(BANG, ch)
		}
	case "<":
		if nextToken(*index) == "=" {
			token.setToken(LESS_EQUAL, "<=")
			*index++
		} else {
			token.setToken(LESS, ch)
		}
	case ">":
		if nextToken(*index) == "=" {
			token.setToken(GREATER_EQUAL, ">=")
			*index++
		} else {
			token.setToken(GREATER, ch)
		}
	case "/":
		if nextToken(*index) != "/" {
			token.setToken(SLASH, ch)
		} else {
			for *index < len(fileContentString) && fileContentString[*index] != '\n' {
				*index++
			}
			line++
		}
	case " ", "\t":
		break
	case "\n":
		line++
	default:
		fmt.Fprintf(os.Stderr, "[line %d] Error: Unexpected character: %s\n", line, ch)
		exitCode = 65
	}

	return token
}

func nextToken(index int) string {
	if index < len(fileContentString)-1 {
		return string(fileContentString[index+1])
	}
	return ""
}
