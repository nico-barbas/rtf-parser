package main

import (
	"fmt"
	"strconv"
)

const (
	ParsingErrorInvalidToken ParsingErrorKind = iota
	ParsingErrorInvalidCharacterSet
	ParsingErrorInvalidANSICodePage
	ParsingErrorInvalidNumberConversion
	ParsingErrorInvalidFormatKind
)

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

	ControlWordParsingFn func(p *Parser, c ControlWord) (Entity, error)
)

var (
	controlWordFnLookup map[string]ControlWordParsingFn
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

func (parser *Parser) expect(k TokenKind) error {
	if parser.current.kind != k {
		return ParsingError{
			kind:  ParsingErrorInvalidToken,
			token: parser.current,
			msg:   fmt.Sprintf("Expected: %s, got: %s", tokenKindStr[k], tokenKindStr[parser.current.kind]),
		}
	}

	return nil
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

func Parse(input string) ([]Entity, error) {
	parser := Parser{
		lexer: makeLexer(input),
	}

	controlWordFnLookup = map[string]ControlWordParsingFn{
		// Character set words
		"ansi":    parseCharacterSet,
		"ansicpg": parseCharacterSet,
		"mac":     parseCharacterSet,
		"pc":      parseCharacterSet,
		"pca":     parseCharacterSet,
		"fbidis":  parseCharacterSet,

		// Color words
		"colortbl": parseColorTableValue,
		"red":      parseColorComponent,
		"green":    parseColorComponent,
		"blue":     parseColorComponent,
		"alpha":    parseColorComponent,

		// Text format words
		"cf": parseTextFormat,
		"f":  parseTextFormat,
		"fs": parseTextFormat,
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

	if fn, exist := controlWordFnLookup[parser.current.text]; exist {
		return fn(parser, word)
	}

	return word, nil
}

func (parser *Parser) parseText() (Text, error) {
	text := Text{
		leadingToken: parser.current,
		tokens:       make([]Token, 0, defaultTextBufferCap),
	}
	text.tokens = append(text.tokens, parser.current)

	for {
		next := parser.peek()

		if next.kind != TokenWhitespace && next.kind != TokenNumber && next.kind != TokenString {
			break
		}

		parser.consume()
		text.tokens = append(text.tokens, next)
	}

	return text, nil
}

func parseCharacterSet(parser *Parser, word ControlWord) (Entity, error) {
	set := CharacterSet{
		ControlWord: word,
	}

	var setExist bool
	set.setKind, setExist = characterSetKindLookup[set.wordToken.text]

	if !setExist {
		return CharacterSet{}, ParsingError{
			kind:  ParsingErrorInvalidCharacterSet,
			token: set.wordToken,
		}
	}

	if set.setKind == CharacterSetANSICPG {
		err := parser.expectNext(TokenNumber)
		if err != nil {
			return CharacterSet{}, err
		}

		codePage, err := strconv.Atoi(parser.current.text)

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

func parseColorTableValue(parser *Parser, word ControlWord) (Entity, error) {
	clr := ColorTableEntry{
		ControlWord: word,
		args:        make([]Entity, 0, 4),
	}

	err := parser.expectNext(TokenSemicolon)
	if err != nil {
		return ColorTableEntry{}, err
	}

parseArgs:
	for {
		nextToken := parser.consume()

		if nextToken.kind == TokenSemicolon {
			break parseArgs
		}

		err = parser.expect(TokenBackslash)
		if err != nil {
			return ColorTableEntry{}, err
		}

		arg, err := parser.parseControlWord()
		if err != nil {
			return ColorTableEntry{}, err
		}

		clr.args = append(clr.args, arg)
	}

	return clr, nil
}

func parseColorComponent(parser *Parser, word ControlWord) (Entity, error) {
	component := ColorComponent{
		ControlWord: word,
	}

	err := parser.expectNext(TokenNumber)
	if err != nil {
		return ColorComponent{}, err
	}

	value, err := strconv.Atoi(parser.current.text)
	if err != nil {
		return ColorComponent{}, ParsingError{
			token: parser.current,
			kind:  ParsingErrorInvalidNumberConversion,
		}
	}

	component.value = uint8(value)

	return component, nil
}

func parseTextFormat(parser *Parser, word ControlWord) (Entity, error) {
	format := TextFormat{
		ControlWord: word,
	}

	formatKind, exist := textFormatKindLookup[format.wordToken.text]
	if !exist {
		return TextFormat{}, ParsingError{
			token: format.wordToken,
			kind:  ParsingErrorInvalidFormatKind,
		}
	}

	format.formatKind = formatKind

	err := parser.expectNext(TokenNumber)
	if err != nil {
		return TextFormat{}, err
	}

	value, err := strconv.Atoi(parser.current.text)
	if err != nil {
		return TextFormat{}, ParsingError{
			token: parser.current,
			kind:  ParsingErrorInvalidNumberConversion,
		}
	}

	format.arg = value

	return format, nil
}
