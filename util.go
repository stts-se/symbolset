package symbolset

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
)

// symbolSlice is used for sorting slices of symbols according to symbol length. Symbols with equal length will be sorted alphabetically.
type symbolSlice []Symbol

func (a symbolSlice) Len() int      { return len(a) }
func (a symbolSlice) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a symbolSlice) Less(i, j int) bool {
	if len(a[i].String) != len(a[j].String) {
		return len(a[i].String) > len(a[j].String)
	}
	return a[i].String < a[j].String
}

// SymbolSetSuffix defines the filename extension for symbol sets
var SymbolSetSuffix = ".sym"

func trimIfNeeded(s string) string {
	trimmed := strings.TrimSpace(s)
	if len(trimmed) > 0 {
		return trimmed
	}
	return s
}

// filterSymbolsByCat is used to filter out specific symbol types from the symbol set (syllabic, non syllabic, etc)
func filterSymbolsByCat(symbols []Symbol, types []SymbolCat) []Symbol {
	var res = make([]Symbol, 0)
	for _, s := range symbols {
		if containsCat(types, s.Cat) {
			res = append(res, s)
		}
	}
	return res
}

func buildIPARegexp(symbols []Symbol) (*regexp.Regexp, error) {
	return buildIPARegexpWithGroup(symbols, false, true)
}

func buildRegexp(symbols []Symbol) (*regexp.Regexp, error) {
	return buildRegexpWithGroup(symbols, false, true)
}

func buildIPARegexpWithGroup(symbols []Symbol, removeEmpty bool, anonGroup bool) (*regexp.Regexp, error) {
	sorted := make([]Symbol, len(symbols))
	copy(sorted, symbols)
	sort.Sort(symbolSlice(sorted))
	var acc = make([]string, 0)
	for _, s := range sorted {
		if removeEmpty {
			if len(s.String) > 0 {
				acc = append(acc, regexp.QuoteMeta(s.IPA.String))
			}
		} else {
			acc = append(acc, regexp.QuoteMeta(s.IPA.String))
		}
	}
	prefix := "(?:"
	if !anonGroup {
		prefix = "("
	}
	s := prefix + strings.Join(acc, "|") + ")"
	regexp.MustCompile(s)
	re, err := regexp.Compile(s)
	if err != nil {
		err = fmt.Errorf("couldn't compile regexp from string '%s' : %v", s, err)
		return nil, err
	}
	return re, nil
}

func buildRegexpWithGroup(symbols []Symbol, removeEmpty bool, anonGroup bool) (*regexp.Regexp, error) {
	sorted := make([]Symbol, len(symbols))
	copy(sorted, symbols)
	sort.Sort(symbolSlice(sorted))
	var acc = make([]string, 0)
	for _, s := range sorted {
		if removeEmpty {
			if len(s.String) > 0 {
				acc = append(acc, regexp.QuoteMeta(s.String))
			}
		} else {
			acc = append(acc, regexp.QuoteMeta(s.String))
		}
	}
	prefix := "(?:"
	if !anonGroup {
		prefix = "("
	}
	s := prefix + strings.Join(acc, "|") + ")"
	regexp.MustCompile(s)
	re, err := regexp.Compile(s)
	if err != nil {
		err = fmt.Errorf("couldn't compile regexp from string '%s' : %v", s, err)
		return nil, err
	}
	return re, nil
}

func containsCat(types []SymbolCat, t SymbolCat) bool {
	for _, t0 := range types {
		if t0 == t {
			return true
		}
	}
	return false
}

func contains(symbols []Symbol, symbol string) bool {
	for _, s := range symbols {
		if s.String == symbol {
			return true
		}
	}
	return false
}

/*
func indexOf(elements []string, element string) int {
	for i, s := range elements {
		if s == element {
			return i
		}
	}
	return -1
}
*/
func string2unicode(s string) string {
	res := ""
	for _, ch := range s {
		res = res + fmt.Sprintf("%U", ch)
	}
	return res
}
