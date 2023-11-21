package parser

import (
	"fmt"
	"jp/lexer"
	"jp/token"
	"slices"
)

func IsValid(json string, log bool) bool {
	tok := token.Token{
		Type:    token.ILLEGAL,
		Literal: "",
	}

	l := lexer.New(json)

	for {
		nextToken := l.NextToken()
		if log {
			fmt.Printf("%+v\n", nextToken)
		}

		if nextToken.Type == token.ILLEGAL {
			return false
		}
		if nextToken.Type == token.EOF {
			break
		}

		tok.Type = nextToken.Type
		tok.Literal = nextToken.Literal
	}

	return slices.Contains([]token.TokenType{token.RBRACE, token.RBRACKET}, tok.Type)
}
