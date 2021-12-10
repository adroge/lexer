package lexer_test

import (
	"context"
	"errors"
	"testing"

	"github.com/adroge/lexer"
	"github.com/stretchr/testify/assert"
)

func TestSetMetaBadLeftMeta(t *testing.T) {
	err := lexer.SetMeta("", ">>", '=', '|')
	assert.NotNil(t, err)
	assert.True(t, errors.Is(lexer.ErrMetaZeroLength, err))
}

func TestSetMetaBadRightMeta(t *testing.T) {
	err := lexer.SetMeta("<<", "", '=', '|')
	assert.NotNil(t, err)
	assert.True(t, errors.Is(lexer.ErrMetaZeroLength, err))
}

func TestSetMetaIndicatorMatch(t *testing.T) {
	err := lexer.SetMeta("<<", ">>", '=', '=')
	assert.NotNil(t, err)
	assert.True(t, errors.Is(lexer.ErrMetaIndicatorMatch, err))
}

func TestSetMetaWorking(t *testing.T) {
	err := lexer.SetMeta("<<", ">>", '=', '|')
	assert.Nil(t, err)

	l := lexer.Create("text <<a=5|b=abc>> end.")
	l.Run(context.Background())

	token := l.NextToken()
	assert.Equal(t, "text ", token.Value)

	token = l.NextToken()
	assert.Equal(t, lexer.TokenLeftMeta, token.Type)

	token = l.NextToken()
	assert.Equal(t, "a", token.Value)
	assert.Equal(t, lexer.TokenMetaIdentifier, token.Type)

	token = l.NextToken()
	assert.Equal(t, "5", token.Value)
	assert.Equal(t, lexer.TokenMetaNumberValue, token.Type)

	token = l.NextToken()
	assert.Equal(t, "b", token.Value)
	assert.Equal(t, lexer.TokenMetaIdentifier, token.Type)

	token = l.NextToken()
	assert.Equal(t, "abc", token.Value)
	assert.Equal(t, lexer.TokenMetaTextValue, token.Type)

	token = l.NextToken()
	assert.Equal(t, lexer.TokenRightMeta, token.Type)

	token = l.NextToken()
	assert.Equal(t, " end.", token.Value)

	token = l.NextToken()
	assert.Equal(t, lexer.TokenEof, token.Type)

	token = l.NextToken()
	assert.Equal(t, lexer.TokenUndefined, token.Type)
}
