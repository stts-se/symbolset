package symbolset

import (
	"fmt"
	"strings"
)

func isTestLine(l string) bool {
	return strings.HasPrefix(l, "TEST	")
}

/// SYMBOL SET TESTS

type ssTest struct {
	testType   string
	symbolType string
	trans      string
}

func parseSSTestLine(l string) (ssTest, error) {
	fs := strings.Split(l, "\t")
	if fs[0] != "TEST" {
		return ssTest{}, fmt.Errorf("symbol set test line must start with TEST; found %s", l)
	}
	if len(fs) != 4 {
		return ssTest{}, fmt.Errorf("mapper test line must have 4 fields, found %s", l)
	}
	tType := fs[1]
	symType := fs[2]
	trans := fs[3]
	if tType != "ACCEPT" && tType != "REJECT" {
		return ssTest{}, fmt.Errorf("invalid test type %s for test line %s", tType, l)
	}
	if symType != "IPA" && symType != "SYMBOLS" {
		return ssTest{}, fmt.Errorf("invalid symbol type %s for test line %s", symType, l)
	}
	return ssTest{
		testType:   tType,
		symbolType: symType,
		trans:      trans}, nil
}

func validateTranscription(ss SymbolSet, trans string) ([]string, error) {
	var messages = make([]string, 0)
	splitted, err := ss.SplitTranscription(trans)
	if err != nil {
		if strings.Contains(fmt.Sprint(err), "unknown phonemes") {
			messages = append(messages, fmt.Sprint(err))
		} else {
			return messages, err
		}
	}
	for _, symbol := range splitted {
		if !ss.ValidSymbol(symbol) {
			messages = append(
				messages,
				fmt.Sprintf("Invalid transcription symbol '%s' in /%s/", symbol, trans))
		}
	}
	return messages, nil
}
func validateIPATranscription(ss SymbolSet, trans string) ([]string, error) {
	var messages = make([]string, 0)
	splitted, err := ss.SplitIPATranscription(trans)
	if err != nil {
		if strings.Contains(fmt.Sprint(err), "unknown phonemes") {
			messages = append(messages, fmt.Sprint(err))
		} else {
			return messages, err
		}
	}
	for _, symbol := range splitted {
		if !ss.ValidIPASymbol(symbol) {
			messages = append(
				messages,
				fmt.Sprintf("Invalid transcription symbol '%s' in /%s/", symbol, trans))
		}
	}
	return messages, nil
}

type testResult struct {
	ok     bool
	errors []string
}

func testSymbolSet(ss SymbolSet, tests []string) (testResult, error) {
	for _, test := range tests {
		t, err := parseSSTestLine(test)
		if err != nil {
			return testResult{}, err
		}
		if t.symbolType == "IPA" {
			res, err := validateIPATranscription(ss, t.trans)
			if err != nil {
				return testResult{ok: false}, err
			}
			if t.testType == "ACCEPT" && len(res) > 0 {
				return testResult{ok: false, errors: []string{fmt.Sprintf("accept test failed: /%s/ : %v", t.trans, res)}}, nil
			}
			if t.testType == "REJECT" && len(res) == 0 {
				return testResult{ok: false, errors: []string{fmt.Sprintf("reject test failed: /%s/ : %v", t.trans, res)}}, nil
			}
		} else if t.symbolType == "SYMBOLS" {
			res, err := validateTranscription(ss, t.trans)
			if err != nil {
				return testResult{ok: false}, err
			}
			if t.testType == "ACCEPT" && len(res) > 0 {
				return testResult{ok: false, errors: []string{fmt.Sprintf("accept test failed: /%s/ : %v", t.trans, res)}}, nil
			}
			if t.testType == "REJECT" && len(res) == 0 {
				return testResult{ok: false, errors: []string{fmt.Sprintf("reject test failed: /%s/ : %v", t.trans, res)}}, nil
			}
		}
	}
	return testResult{ok: true}, nil
}
