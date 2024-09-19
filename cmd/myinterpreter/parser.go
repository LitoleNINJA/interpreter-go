package main

type Parser struct {
	tokens  []Token
	current int
}

func (parser *Parser) parse() Expr {
	return expression(parser)
}

func (parser *Parser) match(tokenTypes ...string) bool {
	for _, tokenType := range tokenTypes {
		if parser.check(tokenType) {
			// fmt.Printf("Match : %s and %s : at %d\n", tokenType, parser.peek().TokenType, parser.current)
			parser.advance()
			return true
		}
	}

	return false
}

func (parser *Parser) check(tokenType string) bool {
	if parser.current >= len(parser.tokens) {
		return false
	}

	return parser.peek().TokenType == tokenType
}

func (parser *Parser) peek() Token {
	return parser.tokens[parser.current]
}

func (parser *Parser) advance() Token {
	parser.current += 1
	return parser.previous()
}

func (parser *Parser) previous() Token {
	return parser.tokens[parser.current-1]
}
