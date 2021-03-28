package lexer

type TokenType int

const (
	TokenUndefined TokenType = iota
	TokenPlainText
	TokenLeftMeta
	TokenMetaIdentifier
	TokenMetaNumberValue
	TokenMetaTextValue
	TokenRightMeta
	TokenError
	TokenEof
)

type Token struct {
	Type  TokenType
	Value string
}

func (t Token) String() string {
	if t.Type == TokenEof {
		return "EOF"
	}
	return t.Value
}

func (t TokenType) String() string {
	switch t {
	case TokenUndefined:
		return "Undefined"
	case TokenPlainText:
		return "PlainText"
	case TokenLeftMeta:
		return "LeftMeta"
	case TokenMetaIdentifier:
		return "MetaIdentifier"
	case TokenMetaNumberValue:
		return "MetaNumberValue"
	case TokenMetaTextValue:
		return "MetaTextValue"
	case TokenRightMeta:
		return "RightMeta"
	case TokenError:
		return "Error"
	case TokenEof:
		return "Eof"
	}
	return "invalid"
}
