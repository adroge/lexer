# Lexer

This is a reusable lexer.

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

Another source of usage are the unit tests.

Rob Pike's Lexer from his presentation was used as inspiration.
