package lexer

import (
	"jp/token"
	"slices"
)

type control struct {
	value int
}

func (c *control) inc() {
	c.value++
}

func (c *control) dec() {
	c.value--
	if c.value < 0 {
		c.value = 0
	}
}

func (c *control) close() {
	c.value = 0
}

func (c *control) isOpen() bool {
	return c.value > 0
}

type Lexer struct {
	input        string
	position     int  // current position in input (points to current char)
	readPosition int  // current reading position in input (after current char)
	ch           byte // current char under examination

	expectKey    bool
	expectValue  bool
	expectObject bool
	expectArray  bool
	expectColon  bool
	expectComma  bool

	openObject *control
	openArray  *control
	openComma  *control
}

func New(input string) *Lexer {
	l := Lexer{
		input: input,

		expectObject: true,
		expectArray:  true,

		openObject: new(control),
		openArray:  new(control),
		openComma:  new(control),
	}
	l.readChar()

	return &l
}

func (l *Lexer) NextToken() token.Token {
	l.skipWhitespace()
	tok := newToken(token.ILLEGAL, l.ch)

	switch l.ch {
	case ':':
		if l.expectColon && !l.openComma.isOpen() {
			tok = newToken(token.COLON, l.ch)
			l.expect(false, true, true, true, false, false)
		}

	case '[':
		if l.expectArray {
			if l.openArray.value == 19 {
				// too deep
				return tok
			}

			tok = newToken(token.LBRACKET, l.ch)
			l.expect(false, true, true, true, false, false)
			l.openArray.inc()
			l.openComma.close()
		}

	case ']':
		if l.openArray.isOpen() && !l.openComma.isOpen() {
			tok = newToken(token.RBRACKET, l.ch)
			l.openArray.dec()
			l.expectComma = l.openArray.isOpen()
		}

	case '{':
		if l.expectObject {
			tok = newToken(token.LBRACE, l.ch)
			l.expect(true, false, false, false, false, false)
			l.openObject.inc()
			l.openComma.close()
		}

	case '}':
		if l.openObject.isOpen() && !l.openComma.isOpen() {
			tok = newToken(token.RBRACE, l.ch)
			l.openObject.dec()
			l.expectKey = l.openObject.isOpen()
			l.expectComma = l.openObject.isOpen() || l.openArray.isOpen()
		}

	case ',':
		if l.expectComma {
			tok = newToken(token.COMMA, l.ch)
			l.expect(l.openObject.isOpen(), l.openArray.isOpen(), l.openArray.isOpen(), l.openArray.isOpen(), false, false)
			l.openComma.inc()
		}

	case 0:
		tok.Literal = ""
		tok.Type = token.EOF

	default:
		if l.expectKey {
			return l.readKey(tok)
		}
		if l.expectValue {
			return l.readValue(tok)
		}
	}

	l.readChar()
	return tok
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}

	l.position = l.readPosition
	l.readPosition += 1
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func (l *Lexer) expect(key, value, object, array, colon, comma bool) {
	l.expectKey = key
	l.expectValue = value
	l.expectObject = object
	l.expectArray = array
	l.expectColon = colon
	l.expectComma = comma
}

func (l *Lexer) readKey(tok token.Token) token.Token {
	if !isQuote(l.ch) {
		return tok
	}

	l.readChar()
	position := l.position

	for {
		if isQuote(l.ch) {
			break
		}

		if l.ch == 0 {
			return tok
		}

		l.readChar()
	}

	tok.Type = token.KEY
	tok.Literal = l.input[position:l.position]

	if position == l.position {
		tok.Literal = "EMPTY"
	}

	l.readChar()
	l.expect(false, false, false, false, true, false)
	l.openComma.close()

	return tok
}

func (l *Lexer) readValue(tok token.Token) token.Token {
	if isDigit(l.ch) {
		return l.getValueNumber(tok)
	}

	if isQuote(l.ch) {
		return l.getValueString(tok)
	}

	return l.getValueSpecial(tok)
}

func (l *Lexer) getValueNumber(tok token.Token) token.Token {
	position := l.position
	initialElement := l.ch
	hasMinus := false
	hasDot := false
	hasE := false
	hasEMinus := false
	hasEPlus := false

	if l.ch == 'e' || l.ch == 'E' || l.ch == '+' {
		return tok
	}

	for {
		if !isDigit(l.ch) {
			if position == l.position {
				return tok
			}

			last := l.input[l.position]
			if last == 'e' || last == '+' || last == '-' {
				return tok
			}
			if len(l.input[position:l.position]) > 1 && initialElement == '0' && (!hasE || !hasDot) {
				return tok
			}

			tok.Type = token.VALUE
			tok.Literal = l.input[position:l.position]

			l.expect(false, false, false, false, false, true)

			break
		}

		if l.ch == '-' {
			if hasMinus {
				return tok
			}

			hasMinus = true
		}

		if l.ch == '.' {
			if hasDot {
				return tok
			}

			hasDot = true
		}

		if l.ch == 'e' || l.ch == 'E' {
			if hasE {
				return tok
			}

			hasE = true
		}

		if position > l.position && l.ch == '-' {
			if hasEMinus {
				return tok
			}

			hasEMinus = true
		}

		if position > l.position && l.ch == '+' {
			if hasEPlus {
				return tok
			}

			hasEPlus = true
		}

		l.readChar()
	}

	return tok
}

func (l *Lexer) getValueString(tok token.Token) token.Token {
	l.readChar()
	position := l.position

LOOP:
	for {
		switch l.ch {
		case '"':
			break LOOP

		case '\\':
			return tok

		case '\t':
			return tok

		case '\n':
			return tok

		case '\r':
			return tok
		}

		l.readChar()
	}

	if position == l.position {
		return tok
	}

	tok.Type = token.VALUE
	tok.Literal = l.input[position:l.position]

	l.readChar()
	l.expect(false, false, false, false, false, true)

	return tok
}

func (l *Lexer) getValueSpecial(tok token.Token) token.Token {
	size := 4
	special := l.input[l.position : l.position+size]
	ok := slices.Contains([]string{"null", "true"}, special)

	if !ok {
		size = 5
		special = l.input[l.position : l.position+size]
		ok = special == "false"
	}

	if ok {
		tok.Type = token.VALUE
		tok.Literal = l.input[l.position : l.position+size]

		for i := 0; i < size; i++ {
			l.readChar()
		}
		l.expect(false, false, false, false, false, true)
	}

	return tok
}

func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{
		Type:    tokenType,
		Literal: string(ch),
	}
}

func isQuote(ch byte) bool {
	return ch == '"'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9' || ch == '.' || ch == '-' || ch == '+' || ch == 'e' || ch == 'E'
}
