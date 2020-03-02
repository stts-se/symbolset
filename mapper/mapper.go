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
	res, err := m.SymbolSet1.ConvertToIPA(input)
	if err != nil {
		return "", fmt.Errorf("couldn't map transcription (1) : %v", err)
	}
	res, err = m.SymbolSet2.ConvertFromIPA(res)
	if err != nil {
		return "", fmt.Errorf("couldn't map transcription (2) : %v", err)
	}
	return res, nil
}

// MapSymbol maps one input transcription symbol into the new symbol set.
func (m Mapper) MapSymbol(input symbolset.Symbol) (symbolset.Symbol, error) {
	ipa := input.IPA.String
	res, err := m.SymbolSet2.GetFromIPA(ipa)
	if err != nil {
		return symbolset.Symbol{}, fmt.Errorf("couldn't map symbol : %v", err)
	}
	return res, nil
}

// MapSymbolString maps one input transcription symbol into the new symbol set.
func (m Mapper) MapSymbolString(input string) (string, error) {
	res, err := m.SymbolSet1.Get(input)
	if err != nil {
		return "", fmt.Errorf("couldn't map transcription : %v", err)
	}
	res, err = m.SymbolSet2.GetFromIPA(res.IPA.String)
	if err != nil {
		return "", fmt.Errorf("couldn't map transcription : %v", err)
	}
	return res.String, nil
}

// // MapTranscriptions maps the input entry's transcriptions (in-place)
// func (m Mapper) MapTranscriptions(e *lex.Entry) error {
// 	err := m.SymbolSet1.ConvertTranscriptionsToIPA(e)
// 	if err != nil {
// 		return fmt.Errorf("couldn't map transcription : %v", err)
// 	}
// 	err = m.SymbolSet2.ConvertTranscriptionsFromIPA(e)
// 	if err != nil {
// 		return fmt.Errorf("couldn't map transcription : %v", err)
// 	}
// 	return nil
// }
