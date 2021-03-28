package lexer

import (
	"strings"
)

type stateFn func(*lexer) stateFn

const (
	_LEFT_META                  string = "{{"
	_RIGHT_META                 string = "}}"
	_IDENTIFIER_VALUE_INDICATOR rune   = ':'
	_IDENTIFIER_SEPARATOR       rune   = ','
)

func isSpace(r rune) bool {
	return r == ' ' || r == '\t'
}

func isNumber(r rune) bool {
	return r >= '0' && r <= '9'
}

func isHex(r rune) bool {
	return r >= '0' && r <= '9' ||
		r >= 'a' && r <= 'f' ||
		r >= 'A' && r <= 'F'
}

func isLetter(r rune) bool {
	return r >= 'a' && r <= 'z' ||
		r >= 'A' && r <= 'Z'
}

func isIdentifierSeparator(r rune) bool {
	return r == _IDENTIFIER_SEPARATOR
}

func isIdentifierValueIndicator(r rune) bool {
	return r == _IDENTIFIER_VALUE_INDICATOR
}

// lexText is the entry point and identifies text outside meta tags
func lexText(l *lexer) stateFn {
	for {
		if strings.HasPrefix(l.input[l.pos:], _LEFT_META) {
			if l.pos > l.start {
				l.emit(TokenPlainText)
			}
			return lexLeftMeta
		}
		if l.next() == _EOF {
			break
		}
	}
	// reached EOF
	if l.pos > l.start {
		l.emit(TokenPlainText)
	}
	l.emit(TokenEof)
	return nil // stop run loop
}

func lexLeftMeta(l *lexer) stateFn {
	l.pos += len(_LEFT_META)
	l.emit(TokenLeftMeta)
	return lexInsideMeta // Now inside {{ }}
}

func lexRightMeta(l *lexer) stateFn {
	l.pos += len(_RIGHT_META)
	l.emit(TokenRightMeta)
	return lexText // now outside {{ }}
}

// lexInsideMeta is inside the defined meta tags
func lexInsideMeta(l *lexer) stateFn {
	for {
		if strings.HasPrefix(l.input[l.pos:], _RIGHT_META) {
			return lexRightMeta
		}
		switch r := l.next(); {
		case r == _EOF || r == _NEWLINE:
			return l.errorf("unclosed meta")
		case isSpace(r):
			l.ignore()
		case isIdentifierSeparator(r):
			l.ignore()
		case isLetter(r):
			l.backup()
			return lexMetaIdentifier
		default:
			return l.errorf("identifier syntax: %q", l.input[l.start:l.pos])
		}
	}
}

// lexMetaIdentifier identifies an identifier inside the metadata
func lexMetaIdentifier(l *lexer) stateFn {
	l.acceptRun(_Letters)
	l.emit(TokenMetaIdentifier)

	for {
		switch r := l.next(); {
		case r == _EOF || r == _NEWLINE:
			return l.errorf("unclosed meta")
		case isSpace(r):
			l.ignore()
		case isIdentifierSeparator(r):
			l.ignore()
			return lexInsideMeta
		case isIdentifierValueIndicator(r):
			l.ignore()
			return lexIdentifierValue
		default:
			l.backup()
			return lexInsideMeta
		}
	}
}

// lexIdentifierValue identifies an identifier value after an identifier
func lexIdentifierValue(l *lexer) stateFn {
	for {
		switch r := l.next(); {
		case r == _EOF || r == _NEWLINE:
			return l.errorf("unclosed meta")
		case isSpace(r):
			l.ignore()
		case r == '+' || r == '-' || '0' <= r && r <= '9':
			l.backup()
			return lexMetaNumberValue
		case isLetter(r):
			l.backup()
			return lexMetaTextValue
		default:
			return l.errorf("value syntax: %q", l.input[l.start:l.pos])
		}
	}
}

// lexNumber identifies a number inside the metadata
func lexMetaNumberValue(l *lexer) stateFn {
	l.accept("+-")
	digits := _Numbers
	if l.accept("0") && l.accept("xX") {
		digits = _Hex
	}
	l.acceptRun(digits)
	if l.accept(".") {
		l.acceptRun(digits)
	}
	// the next rune must not be a letter
	if isLetter(l.peek()) {
		l.next()
		return l.errorf("number syntax: %q", l.input[l.start:l.pos])
	}
	l.emit(TokenMetaNumberValue)
	return lexInsideMeta
}

func lexMetaTextValue(l *lexer) stateFn {
	l.acceptRun(_Letters)
	l.emit(TokenMetaTextValue)
	return lexInsideMeta
}
