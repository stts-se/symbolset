package symbolset

import (
	"reflect"
	"sort"
	"testing"
)

var vfs = "Wanted: '%#v' got: '%#v'"

func Test_LenSort(t *testing.T) {

	s0 := "sr"
	s1 := "shrt"
	s2 := "looong"
	ss := []string{s0, s1, s2}

	sort.Sort(byLength(ss))
	if got, want := ss[0], s2; got != want {
		t.Errorf("Got '%#v' Wanted '%#v'", got, want)
	}
	if got, want := ss[1], s1; got != want {
		t.Errorf("Got '%#v' Wanted '%#v'", got, want)
	}
	if got, want := ss[2], s0; got != want {
		t.Errorf("Got '%#v' Wanted '%#v'", got, want)
	}
}

func Test_splitIntoPhonemes(t *testing.T) {
	phs := []Symbol{
		{"aa", NonSyllabic, "", IPASymbol{"n/a", "n/a"}},
		{"a", NonSyllabic, "", IPASymbol{"n/a", "n/a"}},
		{"bb", NonSyllabic, "", IPASymbol{"n/a", "n/a"}},
		{"b", NonSyllabic, "", IPASymbol{"n/a", "n/a"}},
		{"ddddd", NonSyllabic, "", IPASymbol{"n/a", "n/a"}},
		{"f33", NonSyllabic, "", IPASymbol{"n/a", "n/a"}},
	}
	s1 := "c"
	res, unk, err := splitIntoPhonemes(phs, s1)
	if err != nil {
		t.Errorf("%s", err)
	}
	//fmt.Printf("res: '%#v' unk: '%#v'\n", res, unk)
	if got, want := res[0], s1; got != want {
		t.Errorf("Got '%#v' Wanted '%#v'", got, want)
	}
	if got, want := unk[0], s1; got != want {
		t.Errorf("Got '%#v' Wanted '%#v'", got, want)
	}

	s2 := "a"
	res, unk, err = splitIntoPhonemes(phs, s2)
	if err != nil {
		t.Errorf("%s", err)
	}
	if got, want := res[0], s2; got != want {
		t.Errorf("Got '%#v' Wanted '%#v'", got, want)
	}
	if got, want := len(unk), 0; got != want {
		t.Errorf("Got '%#v' Wanted '%#v'", got, want)
	}

	s3 := "_azbbax"
	res, unk, err = splitIntoPhonemes(phs, s3)
	if err != nil {
		t.Errorf("%s", err)
	}
	if got, want := res[0], "_"; got != want {
		t.Errorf("Got '%#v' Wanted '%#v'", got, want)
	}
	if got, want := res[3], "bb"; got != want {
		t.Errorf("Got '%#v' Wanted '%#v'", got, want)
	}

	if got, want := res[5], "x"; got != want {
		t.Errorf("Got '%#v' Wanted '%#v'", got, want)
	}

	if got, want := len(unk), 3; got != want {
		t.Errorf("Got '%#v' Wanted '%#v'", got, want)
	}
	if got, want := len(res), 6; got != want {
		t.Errorf("Got '%#v' Wanted '%#v'", got, want)
	}

	s4 := "a aa aaa f33"
	res, _, err = splitIntoPhonemes(phs, s4)
	if err != nil {
		t.Errorf("%s", err)
	}
	expect := []string{"a", " ", "aa", " ", "aa", "a", " ", "f33"}
	if !reflect.DeepEqual(expect, res) {
		t.Errorf(vfs, expect, res)
	}
}

func Test_splitIntoPhonemes2(t *testing.T) {
	phs1 := []Symbol{
		{"aa", Syllabic, "", IPASymbol{"n/a", "n/a"}},
		{"b", NonSyllabic, "", IPASymbol{"n/a", "n/a"}},
		{" ", PhonemeDelimiter, "", IPASymbol{" ", "n/a"}},
	}

	// _
	s1 := "a b a"
	_, _, err := splitIntoPhonemes(phs1, s1)
	if err == nil {
		t.Errorf("Expected error for phoneme list containing non empty delimiter")
	}

	phs2 := []Symbol{
		{"aa", Syllabic, "", IPASymbol{"n/a", "n/a"}},
		{"b", NonSyllabic, "", IPASymbol{"n/a", "n/a"}},
		{"", PhonemeDelimiter, "", IPASymbol{" ", "n/a"}},
	}

	// _
	s2 := "a b a"
	_, _, err = splitIntoPhonemes(phs2, s2)
	if err != nil {
		t.Errorf("didn't expect error here. Found %s", err)
	}
}

func Test_splitIntoPhonemesIPA(t *testing.T) {
	phs1 := []Symbol{
		{"a", Syllabic, "", IPASymbol{"a", ""}},
		{"b", NonSyllabic, "", IPASymbol{"b", ""}},
		{"r", NonSyllabic, "", IPASymbol{"r", ""}},
		{"k", NonSyllabic, "", IPASymbol{"k", ""}},
		{"", PhonemeDelimiter, "", IPASymbol{"", ""}},
		{"ɑː", Syllabic, "", IPASymbol{"ɑː", ""}},
		{".", SyllableDelimiter, "", IPASymbol{".", ""}},
		{"\u02C8", Stress, "", IPASymbol{"\u02C8", ""}},
		{"\u02C8\u0300", Stress, "", IPASymbol{"\u02C8\u0300", ""}},
	}

	// _
	s1 := "\u02C8brɑ\u0300ː.ka"
	_, _, err := splitIntoPhonemes(phs1, s1)
	if err != nil {
		t.Errorf("didn't expect error here. Found %s", err)
	}

	phs2 := []Symbol{
		{"aa", Syllabic, "", IPASymbol{"n/a", "n/a"}},
		{"b", NonSyllabic, "", IPASymbol{"n/a", "n/a"}},
		{"", PhonemeDelimiter, "", IPASymbol{" ", "n/a"}},
	}

	// _
	s2 := "\u02C8brɑː.ka"
	_, _, err = splitIntoPhonemes(phs2, s2)
	if err != nil {
		t.Errorf("didn't expect error here. Found %s", err)
	}
}
