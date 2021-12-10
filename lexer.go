package lexer

import (
	"context"
	"fmt"
	"strings"
	"unicode/utf8"
)

const (
	_EOF     rune = 0
	_NEWLINE rune = '\n'
)

type lexer struct {
	input string

	start int
	pos   int
	width int

	tokens chan Token
}

// Create creates a new lexer. input is the string to be tokenized
func Create(input string) lexer {
	return lexer{
		input:  input,
		tokens: make(chan Token, 2),
	}
}

// Run lexes the input by executing state functions until the state is nil
func (l *lexer) Run(parentCtx context.Context) {
	go func() {
		defer l.finishedRun() // no more new tokens will be delivered upon exit
		for state := lexText; state != nil; {
			select {
			case <-parentCtx.Done():
				return
			default:
				state = state(l)
			}
		}
	}()
}

// finishedRun closes the chanel and marks the lexer as done.
func (l *lexer) finishedRun() {
	close(l.tokens)
}

// NextToken returns the next token, and indicates when it is done.
func (l *lexer) NextToken() Token {
	return <-l.tokens
}

// backup steps back one rune and can be called only once per call of next
func (l *lexer) backup() {
	l.pos -= l.width
}

// peek returns but does not consume the next rune in the input
func (l *lexer) peek() (nextRune rune) {
	nextRune = l.next()
	l.backup()
	return
}

// next returns the next rune in the input
func (l *lexer) next() (nextRune rune) {
	if l.pos >= len(l.input) {
		l.width = 0
		return _EOF
	}
	nextRune, l.width = utf8.DecodeRuneInString(l.input[l.pos:])
	l.pos += l.width
	return
}

// ignore skips over the pending input before this point
func (l *lexer) ignore() {
	l.start = l.pos
}

// emit puts a token onto the token channel
func (l *lexer) emit(tokenType TokenType) {
	l.tokens <- Token{Type: tokenType, Value: l.input[l.start:l.pos]}
	l.start = l.pos
}

// error returns an error token and terminates the scan
// by passing back a nil pointer that will be the next
// state, terminating Lexer.Run
func (l *lexer) errorf(format string, args ...interface{}) stateFn {
	l.tokens <- Token{
		Type:  TokenError,
		Value: fmt.Sprintf(format, args...),
	}

	return nil
}

// accept consumes the next rune if it's from the valid set.
// The valid set should always be very small 0 < len(valid) < 5
func (l *lexer) accept(valid string) bool {
	if strings.ContainsRune(valid, l.next()) {
		return true
	}
	l.backup()
	return false
}

// These values are used for acceptRun to determine the type of character that should be accepted.
const (
	_Numbers = iota
	_Hex
	_Letters
)

// acceptRun consumes a run of runes from the valid set
func (l *lexer) acceptRun(acceptType int) {
	type checkFunc func(rune) bool
	var acceptValidCharacter checkFunc

	switch acceptType {
	case _Numbers:
		acceptValidCharacter = isNumber
	case _Hex:
		acceptValidCharacter = isHex
	case _Letters:
		acceptValidCharacter = isLetter
	default:
		panic("Invalid acceptType detected.")
	}

	for acceptValidCharacter(l.next()) {
		// accept characters according to type provided
	}

	l.backup()
}
