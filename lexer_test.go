package lexer_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/adroge/lexer"
)

func TestTokensDone(t *testing.T) {
	l := lexer.Build("start text {{meta}}")
	l.Run(context.Background())

	l.NextToken()          // text
	l.NextToken()          // open left meta
	l.NextToken()          // meta
	l.NextToken()          // close right meta
	token := l.NextToken() // eof
	assert.Equal(t, lexer.TokenEof, token.Type)
	token = l.NextToken() // undefined
	assert.Equal(t, lexer.TokenUndefined, token.Type)
}

func TestContextCancel(t *testing.T) {
	l := lexer.Build("start text {{meta}} more text {{meta}} final text is here.")

	ctx, cancel := context.WithCancel(context.Background())
	l.Run(ctx)

	token := l.NextToken() // text
	assert.Equal(t, lexer.TokenPlainText, token.Type)

	cancel()

	// done may be true or false on some of the following
	l.NextToken()         // open left meta
	l.NextToken()         // meta
	l.NextToken()         // close right meta
	token = l.NextToken() // undefined - but should be more text if not canceled
	assert.Equal(t, lexer.TokenUndefined, token.Type)
}

func TestBasic(t *testing.T) {
	l := lexer.Build("x{{y}}z")

	l.Run(context.Background())

	token := l.NextToken()
	assert.Equal(t, "x", token.String())
	assert.Equal(t, lexer.TokenPlainText, token.Type)

	token = l.NextToken()
	assert.Equal(t, "{{", token.String())
	assert.Equal(t, lexer.TokenLeftMeta, token.Type)

	token = l.NextToken()
	assert.Equal(t, "y", token.String())
	assert.Equal(t, lexer.TokenMetaIdentifier, token.Type)

	token = l.NextToken()
	assert.Equal(t, "}}", token.String())
	assert.Equal(t, lexer.TokenRightMeta, token.Type)

	token = l.NextToken()
	assert.Equal(t, "z", token.String())
	assert.Equal(t, lexer.TokenPlainText, token.Type)

	token = l.NextToken()
	assert.Equal(t, "EOF", token.String())
	assert.Equal(t, "Eof", token.Type.String())
	assert.Equal(t, lexer.TokenEof, token.Type)

	token = l.NextToken()
	assert.Equal(t, "", token.String())
	assert.Equal(t, lexer.TokenUndefined, token.Type)
	assert.Equal(t, "Undefined", token.Type.String())
}

func TestMeta(t *testing.T) {
	l := lexer.Build("{{ abcd : def, R : 10\tG : 144 B : 0x3 }}")

	l.Run(context.Background())

	token := l.NextToken()
	assert.Equal(t, lexer.TokenLeftMeta, token.Type)

	token = l.NextToken()
	assert.Equal(t, lexer.TokenMetaIdentifier, token.Type)
	assert.Equal(t, "abcd", token.String())

	token = l.NextToken()
	assert.Equal(t, lexer.TokenMetaTextValue, token.Type)
	assert.Equal(t, "def", token.String())

	token = l.NextToken()
	assert.Equal(t, lexer.TokenMetaIdentifier, token.Type)
	assert.Equal(t, "R", token.String())

	token = l.NextToken()
	assert.Equal(t, lexer.TokenMetaNumberValue, token.Type)
	assert.Equal(t, "10", token.String())

	token = l.NextToken()
	assert.Equal(t, lexer.TokenMetaIdentifier, token.Type)
	assert.Equal(t, "G", token.String())

	token = l.NextToken()
	assert.Equal(t, lexer.TokenMetaNumberValue, token.Type)
	assert.Equal(t, "144", token.String())

	token = l.NextToken()
	assert.Equal(t, lexer.TokenMetaIdentifier, token.Type)
	assert.Equal(t, "B", token.String())

	token = l.NextToken()
	assert.Equal(t, lexer.TokenMetaNumberValue, token.Type)
	assert.Equal(t, "0x3", token.String())

	token = l.NextToken()
	assert.Equal(t, lexer.TokenRightMeta, token.Type)

	token = l.NextToken()
	assert.Equal(t, lexer.TokenEof, token.Type)

	token = l.NextToken()
	assert.Equal(t, lexer.TokenUndefined, token.Type)
}

func TestUnclosedMeta(t *testing.T) {
	l := lexer.Build("{{z")

	l.Run(context.Background())

	token := l.NextToken()
	assert.Equal(t, lexer.TokenLeftMeta, token.Type)

	token = l.NextToken()
	assert.Equal(t, lexer.TokenMetaIdentifier, token.Type)

	token = l.NextToken()
	assert.Equal(t, lexer.TokenError, token.Type)
	assert.Equal(t, "unclosed meta", token.String())

	token = l.NextToken()
	assert.Equal(t, lexer.TokenUndefined, token.Type)
}

func TestEndWithLeftMeta(t *testing.T) {
	l := lexer.Build("{{")

	l.Run(context.Background())

	token := l.NextToken()
	assert.Equal(t, lexer.TokenLeftMeta, token.Type)

	token = l.NextToken()
	assert.Equal(t, lexer.TokenError, token.Type)
	assert.Equal(t, "unclosed meta", token.String())

	token = l.NextToken()
	assert.Equal(t, lexer.TokenUndefined, token.Type)
}

func TestUnclosedMetaOnValue(t *testing.T) {
	l := lexer.Build("{{ z:")

	l.Run(context.Background())

	token := l.NextToken()
	assert.Equal(t, lexer.TokenLeftMeta, token.Type)

	token = l.NextToken()
	assert.Equal(t, lexer.TokenMetaIdentifier, token.Type)

	token = l.NextToken()
	assert.Equal(t, lexer.TokenError, token.Type)
	assert.Equal(t, "unclosed meta", token.String())

	token = l.NextToken()
	assert.Equal(t, lexer.TokenUndefined, token.Type)
}

func TestBadMetaIdentifierValue(t *testing.T) {
	l := lexer.Build("{{z:*}}")

	l.Run(context.Background())

	token := l.NextToken()
	assert.Equal(t, lexer.TokenLeftMeta, token.Type)

	token = l.NextToken()
	assert.Equal(t, lexer.TokenMetaIdentifier, token.Type)

	token = l.NextToken()
	assert.Equal(t, lexer.TokenError, token.Type)
	assert.Equal(t, "value syntax: \"*\"", token.String())

	token = l.NextToken()
	assert.Equal(t, lexer.TokenUndefined, token.Type)
}

func TestNumericIdentifier(t *testing.T) {
	l := lexer.Build("{{12}}")

	l.Run(context.Background())

	token := l.NextToken()
	assert.Equal(t, lexer.TokenLeftMeta, token.Type)

	token = l.NextToken()
	assert.Equal(t, lexer.TokenError, token.Type)
	assert.Equal(t, "identifier syntax: \"1\"", token.Value)
}

func TestNumericIdentifierValueError(t *testing.T) {
	l := lexer.Build("{{a:12]}}")

	l.Run(context.Background())

	token := l.NextToken()
	assert.Equal(t, lexer.TokenLeftMeta, token.Type)

	token = l.NextToken()
	assert.Equal(t, lexer.TokenMetaIdentifier, token.Type)
	assert.Equal(t, "a", token.Value)

	token = l.NextToken()
	assert.Equal(t, lexer.TokenMetaNumberValue, token.Type)
	assert.Equal(t, "12", token.Value)

	token = l.NextToken()
	assert.Equal(t, lexer.TokenError, token.Type)
	assert.Equal(t, "identifier syntax: \"]\"", token.Value)

	token = l.NextToken()
	assert.Equal(t, lexer.TokenUndefined, token.Type)
}

func TestNumericIdentifierValueBadNumber(t *testing.T) {
	l := lexer.Build("{{a:12e12}}")

	l.Run(context.Background())

	token := l.NextToken()
	assert.Equal(t, lexer.TokenLeftMeta, token.Type)

	token = l.NextToken()
	assert.Equal(t, lexer.TokenMetaIdentifier, token.Type)
	assert.Equal(t, "a", token.Value)

	token = l.NextToken()
	assert.Equal(t, lexer.TokenError, token.Type)
	assert.Equal(t, "number syntax: \"12e\"", token.Value)

	token = l.NextToken()
	assert.Equal(t, lexer.TokenUndefined, token.Type)
}

func TestAcceptMetaFloatNumber(t *testing.T) {
	l := lexer.Build("{{pi:3.14}}")

	l.Run(context.Background())

	token := l.NextToken()
	assert.Equal(t, lexer.TokenLeftMeta, token.Type)

	token = l.NextToken()
	assert.Equal(t, lexer.TokenMetaIdentifier, token.Type)
	assert.Equal(t, "pi", token.Value)

	token = l.NextToken()
	assert.Equal(t, lexer.TokenMetaNumberValue, token.Type)
	assert.Equal(t, "3.14", token.Value)

	token = l.NextToken()
	assert.Equal(t, lexer.TokenRightMeta, token.Type)

	token = l.NextToken()
	assert.Equal(t, lexer.TokenEof, token.Type)

	token = l.NextToken()
	assert.Equal(t, lexer.TokenUndefined, token.Type)
}

func TestTwoIdentifiers(t *testing.T) {
	l := lexer.Build("{{a, b}}")

	l.Run(context.Background())

	token := l.NextToken()
	assert.Equal(t, lexer.TokenLeftMeta, token.Type)

	token = l.NextToken()
	assert.Equal(t, lexer.TokenMetaIdentifier, token.Type)
	assert.Equal(t, "a", token.Value)

	token = l.NextToken()
	assert.Equal(t, lexer.TokenMetaIdentifier, token.Type)
	assert.Equal(t, "b", token.Value)

	token = l.NextToken()
	assert.Equal(t, lexer.TokenRightMeta, token.Type)

	token = l.NextToken()
	assert.Equal(t, lexer.TokenEof, token.Type)

	token = l.NextToken()
	assert.Equal(t, lexer.TokenUndefined, token.Type)
}
