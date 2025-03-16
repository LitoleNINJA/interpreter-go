package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Token struct {
	TokenType string
	lexeme    string
	literal   string
}

func (token *Token) setToken(args ...string) {
	token.TokenType = args[0]
	token.lexeme = args[1]
	if len(args) > 2 {
		token.literal = args[2]
	} else {
		token.literal = "null"
	}
}

func (token *Token) printToken() {
	if token.literal == "" {
		fmt.Printf("%s  %s\n", token.TokenType, token.lexeme)
	} else {
		fmt.Printf("%s %s %s\n", token.TokenType, token.lexeme, token.literal)
	}
}

func tokenizeFile(fileContents []byte) []Token {
	fileContentString = string(fileContents)

	tokens := []Token{}
	for i := 0; i < len(fileContentString); i++ {
		newToken, err := addToken(string(fileContentString[i]), &i)
		if err != nil {
			fmt.Fprintf(os.Stderr, "[line %d] %v\n", line, err)
		}
		// fmt.Printf("Token : %+v\n", newToken)
		if newToken != (Token{}) {
			tokens = append(tokens, newToken)
		}
	}

	tokens = append(tokens, Token{TokenType: "EOF", lexeme: "null"})

	return tokens
}

func addToken(ch string, index *int) (Token, error) {
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
	case " ", "\t", "\r":
		break
	case "\n":
		line++
	case `"`:
		str, err := readString(index)
		if err != nil {
			exitCode = 65
			return token, err
		} else {
			token.setToken(STRING, `"`+str+`"`, str)
		}
	default:
		if isStringDigit(ch) {
			str, frac := readNumber(index)
			if frac == "0" {
				floatVal, _ := strconv.ParseFloat(str, 64)
				intVal := int64(floatVal)
				token.setToken(NUMBER, str, strconv.FormatInt(intVal, 10)+"."+frac)
			} else {
				token.setToken(NUMBER, str, strings.Trim(str, "0"))
			}
			*index--
		} else if isIndentifierStart(ch) {
			str := readIdentifier(index)
			if _, isKeyword := keywords[str]; isKeyword {
				token.setToken(keywords[str], str)
			} else {
				token.setToken(IDENTIFIER, str)
			}
			*index--
		} else {
			exitCode = 65
			return token, fmt.Errorf("Error: Unexpected character: %s", ch)
		}
	}

	return token, nil
}

func nextToken(index int) string {
	if index < len(fileContentString)-1 {
		return string(fileContentString[index+1])
	}
	return ""
}
