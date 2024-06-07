package main

const (
	EntityKindInvalid EntityKind = iota
	EntityKindControlGroup
	EntityKindControlWord
	EntityKindCharacterSet
	EntityKindText
)

var (
	entityKindStr = map[EntityKind]string{
		EntityKindInvalid:      "Invalid",
		EntityKindControlGroup: "Control Group",
		EntityKindControlWord:  "Control Word",
		EntityKindCharacterSet: "Character Set",
		EntityKindText:         "Text",
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
		"ansi":   CharacterSetANSI,
		"mac":    CharacterSetMAC,
		"pc":     CharacterSetPC,
		"pca":    CharacterSetPCA,
		"fbidis": CharacterSetFBIDIS,
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

type (
	EntityKind int

	Entity interface {
		kind() EntityKind
		getToken() Token
	}

	ControlGroupKind uint8
	CharacterSetKind uint8

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

func (t Text) kind() EntityKind {
	return EntityKindText
}

func (t Text) getToken() Token {
	return t.leadingToken
}
