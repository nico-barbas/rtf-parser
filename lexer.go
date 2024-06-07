package main

// FIXME(nb): In some cases whitespaces are significant (in the text)
// The lexer should handle that

const (
	TokenInvalid TokenKind = iota
	TokenNewline
	TokenEOF

	TokenOpenBracket
	TokenCloseBracket
	TokenBackslash
	TokenSemicolon

	TokenString
	TokenNumber
	TokenWhitespace
)

type (
	Lexer struct {
		input   []byte
		current int
	}

	TokenKind int

	Token struct {
		kind  TokenKind
		text  string
		start int
		end   int
	}
)

func makeLexer(input string) Lexer {
	lexer := Lexer{
		input:   []byte(input),
		current: 0,
	}

	return lexer
}

func (lexer *Lexer) NextToken() Token {
	result := Token{
		start: lexer.current,
	}

	// lexer.skipWhitespace()
	if lexer.isEOF() {
		result.kind = TokenEOF
		result.end = lexer.current
		return result
	}

	c := lexer.advance()

	switch c {
	case '\n':
		result.kind = TokenNewline
	case '\\':
		result.kind = TokenBackslash
	case '{':
		result.kind = TokenOpenBracket
	case '}':
		result.kind = TokenCloseBracket
	case ';':
		result.kind = TokenSemicolon
	case ' ':
	lexWhitespace:
		for {
			if lexer.isEOF() {
				break lexWhitespace
			}

			if !isWhitespace(lexer.peek()) {
				break lexWhitespace
			}

			lexer.advance()
		}
		result.kind = TokenWhitespace
	default:
		if isLetter(c) {
		lexString:
			for {
				if lexer.isEOF() {
					break lexString
				}

				next := lexer.peek()
				if !isLetter(next) {
					break lexString
				}

				lexer.advance()
			}
			result.kind = TokenString
		} else if isNumber(c) {
		lexNumber:
			for {
				if lexer.isEOF() {
					break lexNumber
				}

				next := lexer.peek()
				if !isNumber(next) {
					break lexNumber
				}

				lexer.advance()
			}
			result.kind = TokenNumber

		} else {
			result.kind = TokenInvalid
		}
	}

	result.end = lexer.current
	result.text = string(lexer.input)[result.start:result.end]

	return result
}

func (lexer *Lexer) skipWhitespace() {
	for {
		if lexer.isEOF() {
			return
		}

		c := lexer.peek()
		if !(c == '\r' || c == '\t' || c == '\b' || c == ' ') {
			break
		}
		lexer.advance()
	}
}

func (lexer *Lexer) isEOF() bool {
	return lexer.current >= len(lexer.input)
}

func (lexer *Lexer) advance() byte {
	c := lexer.input[lexer.current]
	lexer.current += 1
	return c
}

func (lexer *Lexer) peek() byte {
	return lexer.input[lexer.current]
}

func isLetter(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
}

func isNumber(c byte) bool {
	return c >= '0' && c <= '9'
}

func isWhitespace(c byte) bool {
	return c == '\t' || c == ' '
}
