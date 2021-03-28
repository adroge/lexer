# Lexer

This is a lexer in Go I use as a starting point for other projects.

To install:

```sh
go get -u github.com/adroge/lexer
```

This is an example of usage:

```go
package main

import (
	"context"
	"fmt"

	"github.com/adroge/lexer"
)

func main() {
	lex := lexer.Build("start {{setting, x:y}} middle {{pi:3.14}} end text.")
	lex.Run(context.Background())

	for token := lex.NextToken(); token.Type != lexer.TokenUndefined; token = lex.NextToken() {
		fmt.Printf("token: %s\n text: \"%s\"\n\n", token.Type, token.Value)
	}
}
```

Rob Pike's Lexer from his presentation was used as inspiration.
