package symbolset

// symbol set filters for accent/stress placement

import (
	"fmt"
	"regexp"
	"strings"
)

func preFilter(ss SymbolSet, trans string, fromType Type) (string, error) {
	if fromType == InternalIPA {
		return filterBeforeMappingFromInternalIPA(ss, trans)
	} else if fromType == IPA {
		return filterBeforeMappingFromIPA(ss, trans)
	} else if fromType == CMU {
		return filterBeforeMappingFromCMU(ss, trans)
	}
	return trans, nil
}

func postFilter(ss SymbolSet, trans string, toType Type) (string, error) {
	if toType == InternalIPA {
		return filterAfterMappingToInternalIPA(ss, trans)
	} else if toType == IPA {
		return filterAfterMappingToIPA(ss, trans)
	} else if toType == CMU {
		return filterAfterMappingToCMU(ss, trans)
	}
	return trans, nil
}

var ipaIndepStressRe = fmt.Sprintf("[%s%s]", ipaAccentI, ipaSecStress)
var ipaAccentI = "\u02C8"
var ipaAccentII = "\u0300"
var ipaSecStress = "\u02CC"
var ipaLength = "\u02D0"

//var ipaSyllDelim = "."
//var cmuString = "cmu"

func filterBeforeMappingFromInternalIPA(ss SymbolSet, trans string) (string, error) {
	// IPA: ˈba`ŋ.ka => ˈ`baŋ.ka"
	// IPA: ˈɑ̀ː.pa => ˈ`ɑː.pa
	trans = strings.Replace(trans, ipaAccentII+ipaLength, ipaLength+ipaAccentII, -1)
	s := ipaAccentI + "(" + ss.ipaPhonemeRe.String() + "+)" + ipaAccentII
	repl, err := regexp.Compile(s)
	if err != nil {
		return "", fmt.Errorf("couldn't compile regexp from string '%s' : %w", s, err)
	}
	res := repl.ReplaceAllString(trans, ipaAccentI+ipaAccentII+"$1")
	return res, nil
}

/*
	func canMapToIPA(ss SymbolSet, trans string) (bool, error) {
		symbols, err := ss.SplitIPATranscription(trans)
		if err != nil {
			return false, err
		}
		nSyllabic := 0
		foundDelim := false
		syllabicReS := "^(" + ss.ipaSyllabicRe.String() + ")$"
		syllabicRe, err := regexp.Compile(syllabicReS)
		if err != nil {
			return false, fmt.Errorf("cannot create ipa syllabic regexp from string /%s/", syllabicReS)
		}
		for _, sym := range symbols {
			if syllabicRe.MatchString(sym) {
				nSyllabic = nSyllabic + 1
			} else if sym == ipaSyllDelim {
				foundDelim = true
			}
		}
		if foundDelim == false && nSyllabic > 1 {
			return false, nil
		}
		return true, nil
	}
*/

func filterBeforeMappingFromIPA(ss SymbolSet, trans string) (string, error) {
	return filterBeforeMappingFromInternalIPA(ss, trans)
}

func filterAfterMappingToIPA(ss SymbolSet, trans string) (string, error) {
	// filter stress differently if the symbol set has a syllable delimiter
	// hasSyllDelim := false
	// for _, s := range filterSymbolsByCat(ss.Symbols, []SymbolCat{SyllableDelimiter}) {
	// 	if len(s.String) > 0 {
	// 		hasSyllDelim = true
	// 	}
	// }
	// if !hasSyllDelim {
	// 	return trans, nil
	// }

	// IPA: /t°Ɑlsyn`tEs/ => /t°Ɑlsynt`Es/
	s := "(" + ss.StressRe.String() + ")(" + ss.NonSyllabicRe.String() + "*)(" + ss.SyllabicRe.String() + ")"
	repl, err := regexp.Compile(s)
	if err != nil {
		return "", fmt.Errorf("couldn't compile regexp from string '%s' : %w", s, err)
	}
	trans = repl.ReplaceAllString(trans, "$2$1$3")

	// IPA: əs.ˈ̀̀e ...
	// IPA: /'`pa.pa/ => /'pa`.pa/
	accentIIConditionForAfterMapping := ipaAccentI + ipaAccentII
	if strings.Contains(trans, accentIIConditionForAfterMapping) {
		s := ipaAccentI + ipaAccentII + "(" + ss.NonSyllabicRe.String() + "*)(" + ss.SyllabicRe.String() + ")"
		repl, err := regexp.Compile(s)
		if err != nil {
			return "", fmt.Errorf("couldn't compile regexp from string '%s' : %w", s, err)
		}
		res := repl.ReplaceAllString(trans, ipaAccentI+"$1$2"+ipaAccentII)
		trans = res
	}
	// IPA: /'paː`.pa/ => /'pa`ː.pa/
	trans = strings.Replace(trans, ipaLength+ipaAccentII, ipaAccentII+ipaLength, -1)
	return trans, nil
}

func filterAfterMappingToInternalIPA(ss SymbolSet, trans string) (string, error) {

	// filter stress differently if the symbol set has a syllable delimiter
	hasSyllDelim := false
	for _, s := range filterSymbolsByCat(ss.Symbols, []SymbolCat{SyllableDelimiter}) {
		if len(s.String) > 0 {
			hasSyllDelim = true
		}
	}
	if !hasSyllDelim {
		return trans, nil
	}
	// create an error if the input transcription contains more than one syllabic, but no syllable delimiter
	// canMap, err := canMapToIPA(ss, trans)
	// if err != nil {
	// 	return "", err
	// }
	// if !canMap {
	// 	return trans, fmt.Errorf("cannot map transcription to IPA /%s/", trans)
	// }

	// IPA: /ə.ba⁀ʊˈt/ => /ə.ˈba⁀ʊt/
	s := "(" + ss.ipaNonSyllabicRe.String() + "*)(" + ss.ipaSyllabicRe.String() + ")(" + ipaIndepStressRe + ")"
	repl, err := regexp.Compile(s)
	if err != nil {
		return "", fmt.Errorf("couldn't compile regexp from string '%s' : %w", s, err)
	}
	trans = repl.ReplaceAllString(trans, "$3$1$2")

	// IPA: /ə.bˈa⁀ʊt/ => /ə.ˈba⁀ʊt/
	s = "(" + ss.ipaNonSyllabicRe.String() + "*)(" + ipaIndepStressRe + ")(" + ss.ipaSyllabicRe.String() + ")"
	repl, err = regexp.Compile(s)
	if err != nil {
		return "", fmt.Errorf("couldn't compile regexp from string '%s' : %w", s, err)
	}
	trans = repl.ReplaceAllString(trans, "$2$1$3")

	// IPA: əs.ˈ̀̀e ...
	// IPA: /'`pa.pa/ => /'pa`.pa/
	accentIIConditionForAfterMapping := ipaAccentI + ipaAccentII
	if strings.Contains(trans, accentIIConditionForAfterMapping) {
		s := ipaAccentI + ipaAccentII + "(" + ss.ipaNonSyllabicRe.String() + "*)(" + ss.ipaSyllabicRe.String() + ")"
		repl, err := regexp.Compile(s)
		if err != nil {
			return "", fmt.Errorf("couldn't compile regexp from string '%s' : %w", s, err)
		}
		res := repl.ReplaceAllString(trans, ipaAccentI+"$1$2"+ipaAccentII)
		trans = res
	}
	// IPA: /'paː`.pa/ => /'pa`ː.pa/
	trans = strings.Replace(trans, ipaLength+ipaAccentII, ipaAccentII+ipaLength, -1)
	return trans, nil
}

func filterBeforeMappingFromCMU(ss SymbolSet, trans string) (string, error) {
	re, err := regexp.Compile("([^ ]+)([012])")
	if err != nil {
		return "", err
	}
	trans = re.ReplaceAllString(trans, "$2 $1")
	return trans, nil
}

func filterAfterMappingToCMU(ss SymbolSet, trans string) (string, error) {
	s := "([012]) ((?:" + ss.NonSyllabicRe.String() + " )*)(" + ss.SyllabicRe.String() + ")"
	repl, err := regexp.Compile(s)
	if err != nil {
		return "", fmt.Errorf("couldn't compile regexp from string '%s' : %w", s, err)
	}
	trans = repl.ReplaceAllString(trans, "$2$3$1")

	trans = strings.Replace(trans, " 1", "1", -1)
	trans = strings.Replace(trans, " 2", "2", -1)
	trans = strings.Replace(trans, " 0", "0", -1)
	return trans, nil
}
