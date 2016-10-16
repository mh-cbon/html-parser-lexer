package htmlparserlexer

import (
	// "fmt"
	// "time"
	"github.com/mh-cbon/state-lexer"
	"io"
)

const (
	TagOpenToken  lexer.TokenType = iota // 0
	TagCloseToken                        // 1

	TagOpenStartToken  // 2
	TagOpenEndToken    // 3
	TagCloseStartToken // 4
	TagCloseEndToken   // 5

	TagAttrNameToken  // 6
	TagAttrEqToken    // 7
	TagAttrQuoteToken // 8
	TagAttrValueToken // 9

	WsToken   // 10
	TextToken // 11

	CommentStartToken // 12
	CommentToken      // 13
	CommentEndToken   // 14
)

const (
	NOTEOFRune rune = -2
)

func NewHtmlParserLexer(r io.Reader) *lexer.L {
	return lexer.New(r, TextState)
}

func isWhitespace(ch rune) bool { return ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r' }
func isQuote(ch rune) bool      { return ch == '"' || ch == '\'' }

func eatSomeWs(l *lexer.L) rune {
	hasWs := false
	r := l.Next()
	for isWhitespace(r) {
		hasWs = true
		if r == lexer.EOFRune {
			l.Emit(WsToken)
			return lexer.EOFRune
		}
		r = l.Next()
	}
	l.Rewind()
	if hasWs {
		l.Emit(WsToken)
	}
	return NOTEOFRune
}

func isATagCloseStart(l *lexer.L) bool {
	r := l.Next()
	if r != '<' {
		l.Rewind()
		return false
	}
	r = l.Next()
	if r != '/' {
		l.Rewind()
		l.Rewind()
		return false
	}
	l.Rewind()
	l.Rewind()
	return true
}

func isATagOpenStart(l *lexer.L) bool {
	r := l.Next()
	if r != '<' {
		l.Rewind()
		return false
	}
	r = l.Next()
	if (r < 'a' || r > 'z') && (r < 'A' || r > 'Z') {
		l.Rewind()
		l.Rewind()
		return false
	}
	l.Rewind()
	l.Rewind()
	return true
}

func isATagCommentStartToken(l *lexer.L) bool {
	r := l.Next()
	if r != '<' {
		l.Rewind()
		return false
	}
	r = l.Next()
	if r != '!' {
		l.Rewind()
		l.Rewind()
		return false
	}
	r = l.Next()
	if r != '-' {
		l.Rewind()
		l.Rewind()
		l.Rewind()
		return false
	}
	r = l.Next()
	if r != '-' {
		l.Rewind()
		return false
	}
	l.Rewind()
	l.Rewind()
	l.Rewind()
	l.Rewind()
	return true
}

func isATagCommentEndToken(l *lexer.L) bool {

	if r := l.Next(); r != '-' {
		l.Rewind()
		return false
	}
	if r := l.Next(); r != '-' {
		l.Rewind()
		l.Rewind()
		return false
	}
	if r := l.Next(); r != '>' {
		l.Rewind()
		l.Rewind()
		l.Rewind()
		return false
	}
	l.Rewind()
	l.Rewind()
	l.Rewind()
	return true
}

func TextState(l *lexer.L) lexer.StateFunc {

	// eat some ws
	if r := eatSomeWs(l); r == lexer.EOFRune {
		return nil
	}

	hasText := false
	r := l.Next()
	for r != '<' {
		if r == lexer.EOFRune {
			if hasText {
				l.Emit(TextToken)
			}
			return nil
		}
		hasText = true
		r = l.Next()
	}
	l.Rewind()
	if hasText {
		l.Emit(TextToken)
	}

	// is it goind to be a </
	if isATagCloseStart(l) {
		return TagCloseState
	}
	// or a <
	if isATagOpenStart(l) {
		return TagStartState
	}
	// or a <!--
	if isATagCommentStartToken(l) {
		return TagCommentStartState
	}
	// dunno... :s :s :s yet!
	return TextState
}

func TagCommentStartState(l *lexer.L) lexer.StateFunc {
	// eat <
	if r := l.Next(); r != '<' {
		panic("ohoh not a comment start")
	}
	// eat !
	if r := l.Next(); r != '!' {
		panic("ohoh not a comment start")
	}
	// eat -
	if r := l.Next(); r != '-' {
		panic("ohoh not a comment start")
	}
	// eat -
	if r := l.Next(); r != '-' {
		panic("ohoh not a comment start")
	}
	l.Emit(CommentStartToken)

	for { // @todo do something if the comment is unclosed. possible ?
		if r := l.Next(); r == lexer.EOFRune {
			return nil
		}
		// <-time.After(time.Millisecond)
		if isATagCommentEndToken(l) {
			break
		}
	}
	l.Emit(CommentToken)

	// eat -
	if r := l.Next(); r != '-' {
		panic("ohoh not a comment end")
	}
	// eat -
	if r := l.Next(); r != '-' {
		panic("ohoh not a comment end")
	}
	// eat >
	if r := l.Next(); r != '>' {
		panic("ohoh not a comment end")
	}
	l.Emit(CommentEndToken)

	return TextState
}

func TagStartState(l *lexer.L) lexer.StateFunc {
	r := l.Next()
	// eat <
	if r != '<' {
		panic("ohoh not a tag start")
	}
	l.Emit(TagOpenStartToken)

	// // eat in between ws
	if r := eatSomeWs(l); r == lexer.EOFRune {
		return nil
	}

	// eat the tag name
	r = l.Next()
	for (r >= 'a' && r <= 'z') || r == ':' {
		r = l.Next()
	}
	l.Rewind()
	l.Emit(TagOpenToken)

	// // eat in between ws
	if r := eatSomeWs(l); r == lexer.EOFRune {
		return nil
	}

	// eat >
	r = l.Next()
	if r == '>' {
		l.Emit(TagOpenEndToken)
		return TextState
	}
	l.Rewind()

	//badass html, unclosed html tag name.
	r = l.Next()
	if r == '<' {
		l.Rewind()
		l.Emit(TagOpenEndToken) //empty token!!
		return TagStartState
	}
	l.Rewind()

	// go for tag attributes
	return TagAttrNameState
}

func TagCloseState(l *lexer.L) lexer.StateFunc {
	r := l.Next()
	// eat <
	if r != '<' {
		panic("ohoh not a tag close")
	}
	r = l.Next()
	if r != '/' {
		panic("ohoh not a tag close")
	}
	l.Emit(TagCloseStartToken)

	// // eat in between ws
	if r := eatSomeWs(l); r == lexer.EOFRune {
		return nil
	}

	// eat the tag name
	r = l.Next()
	for r >= 'a' && r <= 'z' {
		r = l.Next()
	}
	l.Rewind()
	l.Emit(TagCloseToken)

	// // eat in between ws
	if r := eatSomeWs(l); r == lexer.EOFRune {
		return nil
	}

	// eat >
	r = l.Next()
	if r != '>' {
		l.Rewind()
	}
	l.Emit(TagCloseEndToken)
	return TextState
}

func TagAttrNameState(l *lexer.L) lexer.StateFunc {
	// eat attribute name
	r := l.Next()
	for r != '=' && !isWhitespace(r) {
		r = l.Next()
	}
	l.Rewind()
	l.Emit(TagAttrNameToken)

	// // eat in between ws
	if r := eatSomeWs(l); r == lexer.EOFRune {
		return nil
	}

	// eat eq sign
	withValue := false
	r = l.Next()
	if r == '=' {
		withValue = true
		l.Emit(TagAttrEqToken)
	} else {
		l.Rewind()
	}

	// // eat in between ws
	if r := eatSomeWs(l); r == lexer.EOFRune {
		return nil
	}

	if withValue {
		// eat tag value quote, if any
		r = l.Next()
		if isQuote(r) {
			quote := r // save quote type
			l.Emit(TagAttrQuoteToken)
			// eat the attribute value until a quote
			var p rune // holds previous read to test for escaped quotes
			for {
				r = l.Next()
				for !isQuote(r) {
					p = r
					r = l.Next()
				}
				if r == quote && p != rune('\\') {
					break
				}
			}
			l.Rewind()
			l.Emit(TagAttrValueToken)
			// eat the quote
			l.Next()
			l.Emit(TagAttrQuoteToken)

		} else {
			// eat the attribute value until a whitespace or a >
			r = l.Next()
			for r != '>' && !isWhitespace(r) {
				r = l.Next()
			}
			l.Rewind()
			l.Emit(TagAttrValueToken)
		}

		// // eat in between ws
		if r := eatSomeWs(l); r == lexer.EOFRune {
			return nil
		}

	}

	// close tag, maybe
	r = l.Next()
	if r == '/' {
		r = l.Next()
	}
	if r == '>' {
		l.Emit(TagOpenEndToken)
		return TextState
	}

	// keep looking for attributes
	return TagAttrNameState
}

// Helper function
func TokenName(tok lexer.Token) string {
	switch tok.Type {
	case TagOpenToken:
		return "TagOpenToken"
	case TagCloseToken:
		return "TagCloseToken"
	case TagOpenStartToken:
		return "TagOpenStartToken"
	case TagOpenEndToken:
		return "TagOpenEndToken"
	case TagCloseStartToken:
		return "TagCloseStartToken"
	case TagCloseEndToken:
		return "TagCloseEndToken"
	case TagAttrNameToken:
		return "TagAttrNameToken"
	case TagAttrEqToken:
		return "TagAttrEqToken"
	case TagAttrQuoteToken:
		return "TagAttrQuoteToken"
	case TagAttrValueToken:
		return "TagAttrValueToken"
	case WsToken:
		return "WsToken"
	case TextToken:
		return "TextToken"
	}
	return "undefined"
}
