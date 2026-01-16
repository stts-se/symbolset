package mapper

import (
	"fmt"

	"github.com/stts-se/symbolset"
)

// Mapper is a struct for package private usage. To create a new instance of Mapper, use LoadMapper.
type Mapper struct {
	Name       string
	SymbolSet1 symbolset.SymbolSet
	SymbolSet2 symbolset.SymbolSet
}

// MapTranscription maps one input transcription string into the new symbol set.
func (m Mapper) MapTranscription(input string) (string, error) {
	res, err := m.SymbolSet1.ConvertToInternalIPA(input)
	if err != nil {
		return "", fmt.Errorf("couldn't map transcription (1) : %w", err)
	}
	res, err = m.SymbolSet2.ConvertFromInternalIPA(res)
	if err != nil {
		return "", fmt.Errorf("couldn't map transcription (2) : %w", err)
	}
	return res, nil
}

// MapSymbol maps one input transcription symbol into the new symbol set.
func (m Mapper) MapSymbol(input symbolset.Symbol) (symbolset.Symbol, error) {
	ipa := input.IPA.String
	res, err := m.SymbolSet2.GetFromInternalIPA(ipa)
	if err != nil {
		return symbolset.Symbol{}, fmt.Errorf("couldn't map symbol : %w", err)
	}
	return res, nil
}

// MapSymbolString maps one input transcription symbol into the new symbol set.
func (m Mapper) MapSymbolString(input string) (string, error) {
	res, err := m.SymbolSet1.Get(input)
	if err != nil {
		return "", fmt.Errorf("couldn't map transcription : %w", err)
	}
	res, err = m.SymbolSet2.GetFromInternalIPA(res.IPA.String)
	if err != nil {
		return "", fmt.Errorf("couldn't map transcription : %w", err)
	}
	return res.String, nil
}

// MapTranscriptions maps the input transcriptions
func (m Mapper) MapTranscriptions(input []string) ([]string, error) {
	var res []string
	for _, t := range input {
		tNew, err := m.MapTranscription(t)
		if err != nil {
			return res, fmt.Errorf("couldn't map transcription : %w", err)
		}
		res = append(res, tNew)
	}
	return res, nil
}
