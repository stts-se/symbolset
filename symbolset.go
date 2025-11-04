package symbolset

import (
	"fmt"
	"regexp"
	"slices"
	"strings"
)

// Type is used for accent placement, etc.
type Type int

const (
	// CMU is used for the phone set used in the CMU lexicon
	CMU Type = iota

	// SAMPA is used for SAMPA transcriptions (http://www.phon.ucl.ac.uk/home/sampa/)
	SAMPA

	// IPA is used for IPA transcriptions
	IPA

	// Other is used for symbol sets not defined in the types above
	Other
)

// SymbolCat is used to categorize transcription symbols.
type SymbolCat int

const (
	// Syllabic is used for syllabic phonemes (typically vowels and syllabic consonants)
	Syllabic SymbolCat = iota

	// NonSyllabic is used for non-syllabic phonemes (typically consonants)
	NonSyllabic

	// Stress is used for stress and accent symbols (primary, secondary, tone accents, etc)
	Stress

	// PhonemeDelimiter is used for phoneme delimiters (white space, empty string, etc)
	PhonemeDelimiter

	// SyllableDelimiter is used for syllable delimiters
	SyllableDelimiter

	// MorphemeDelimiter is used for morpheme delimiters that need not align with
	// morpheme boundaries in the decompounded orthography
	MorphemeDelimiter

	// CompoundDelimiter is used for compound delimiters that should be aligned
	// with compound boundaries in the decompounded orthography
	CompoundDelimiter

	// WordDelimiter is used for word delimiters
	WordDelimiter
)

// IPASymbol ipa symbol string with Unicode representation
type IPASymbol struct {
	String  string
	Unicode string
}

// Symbol represent a phoneme, stress or delimiter symbol used in transcriptions, including the IPA symbol with unicode
type Symbol struct {
	String string
	Cat    SymbolCat
	Desc   string
	IPA    IPASymbol
}

// SymbolSet is a struct for package private usage.
// To create a new 'SymbolSet' instance, use NewSymbolSet
type SymbolSet struct {
	Name    string
	Type    Type
	Symbols []Symbol

	// to check if the struct has been initialized properly
	isInit bool

	// derived values computed upon initialization

	// Phonemes: actual phonemes (syllabic non-syllabic)
	Phonemes []Symbol

	// PhoneticSymbols: Phonemes and stress
	PhoneticSymbols []Symbol

	stressSymbols []Symbol
	syllabic      []Symbol
	nonSyllabic   []Symbol

	PhonemeRe     *regexp.Regexp
	SyllabicRe    *regexp.Regexp
	NonSyllabicRe *regexp.Regexp
	SymbolRe      *regexp.Regexp

	ipaPhonemeRe     *regexp.Regexp
	ipaSyllabicRe    *regexp.Regexp
	ipaNonSyllabicRe *regexp.Regexp

	PhonemeDelimiter          Symbol
	phonemeDelimiterRe        *regexp.Regexp
	repeatedPhonemeDelimiters *regexp.Regexp
}

// ValidSymbol checks if a string is a valid symbol or not
func (ss SymbolSet) ValidSymbol(symbol string) bool {
	return contains(ss.Symbols, symbol)
}

// ValidIPASymbol checks if a string is a valid symbol or not
func (ss SymbolSet) ValidIPASymbol(symbol string) bool {
	for _, s := range ss.Symbols {
		if s.IPA.String == symbol {
			return true
		}
	}
	return false
}

// ContainsSymbols checks if a transcription contains a certain phoneme symbol
func (ss SymbolSet) ContainsSymbols(trans string, symbols []Symbol) (bool, error) {
	splitted, err := ss.SplitTranscription(trans)
	if err != nil {
		return false, err
	}
	for _, phn := range splitted {
		for _, symbol := range symbols {
			if phn == symbol.String {
				return true, nil
			}
		}
	}
	return false, nil
}

// Get searches the SymbolSet for a symbol with the given string
func (ss SymbolSet) Get(symbol string) (Symbol, error) {
	for _, s := range ss.Symbols {
		if s.String == symbol {
			return s, nil
		}
	}
	return Symbol{}, fmt.Errorf("no symbol /%s/ in symbol set", symbol)
}

// GetFromIPA searches the SymbolSet for a symbol with the given IPA symbol string
func (ss SymbolSet) GetFromIPA(ipa string) (Symbol, error) {
	for _, s := range ss.Symbols {
		if s.IPA.String == ipa {
			return s, nil
		}
	}
	return Symbol{}, fmt.Errorf("no ipa symbol /%s/ in symbol set", ipa)
}

// SplitTranscription splits the input transcription into separate symbols
func (ss SymbolSet) SplitTranscription(input string) ([]string, error) {
	if !ss.isInit {
		panic("symbolSet " + ss.Name + " has not been initialized properly!")
	}
	delim := ss.phonemeDelimiterRe
	if delim.FindStringIndex("") != nil {
		splitted, unknown, err := splitIntoPhonemes(ss.Symbols, input)
		if err != nil {
			return []string{}, err
		}
		if len(unknown) > 0 {
			ssErr := UnknownInputSymbol(unknown)
			return []string{}, ssErr
			//return []string{}, fmt.Errorf("found unknown phonemes in transcription /%s/: %w", input, unknown)
		}
		return splitted, nil
	}
	tmpRes := delim.Split(input, -1)
	res := []string{}
	// remove leading/trailing empty space
	for i, sym := range tmpRes {
		if (i == 0 || i == len(tmpRes)-1) && sym == "" {
			continue
		}
		res = append(res, sym)
	}
	return res, nil
}

// SplitIPATranscription splits the input transcription into separate symbols
func (ss SymbolSet) SplitIPATranscription(input string) ([]string, error) {
	if !ss.isInit {
		panic("symbolSet " + ss.Name + " has not been initialized properly!")
	}
	input, err := preFilter(ss, input, IPA)
	if err != nil {
		return []string{}, err
	}
	delim := ss.PhonemeDelimiter.IPA.String
	if delim == "" {
		symbols := []Symbol{}
		for _, s := range ss.Symbols {
			ipa := s
			ipa.String = ipa.IPA.String
			symbols = append(symbols, ipa)
		}
		splitted, unknown, err := splitIntoPhonemes(symbols, input)
		if err != nil {
			return []string{}, err
		}
		if len(unknown) > 0 {
			ssErr := UnknownInputSymbol(unknown)
			return []string{}, ssErr
		}
		return splitted, nil
	}
	return strings.Split(input, delim), nil
}

// ConvertToIPA maps one input transcription string into an IPA transcription
func (ss SymbolSet) ConvertToIPA(trans string) (string, error) {
	var unknownInputSymbols = []string{}
	res, err := preFilter(ss, trans, ss.Type)
	if err != nil {
		return "", err
	}
	splitted, err := ss.SplitTranscription(res)
	if err != nil {
		return "", err
	}
	var mapped = make([]string, 0)
	for _, fromS := range splitted {
		symbol, err := ss.Get(fromS)
		if err != nil {
			if !slices.Contains(unknownInputSymbols, fromS) {
				unknownInputSymbols = append(unknownInputSymbols, fromS)
			}
			continue
		}
		to := symbol.IPA.String
		if len(to) > 0 {
			mapped = append(mapped, to)
		}
	}
	if len(unknownInputSymbols) > 0 {
		ssErr := UnknownInputSymbol(unknownInputSymbols)
		return "", ssErr
	}

	res = strings.Join(mapped, ss.PhonemeDelimiter.IPA.String)

	res, err = postFilter(ss, res, IPA)
	return res, err
}

// ConvertFromIPA maps one input IPA transcription into the current symbol set
func (ss SymbolSet) ConvertFromIPA(trans string) (string, error) {
	res := trans
	splitted, err := ss.SplitIPATranscription(res)
	if err != nil {
		return "", err
	}
	var unknownInputSymbols = []string{}
	var mapped = make([]string, 0)
	for _, fromS := range splitted {
		symbol, err := ss.GetFromIPA(fromS)
		if err != nil {
			unknownInputSymbols = append(unknownInputSymbols, fromS)
			continue
			//return "", fmt.Errorf("input symbol /%s/ is undefined : %w", fromS, err)
		}
		to := symbol.String
		if len(to) > 0 {
			mapped = append(mapped, to)
		}
	}
	if len(unknownInputSymbols) > 0 {
		ssErr := UnknownInputSymbol(unknownInputSymbols)
		return "", ssErr
	}
	res = strings.Join(mapped, ss.PhonemeDelimiter.String)

	// remove repeated phoneme delimiters, if any
	res = ss.repeatedPhonemeDelimiters.ReplaceAllString(res, ss.PhonemeDelimiter.String)
	res, err = postFilter(ss, res, ss.Type)
	return res, err
}
