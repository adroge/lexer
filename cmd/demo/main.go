package main

import (
	"context"
	"fmt"

	"github.com/adroge/lexer"
)

func main() {
	lex := lexer.Build("start text {{meta, x:y}} middle text {{pi:3.14}} final text is here.")
	lex.Run(context.Background())

	done := false
	for !done {
		token := lex.NextToken()
		fmt.Printf("token: %s\n text: \"%s\"\n", token.Type, token.Value)
		if token.Type == lexer.TokenUndefined {
			done = true
		}
	}
}
