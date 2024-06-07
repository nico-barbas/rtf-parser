package main

const (
	EntityKindInvalid EntityKind = iota
	EntityKindControlGroup
	EntityKindControlWord
	EntityKindCharacterSet
	EntityKindColorTableEntry
	EntityKindColorComponent
	EntityKindTextFormat
	EntityKindText
)

var (
	entityKindStr = map[EntityKind]string{
		EntityKindInvalid:         "Invalid",
		EntityKindControlGroup:    "Control Group",
		EntityKindControlWord:     "Control Word",
		EntityKindCharacterSet:    "Character Set",
		EntityKindColorTableEntry: "Color Table Entry",
		EntityKindColorComponent:  "Color Component",
		EntityKindTextFormat:      "Text Format",
		EntityKindText:            "Text",
	}
)

const (
	ControlGroupKindBegin ControlGroupKind = iota
	ControlGroupKindEnd
)

var (
	controlGroupKindStr = map[ControlGroupKind]string{
		ControlGroupKindBegin: "Begin",
		ControlGroupKindEnd:   "End",
	}
)

const (
	CharacterSetInvalid CharacterSetKind = iota
	CharacterSetANSI
	CharacterSetANSICPG
	CharacterSetMAC
	CharacterSetPC
	CharacterSetPCA
	CharacterSetFBIDIS
)

var (
	characterSetKindLookup = map[string]CharacterSetKind{
		"ansi":    CharacterSetANSI,
		"ansicpg": CharacterSetANSICPG,
		"mac":     CharacterSetMAC,
		"pc":      CharacterSetPC,
		"pca":     CharacterSetPCA,
		"fbidis":  CharacterSetFBIDIS,
	}

	characterSetKindStr = map[CharacterSetKind]string{
		CharacterSetANSI:    "ansi",
		CharacterSetANSICPG: "ansicpg",
		CharacterSetMAC:     "mac",
		CharacterSetPC:      "pc",
		CharacterSetPCA:     "pca",
		CharacterSetFBIDIS:  "fbidi",
	}
)

const (
	TextFormatColor TextFormatKind = iota
	TextFormatFontIndex
	TextFormatFontSize
)

var (
	textFormatKindLookup = map[string]TextFormatKind{
		"cf": TextFormatColor,
		"f":  TextFormatFontIndex,
		"fs": TextFormatFontSize,
	}

	textFormatKindStr = map[TextFormatKind]string{
		TextFormatColor:     "Color",
		TextFormatFontIndex: "Font",
		TextFormatFontSize:  "Font size",
	}
)

type (
	EntityKind int

	Entity interface {
		kind() EntityKind
		getToken() Token
	}

	ControlGroupKind uint8
	CharacterSetKind uint8
	TextFormatKind   uint8

	ControlGroup struct {
		token     Token
		groupKind ControlGroupKind
	}

	ControlWord struct {
		token     Token
		wordToken Token
	}

	CharacterSet struct {
		ControlWord
		setKind  CharacterSetKind
		codePage int
	}

	ColorTableEntry struct {
		ControlWord
		args []Entity
	}

	ColorComponent struct {
		ControlWord
		value uint8
	}

	TextFormat struct {
		ControlWord
		formatKind TextFormatKind
		arg        int
	}

	Text struct {
		leadingToken Token
		tokens       []Token
	}
)

func (c ControlGroup) kind() EntityKind {
	return EntityKindControlGroup
}

func (c ControlGroup) getToken() Token {
	return c.token
}

func (c ControlWord) kind() EntityKind {
	return EntityKindControlWord
}

func (c ControlWord) getToken() Token {
	return c.token
}

func (c CharacterSet) kind() EntityKind {
	return EntityKindCharacterSet
}

func (c CharacterSet) getToken() Token {
	return c.token
}

func (c ColorTableEntry) kind() EntityKind {
	return EntityKindColorTableEntry
}

func (c ColorTableEntry) getToken() Token {
	return c.token
}

func (c ColorComponent) kind() EntityKind {
	return EntityKindColorComponent
}

func (c ColorComponent) getToken() Token {
	return c.token
}

func (t Text) kind() EntityKind {
	return EntityKindText
}

func (c TextFormat) kind() EntityKind {
	return EntityKindTextFormat
}

func (c TextFormat) getToken() Token {
	return c.token
}

func (t Text) getToken() Token {
	return t.leadingToken
}
