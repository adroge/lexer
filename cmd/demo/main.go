package main

import (
	"context"
	"fmt"

	"github.com/adroge/lexer"
)

func main() {
	lex := lexer.Create("start {{setting, x:y}} middle {{pi:3.14}} end text.")
	lex.Run(context.Background())

	for token := lex.NextToken(); token.Type != lexer.TokenUndefined; token = lex.NextToken() {
		fmt.Printf("token: %s\n text: \"%s\"\n\n", token.Type, token.Value)
	}
}
