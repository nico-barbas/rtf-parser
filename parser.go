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
		opt              ParsingOptions
		ops              []Entity
		lexer            Lexer
		previous         Token
		current          Token
		textEscapeTokens []TokenKind
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
	controlWordFnLookup     map[string]ControlWordParsingFn
	defaultTextEscapeTokens = []TokenKind{TokenOpenBracket, TokenCloseBracket, TokenBackslash}
	fontTextEscapeTokens    = []TokenKind{TokenOpenBracket, TokenCloseBracket, TokenBackslash, TokenSemicolon}
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

func (parser *Parser) peekNext() Token {
	idx := parser.lexer.current
	parser.lexer.NextToken()
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
		lexer:            makeLexer(input),
		textEscapeTokens: defaultTextEscapeTokens,
	}

	controlWordFnLookup = map[string]ControlWordParsingFn{
		// Character set words
		"ansi":    parseCharacterSet,
		"ansicpg": parseCharacterSet,
		"mac":     parseCharacterSet,
		"pc":      parseCharacterSet,
		"pca":     parseCharacterSet,
		"fbidis":  parseCharacterSet,

		// Font words
		"fonttbl": parseFontTable,

		// Color words
		"colortbl": parseColorTable,
		"red":      parseColorComponent,
		"green":    parseColorComponent,
		"blue":     parseColorComponent,
		"alpha":    parseColorComponent,

		// Text format words
		"cf":   parseTextFormat,
		"f":    parseTextFormat,
		"fs":   parseTextFormat,
		"li":   parseTextFormat,
		"fi":   parseTextFormat,
		"pard": parseTextFormatNoArg,
		"par":  parseTextFormatNoArg,
		"b":    parseTextFormatNoArg,
		"qc":   parseTextFormatNoArg,
		"qj":   parseTextFormatNoArg,
		"qr":   parseTextFormatNoArg,
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

parseSequence:
	for {
		next := parser.peek()

		for _, escape := range parser.textEscapeTokens {
			if next.kind == escape {
				break parseSequence
			}
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

func parseFontTable(parser *Parser, word ControlWord) (Entity, error) {
	tbl := FontTable{
		ControlWord: word,
	}

parseFonts:
	for {
		nextToken := parser.peek()

		if nextToken.kind != TokenOpenBracket {
			break parseFonts
		}

		parser.consume()

		f, err := parseFontTableEntry(parser)
		if err != nil {
			return FontTable{}, err
		}

		tbl.fonts = append(tbl.fonts, f)
	}

	return tbl, nil
}

func parseFontTableEntry(parser *Parser) (Entity, error) {
	fnt := FontTableEntry{}

	err := parser.expect(TokenOpenBracket)
	if err != nil {
		return FontTableEntry{}, err
	}

	parser.textEscapeTokens = fontTextEscapeTokens

parseArgs:
	for {
		nextToken := parser.consume()

		switch nextToken.kind {
		case TokenSemicolon:
			break parseArgs
		case TokenString:
			fnt.fontName, err = parser.parseText()
			if err != nil {
				return FontTableEntry{}, err
			}
			fallthrough
		case TokenWhitespace:
			continue
		}

		err = parser.expect(TokenBackslash)
		if err != nil {
			return FontTableEntry{}, err
		}

		err = parser.expectNext(TokenString)
		if err != nil {
			return FontTableEntry{}, err
		}

		switch parser.current.text {
		case "f":
			err = parser.expectNext(TokenNumber)
			if err != nil {
				return FontTableEntry{}, err
			}

			fnt.index, err = strconv.Atoi(parser.current.text)
			if err != nil {
				return FontTableEntry{}, ParsingError{
					token: parser.current,
					kind:  ParsingErrorInvalidNumberConversion,
				}
			}

		case "fnil":
			fnt.defaultFallback = true

		case "fcharset":
			err = parser.expectNext(TokenNumber)
			if err != nil {
				return FontTableEntry{}, err
			}

			fnt.charset, err = strconv.Atoi(parser.current.text)
			if err != nil {
				return FontTableEntry{}, ParsingError{
					token: parser.current,
					kind:  ParsingErrorInvalidNumberConversion,
				}
			}

		default:
			return FontTableEntry{}, ParsingError{
				token: parser.current,
				kind:  ParsingErrorInvalidToken,
			}
		}
	}

	err = parser.expectNext(TokenCloseBracket)
	if err != nil {
		return FontTableEntry{}, err
	}
	parser.textEscapeTokens = defaultTextEscapeTokens

	return fnt, nil
}

func parseColorTable(parser *Parser, word ControlWord) (Entity, error) {
	table := ColorTable{
		ControlWord: word,
	}

parseColors:
	for {
		nextToken := parser.peek()

		switch nextToken.kind {
		case TokenCloseBracket:
			break parseColors
		case TokenSemicolon:
			parser.consume()
			if parser.peek().kind != TokenCloseBracket {
				clr, err := parseColorTableEntry(parser)
				if err != nil {
					return ColorTable{}, err
				}

				table.colors = append(table.colors, clr)
			}
		}
	}

	return table, nil
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

	nextToken := parser.consume()
	negateNumber := false
	if nextToken.kind == TokenDash {
		negateNumber = true
		parser.consume()
	}

	err := parser.expect(TokenNumber)
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

	if negateNumber {
		format.arg = -format.arg
	}

	return format, nil
}

func parseTextFormatNoArg(parser *Parser, word ControlWord) (Entity, error) {
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
	format.arg = -1
	return format, nil
}

func parseColorTableEntry(parser *Parser) (Entity, error) {
	clr := ColorTableEntry{
		startToken: parser.current,
	}

parseComponents:
	for i := 0; i < 4; i += 1 {
		nextToken := parser.peek()

		if nextToken.kind == TokenSemicolon {
			break parseComponents
		}

		err := parser.expectNext(TokenBackslash)
		if err != nil {
			return ColorTable{}, err
		}

		channel, err := parser.parseControlWord()
		if err != nil {
			return ColorTable{}, err
		}

		clr.channels[i] = channel
	}

	return clr, nil
}
