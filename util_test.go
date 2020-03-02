package symbolset

import "testing"

var fsExp = "Expected: '%v' got: '%v'"

func testEqStrings(t *testing.T, expect []string, result []string) {
	if len(expect) != len(result) {
		t.Errorf(fsExp, expect, result)
		return
	}
	for i, ex := range expect {
		re := result[i]
		if ex != re {
			t.Errorf(fsExp, expect, result)
			return
		}
	}
}

func testEqSymbols(t *testing.T, expect []Symbol, result []Symbol) {
	if len(expect) != len(result) {
		t.Errorf(fsExp, expect, result)
		return
	}
	for i, ex := range expect {
		re := result[i]
		if ex != re {
			t.Errorf(fsExp, expect, result)
			return
		}
	}
}

func Test_buildRegexp1(t *testing.T) {
	symbols := []Symbol{
		{"a", Syllabic, "", IPASymbol{"", ""}},
		{"t", NonSyllabic, "", IPASymbol{"", ""}},
		{".", SyllableDelimiter, "", IPASymbol{"", ""}},
		{"s", NonSyllabic, "", IPASymbol{"", ""}},
		{"t_s", NonSyllabic, "", IPASymbol{"", ""}},
		{"", PhonemeDelimiter, "", IPASymbol{"", ""}},
		{" ", PhonemeDelimiter, "", IPASymbol{"", ""}},
		{"e", Syllabic, "", IPASymbol{"", ""}},
	}
	re, err := buildRegexp(symbols)
	if err != nil {
		t.Errorf("buildRegexp() didn't expect error here : %v", err)
	}
	expect := "(?:t_s| |\\.|a|e|s|t|)"
	if expect != re.String() {
		t.Errorf(fsExp, expect, re.String())
	}
}

func Test_buildRegexp2(t *testing.T) {
	symbols := []Symbol{
		{"a", Syllabic, "", IPASymbol{"", ""}},
		{"t", NonSyllabic, "", IPASymbol{"", ""}},
		{".", SyllableDelimiter, "", IPASymbol{"", ""}},
		{"s", NonSyllabic, "", IPASymbol{"", ""}},
		{"t_s", NonSyllabic, "", IPASymbol{"", ""}},
		{"", PhonemeDelimiter, "", IPASymbol{"", ""}},
		{" ", PhonemeDelimiter, "", IPASymbol{"", ""}},
		{"e", Syllabic, "", IPASymbol{"", ""}},
	}
	re, err := buildRegexpWithGroup(symbols, true, false)
	if err != nil {
		t.Errorf("buildRegexp() didn't expect error here : %v", err)
	}
	expect := "(t_s| |\\.|a|e|s|t)"
	if expect != re.String() {
		t.Errorf(fsExp, expect, re.String())
	}
}

func Test_buildRegexp3(t *testing.T) {
	symbols := []Symbol{
		{"a", Syllabic, "", IPASymbol{"", ""}},
		{"t", NonSyllabic, "", IPASymbol{"", ""}},
		{".", SyllableDelimiter, "", IPASymbol{"", ""}},
		{"s", NonSyllabic, "", IPASymbol{"", ""}},
		{"t_s", NonSyllabic, "", IPASymbol{"", ""}},
		{"$", PhonemeDelimiter, "", IPASymbol{"", ""}},
		{" ", PhonemeDelimiter, "", IPASymbol{"", ""}},
		{"", PhonemeDelimiter, "", IPASymbol{"", ""}},
		{"e", Syllabic, "", IPASymbol{"", ""}},
	}
	re, err := buildRegexpWithGroup(symbols, false, false)
	if err != nil {
		t.Errorf("buildRegexp() didn't expect error here : %v", err)
	}
	expect := "(t_s| |\\$|\\.|a|e|s|t|)"
	if expect != re.String() {
		t.Errorf(fsExp, expect, re.String())
	}
}

func Test_FilterSymbolsByCat(t *testing.T) {
	symbols := []Symbol{
		{"a", Syllabic, "", IPASymbol{"", ""}},
		{"t", NonSyllabic, "", IPASymbol{"", ""}},
		{"%", Stress, "", IPASymbol{"", ""}},
		{".", SyllableDelimiter, "", IPASymbol{"", ""}},
		{"s", NonSyllabic, "", IPASymbol{"", ""}},
		{"t_s", NonSyllabic, "", IPASymbol{"", ""}},
		{"$", PhonemeDelimiter, "", IPASymbol{"", ""}},
		{"\"", Stress, "", IPASymbol{"", ""}},
		{"e", Syllabic, "", IPASymbol{"", ""}},
		{"", PhonemeDelimiter, "", IPASymbol{"", ""}},
		{"+", MorphemeDelimiter, "", IPASymbol{"", ""}},
	}
	stressE := []Symbol{
		{"%", Stress, "", IPASymbol{"", ""}},
		{"\"", Stress, "", IPASymbol{"", ""}},
	}
	stressR := filterSymbolsByCat(symbols, []SymbolCat{Stress})
	testEqSymbols(t, stressE, stressR)

	delimE := []Symbol{
		{".", SyllableDelimiter, "", IPASymbol{"", ""}},
		{"$", PhonemeDelimiter, "", IPASymbol{"", ""}},
		{"", PhonemeDelimiter, "", IPASymbol{"", ""}},
		{"+", MorphemeDelimiter, "", IPASymbol{"", ""}},
	}
	delimR := filterSymbolsByCat(symbols, []SymbolCat{SyllableDelimiter, PhonemeDelimiter, MorphemeDelimiter})
	testEqSymbols(t, delimE, delimR)
}

func Test_contains(t *testing.T) {
	symbols := []Symbol{
		{"a", Syllabic, "", IPASymbol{"", ""}},
		{"t", NonSyllabic, "", IPASymbol{"", ""}},
		{".", SyllableDelimiter, "", IPASymbol{"", ""}},
		{"s", NonSyllabic, "", IPASymbol{"", ""}},
		{"t_s", NonSyllabic, "", IPASymbol{"", ""}},
		{"$", PhonemeDelimiter, "", IPASymbol{"", ""}},
		{" ", PhonemeDelimiter, "", IPASymbol{"", ""}},
		{"", PhonemeDelimiter, "", IPASymbol{"", ""}},
		{"e", Syllabic, "", IPASymbol{"", ""}},
	}
	var s string

	s = "t_s"
	if !contains(symbols, s) {
		t.Errorf("contains() Expected true for symbol %s in %v", s, symbols)
	}
	s = "_"
	if contains(symbols, s) {
		t.Errorf("contains() Expected false for symbol %s in %v", s, symbols)
	}
}

func Test_string2unicode(t *testing.T) {

	// --
	s := "a"
	expect := "U+0061"
	result := string2unicode(s)
	if result != expect {
		t.Errorf("For /%s/, expected '%s', got '%s'", s, expect, result)
	}

	// --
	s = "_"
	expect = "U+005F"
	result = string2unicode(s)
	if result != expect {
		t.Errorf("For /%s/, expected '%s', got '%s'", s, expect, result)
	}
}
