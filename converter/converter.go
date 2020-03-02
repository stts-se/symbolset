package converter

import (
	"fmt"
	"strings"

	"github.com/dlclark/regexp2"
	"github.com/stts-se/symbolset"
)

// Converter is used to convert between symbol sets from different languages.
type Converter struct {
	Name  string
	From  symbolset.SymbolSet
	To    symbolset.SymbolSet
	Rules []Rule
}

// Convert : converts the input transcription string
func (c Converter) Convert(trans string) (string, error) {
	var res = trans
	var err error
	for _, r := range c.Rules {
		res, err = r.Convert(res, c.From)
		if err != nil {
			return "", err
		}
	}
	invalid, err := c.getInvalidSymbols(res, c.To)
	if err != nil {
		return "", err
	}
	if len(invalid) > 0 {
		return res, fmt.Errorf("invalid symbol(s) in output transcription /%s/: %v", res, invalid)
	}
	return res, nil
}

type test struct {
	from string
	to   string
}

// TestResult a test result container
type TestResult struct {
	OK     bool
	Errors []string
}

func (c Converter) getInvalidSymbols(trans string, symbolset symbolset.SymbolSet) ([]string, error) {
	if trans == symbolset.PhonemeDelimiter.String {
		return []string{}, nil
	}
	invalid := []string{}
	splitted, err := symbolset.SplitTranscription(trans)
	if err != nil {
		return invalid, err
	}
	for _, phn := range splitted {
		if !symbolset.ValidSymbol(phn) {
			invalid = append(invalid, phn)
		}
	}
	return invalid, nil
}

// Rule is a rule interface for transcription converters
type Rule interface {

	// FromString returns a string representation of the rule's input field
	FromString() string

	// ToString returns a string representation of the rule's output field
	ToString() string

	// Type returns the rule type (SYMBOL or RE)
	Type() string

	// Convert is used to execute the conversion for this rule
	Convert(trans string, symbolset symbolset.SymbolSet) (string, error)

	// String returns a tab separated string representation of the rule
	String() string
}

// SymbolRule is a simple rule that maps from one phoneme symbol to another
type SymbolRule struct {
	From string
	To   string
}

// String returns a tab separated string representation of the rule
func (r SymbolRule) String() string {
	return fmt.Sprintf("%s\t%s\t%s", "SYMBOL", r.From, r.To)
}

// FromString returns a string representation of the rule's input field
func (r SymbolRule) FromString() string {
	return r.From
}

// ToString returns a string representation of the rule's output field
func (r SymbolRule) ToString() string {
	return r.To
}

// Type returns the rule type (SYMBOL or RE)
func (r SymbolRule) Type() string {
	return "SYMBOL"
}

// Convert is used to execute the conversion for this rule
func (r SymbolRule) Convert(trans string, symbolset symbolset.SymbolSet) (string, error) {
	splitted, err := symbolset.SplitTranscription(trans)
	if err != nil {
		return "", err
	}
	res := []string{}
	for _, phn := range splitted {
		if phn == r.From {
			res = append(res, r.To)
		} else {
			res = append(res, phn)
		}

	}
	return strings.Join(res, symbolset.PhonemeDelimiter.String), nil
}

// RegexpRule is used to convert from one symbol set to another using regular expressions
type RegexpRule struct {
	From *regexp2.Regexp
	To   string
}

// String returns a tab separated string representation of the rule
func (r RegexpRule) String() string {
	return fmt.Sprintf("%s\t%s\t%s", "RE", r.From, r.To)
}

// FromString returns a string representation of the rule's input field
func (r RegexpRule) FromString() string {
	return r.From.String()
}

// ToString returns a string representation of the rule's output field
func (r RegexpRule) ToString() string {
	return r.To
}

// Type returns the rule type (SYMBOL or RE)
func (r RegexpRule) Type() string {
	return "RE"
}

// Convert is used to execute the conversion for this rule
func (r RegexpRule) Convert(trans string, symbolset symbolset.SymbolSet) (string, error) {
	res, err := r.From.Replace(trans, r.To, -1, -1)
	if err != nil {
		return "", err
	}
	return res, nil
}
