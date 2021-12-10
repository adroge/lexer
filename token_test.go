package lexer_test

import (
	"testing"

	"github.com/adroge/lexer"
	"github.com/stretchr/testify/assert"
)

func TestTokenTypeStringUndefined(t *testing.T) {
	tok := lexer.TokenUndefined
	assert.Equal(t, "Undefined", tok.String())
}

func TestTokenTypeStringPlainText(t *testing.T) {
	tok := lexer.TokenPlainText
	assert.Equal(t, "PlainText", tok.String())
}

func TestTokenTypeStringLeftMeta(t *testing.T) {
	tok := lexer.TokenLeftMeta
	assert.Equal(t, "LeftMeta", tok.String())
}

func TestTokenTypeStringMetaIdentifier(t *testing.T) {
	tok := lexer.TokenMetaIdentifier
	assert.Equal(t, "MetaIdentifier", tok.String())
}

func TestTokenTypeStringMetaNumberValue(t *testing.T) {
	tok := lexer.TokenMetaNumberValue
	assert.Equal(t, "MetaNumberValue", tok.String())
}

func TestTokenTypeStringMetaTextValue(t *testing.T) {
	tok := lexer.TokenMetaTextValue
	assert.Equal(t, "MetaTextValue", tok.String())
}

func TestTokenTypeStringRightMeta(t *testing.T) {
	tok := lexer.TokenRightMeta
	assert.Equal(t, "RightMeta", tok.String())
}

func TestTokenTypeStringError(t *testing.T) {
	tok := lexer.TokenError
	assert.Equal(t, "Error", tok.String())
}

func TestTokenTypeStringEof(t *testing.T) {
	tok := lexer.TokenEof
	assert.Equal(t, "Eof", tok.String())
}

func TestTokenTypeStringInvalid(t *testing.T) {
	tok := lexer.TokenEof + 10000
	assert.Equal(t, "invalid", tok.String())
}
