package main

import (
	"strconv"
	"strings"
)

const (
	ParsingErrorInvalidToken ParsingErrorKind = iota
	ParsingErrorInvalidANSICodePage
)

var ()

const (
	defaultTextBufferCap = 4
)

type (
	Parser struct {
		opt      ParsingOptions
		ops      []Entity
		lexer    Lexer
		previous Token
		current  Token
	}

	ParsingErrorKind int

	ParsingError struct {
		kind  ParsingErrorKind
		token Token
		msg   string
	}

	ParsingOptions struct {
	}
)

func (e ParsingError) Error() string {
	return e.msg
}

func (parser *Parser) peek() Token {
	idx := parser.lexer.current
	token := parser.lexer.NextToken()
	parser.lexer.current = idx

	return token
}

func (parser *Parser) consume() Token {
	parser.previous = parser.current
	parser.current = parser.lexer.NextToken()
	return parser.current
}

func (parser *Parser) expectNext(k TokenKind) error {
	token := parser.consume()

	if token.kind != k {
		return ParsingError{
			kind:  ParsingErrorInvalidToken,
			token: token,
			msg:   "Invalid Token Found",
		}
	}

	return nil
}

func (parser *Parser) isWordCharacterSet() CharacterSetKind {
	token := parser.current
	if setKind, exist := characterSetKindLookup[token.text]; exist {
		return setKind
	} else if strings.HasPrefix(token.text, "ansicpg") {
		return CharacterSetANSICPG
	}

	return CharacterSetInvalid
}

func Parse(input string) ([]Entity, error) {
	parser := Parser{
		lexer: makeLexer(input),
	}

parseDocument:
	for {
		token := parser.consume()

		switch token.kind {
		case TokenEOF:
			break parseDocument

		case TokenOpenBracket:
			fallthrough
		case TokenCloseBracket:
			parser.ops = append(parser.ops, parser.parseControlGroup())

		case TokenBackslash:
			word, err := parser.parseControlWord()
			if err != nil {
				return []Entity{}, err
			}

			parser.ops = append(parser.ops, word)

		case TokenString:
			text, err := parser.parseText()
			if err != nil {
				return []Entity{}, err
			}
			parser.ops = append(parser.ops, text)

		default:
		}
	}

	return parser.ops, nil
}

func (parser *Parser) parseControlGroup() ControlGroup {
	group := ControlGroup{
		token: parser.current,
	}

	if parser.current.kind == TokenOpenBracket {
		group.groupKind = ControlGroupKindBegin
	} else {
		group.groupKind = ControlGroupKindEnd
	}

	return group
}

func (parser *Parser) parseControlWord() (Entity, error) {
	word := ControlWord{
		token: parser.current,
	}

	err := parser.expectNext(TokenString)

	if err != nil {
		return ControlWord{}, err
	}

	word.wordToken = parser.current

	if s := parser.isWordCharacterSet(); s != CharacterSetInvalid {
		return parser.parseCharacterSet(word, s)
	}

	return word, nil
}

func (parser *Parser) parseCharacterSet(word ControlWord, s CharacterSetKind) (CharacterSet, error) {
	set := CharacterSet{
		ControlWord: word,
		setKind:     s,
	}

	if s == CharacterSetANSICPG {
		codePageStr, _ := strings.CutPrefix(set.wordToken.text, "ansicpg")
		codePage, err := strconv.Atoi(codePageStr)

		if err != nil {
			return CharacterSet{}, ParsingError{
				kind:  ParsingErrorInvalidANSICodePage,
				token: set.wordToken,
				msg:   err.Error(),
			}
		}

		set.codePage = codePage
	}

	return set, nil
}

func (parser *Parser) parseText() (Text, error) {
	text := Text{
		leadingToken: parser.current,
		tokens:       make([]Token, 0, defaultTextBufferCap),
	}
	text.tokens = append(text.tokens, parser.current)

	for {
		next := parser.peek()

		if next.kind != TokenString {
			break
		}

		parser.consume()
		text.tokens = append(text.tokens, next)
	}

	return text, nil
}
