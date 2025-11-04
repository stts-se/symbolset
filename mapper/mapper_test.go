package mapper

import (
	"testing"

	"github.com/stts-se/symbolset"
)

var fsExpTrans = "Expected: /%v/ got: /%v/"

func testMapTranscription(t *testing.T, mapper Mapper, input string, expect string) {
	result, err := mapper.MapTranscription(input)
	if err != nil {
		t.Errorf("MapTranscription() didn't expect error here; input=/%s/, expect=/%s/ : %v", input, expect, err)
		return
	} else if result != expect {
		t.Errorf(fsExpTrans, expect, result)
	}
}

func Test_MapTranscription_EmptyDelimiterInInput1(t *testing.T) {
	symbols1 := []symbolset.Symbol{
		{"a", symbolset.Syllabic, "", symbolset.IPASymbol{"a", "U+0061"}},
		{"r", symbolset.NonSyllabic, "", symbolset.IPASymbol{"r", "U+0072"}},
		{"t", symbolset.NonSyllabic, "", symbolset.IPASymbol{"t", "U+0074"}},
		{"r*t", symbolset.NonSyllabic, "", symbolset.IPASymbol{"R", "U+0052"}},
		{"", symbolset.PhonemeDelimiter, "", symbolset.IPASymbol{"", ""}},
	}
	symbols2 := []symbolset.Symbol{
		{"A", symbolset.Syllabic, "", symbolset.IPASymbol{"a", "U+0061"}},
		{"R", symbolset.NonSyllabic, "", symbolset.IPASymbol{"r", "U+0072"}},
		{"T", symbolset.NonSyllabic, "", symbolset.IPASymbol{"t", "U+0074"}},
		{"RT", symbolset.NonSyllabic, "", symbolset.IPASymbol{"R", "U+0052"}},
		{" ", symbolset.PhonemeDelimiter, "", symbolset.IPASymbol{"", ""}},
	}
	ss1, err := symbolset.NewSymbolSet("sampa1", symbols1)
	if err != nil {
		t.Errorf("MapTranscription() didn't expect error here : %v", err)
		return
	}

	ss2, err := symbolset.NewSymbolSet("sampa2", symbols2)
	if err != nil {
		t.Errorf("MapTranscription() didn't expect error here : %v", err)
		return
	}
	ssm, err := LoadMapper(ss1, ss2)
	if err != nil {
		t.Errorf("MapTranscription() didn't expect error here : %v", err)
		return
	}

	// --
	input := "ar*ttr"
	expect := "A RT T R"
	result, err := ssm.MapTranscription(input)
	if err != nil {
		t.Errorf("MapTranscription() didn't expect error here : %v", err)
	}
	if result != expect {
		t.Errorf(fsExpTrans, expect, result)
	}

	// --
	input = "ar*trt"
	expect = "A RT R T"
	result, err = ssm.MapTranscription(input)
	if err != nil {
		t.Errorf("MapTranscription() didn't expect error here : %v", err)
	}
	if result != expect {
		t.Errorf(fsExpTrans, expect, result)
	}
}

// func Test_MapTranscription_EmptyDelimiterInInput2(t *testing.T) {
// 	fromName := "ssLC"
// 	toName := "ssIPA"
// 	symbols := []Symbol{
// 		{"a", Syllabic, "", IPASymbol{"A", ""}},
// 		{"r", NonSyllabic, "", IPASymbol{"R", ""}},
// 		{"t", NonSyllabic, "", IPASymbol{"T", ""}},
// 		{"r*t", NonSyllabic, "", IPASymbol{"RT", ""}},
// 		{"", PhonemeDelimiter, "", IPASymbol{" ", ""}},
// 	}
// 	ssm, err := NewSymbolSet("test", symbols)
// 	if err != nil {
// 		t.Errorf("MapTranscription() didn't expect error here : %v", err)
// 	}
// 	input := "ar*ttrrt"
// 	expect := "A RT T R R T"
// 	result, err := ssm.MapTranscription(input)
// 	if err != nil {
// 		t.Errorf("MapTranscription() didn't expect error here : %v", err)
// 	}
// 	if result != expect {
// 		t.Errorf(fsExpTrans, expect, result)
// 	}
// }

// func Test_MapTranscription_EmptyDelimiterInOutput(t *testing.T) {
// 	fromName := "ssLC"
// 	toName := "ssIPA"
// 	symbols := []Symbol{
// 		{"a", Syllabic, "", IPASymbol{"A", ""}},
// 		{"r", NonSyllabic, "", IPASymbol{"R", ""}},
// 		{"t", NonSyllabic, "", IPASymbol{"T", ""}},
// 		{"rt", NonSyllabic, "", IPASymbol{"R*T", ""}},
// 		{" ", PhonemeDelimiter, "", IPASymbol{"", ""}},
// 	}
// 	ssm, err := NewSymbolSet("test", symbols)
// 	if err != nil {
// 		t.Errorf("MapTranscription() didn't expect error here : %v", err)
// 	}
// 	input := "a rt r t"
// 	expect := "AR*TRT"
// 	result, err := ssm.MapTranscription(input)
// 	if err != nil {
// 		t.Errorf("MapTranscription() didn't expect error here : %v", err)
// 	}
// 	if result != expect {
// 		t.Errorf(fsExpTrans, expect, result)
// 	}
// }

// func Test_MapTranscription_FailWithUnknownSymbols_EmptyDelim(t *testing.T) {
// 	fromName := "sampa1"
// 	toName := "ipa2"
// 	symbols := []Symbol{
// 		{"a", Syllabic, "", IPASymbol{"A", ""}},
// 		{"b", NonSyllabic, "", IPASymbol{"b", ""}},
// 		{"ŋ", NonSyllabic, "", IPASymbol{"N", ""}},
// 		{"", PhonemeDelimiter, "", IPASymbol{"", ""}},
// 		{".", SyllableDelimiter, "", IPASymbol{"$", SyllableDelimiter, ""}},
// 		{"\"", Stress, "", IPASymbol{"\"", Stress, ""}},
// 		{"\"\"", Stress, "", IPASymbol{"\"\"", Stress, ""}},
// 	}
// 	ssm, err := NewSymbolSet("test", symbols)
// 	if err != nil {
// 		t.Errorf("MapTranscription() didn't expect error here : %v", err)
// 	}
// 	input := "\"\"baŋ.ka"
// 	result, err := ssm.MapTranscription(input)
// 	if err == nil {
// 		t.Errorf("NewSymbolSet() expected error here, but got %s", result)
// 	}
// }

// func Test_MapTranscription_Ipa2Sampa_WithSwedishStress_1(t *testing.T) {
// 	fromName := "ipa"
// 	toName := "sampa"
// 	symbols := []Symbol{
// 		{"a", Syllabic, "", IPASymbol{"a", ""}},
// 		{"b", NonSyllabic, "", IPASymbol{"b", ""}},
// 		{"k", NonSyllabic, "", IPASymbol{"k", ""}},
// 		{"ŋ", NonSyllabic, "", IPASymbol{"N", ""}},
// 		{"", PhonemeDelimiter, "", IPASymbol{"", ""}},
// 		{".", SyllableDelimiter, "", IPASymbol{"$", SyllableDelimiter, ""}},
// 		{"\u02C8", Stress, "", IPASymbol{"\"", Stress, ""}},
// 		{"\u02C8\u0300", Stress, "", IPASymbol{"\"\"", Stress, ""}},
// 	}
// 	ssm, err := NewSymbolSet("test", symbols)
// 	if err != nil {
// 		t.Errorf("MapTranscription() didn't expect error here : %v", err)
// 	}
// 	input := "\u02C8ba\u0300ŋ.ka" // => ˈ`baŋ.ka before mapping
// 	expect := "\"\"baN$ka"
// 	result, err := ssm.MapTranscription(input)
// 	if err != nil {
// 		t.Errorf("MapTranscription() didn't expect error here : %v", err)
// 	}
// 	if result != expect {
// 		t.Errorf(fsExpTrans, expect, result)
// 	}
// }

// func Test_MapTranscription_Ipa2Sampa_WithSwedishStress_2(t *testing.T) {
// 	fromName := "ipa"
// 	toName := "sampa"
// 	symbols := []Symbol{
// 		{"a", Syllabic, "", IPASymbol{"a", ""}},
// 		{"b", NonSyllabic, "", IPASymbol{"b", ""}},
// 		{"k", NonSyllabic, "", IPASymbol{"k", ""}},
// 		{"ŋ", NonSyllabic, "", IPASymbol{"N", ""}},
// 		{"", PhonemeDelimiter, "", IPASymbol{"", ""}},
// 		{".", SyllableDelimiter, "", IPASymbol{"$", SyllableDelimiter, ""}},
// 		{"\u02C8", Stress, "", IPASymbol{"\"", Stress, ""}},
// 		{"\u02C8\u0300", Stress, "", IPASymbol{"\"\"", Stress, ""}},
// 	}
// 	ssm, err := NewSymbolSet("test", symbols)
// 	if err != nil {
// 		t.Errorf("MapTranscription() didn't expect error here : %v", err)
// 	}
// 	input := "\u02C8a\u0300ŋ.ka"
// 	expect := "\"\"aN$ka"
// 	result, err := ssm.MapTranscription(input)
// 	if err != nil {
// 		t.Errorf("MapTranscription() didn't expect error here : %v", err)
// 	}
// 	if result != expect {
// 		t.Errorf(fsExpTrans, expect, result)
// 	}
// }

// func Test_MapTranscription_Ipa2Sampa_WithSwedishStress_3(t *testing.T) {
// 	fromName := "ipa"
// 	toName := "sampa"
// 	symbols := []Symbol{
// 		{"a", Syllabic, "", IPASymbol{"a", ""}},
// 		{"b", NonSyllabic, "", IPASymbol{"b", ""}},
// 		{"r", NonSyllabic, "", IPASymbol{"r", ""}},
// 		{"k", NonSyllabic, "", IPASymbol{"k", ""}},
// 		{"ŋ", NonSyllabic, "", IPASymbol{"N", ""}},
// 		{"", PhonemeDelimiter, "", IPASymbol{"", ""}},
// 		{".", SyllableDelimiter, "", IPASymbol{"$", SyllableDelimiter, ""}},
// 		{"\u02C8", Stress, "", IPASymbol{"\"", Stress, ""}},
// 		{"\u02C8\u0300", Stress, "", IPASymbol{"\"\"", Stress, ""}},
// 	}
// 	ssm, err := NewSymbolSet("test", symbols)
// 	if err != nil {
// 		t.Errorf("MapTranscription() didn't expect error here : %v", err)
// 	}
// 	input := "\u02C8bra\u0300ŋ.ka"
// 	expect := "\"\"braN$ka"
// 	result, err := ssm.MapTranscription(input)
// 	if err != nil {
// 		t.Errorf("MapTranscription() didn't expect error here : %v", err)
// 	}
// 	if result != expect {
// 		t.Errorf(fsExpTrans, expect, result)
// 	}
// }

// func Test_MapTranscription_NstXSAMPA_To_WsSAMPA_1(t *testing.T) {
// 	fromName := "NST-XSAMPA"
// 	toName := "WS-SAMPA_IPADUMMY"
// 	symbols := []Symbol{
// 		{"a", Syllabic, "", IPASymbol{"a", ""}},
// 		{"b", NonSyllabic, "", IPASymbol{"b", ""}},
// 		{"r", NonSyllabic, "", IPASymbol{"r", ""}},
// 		{"k", NonSyllabic, "", IPASymbol{"k", ""}},
// 		{"N", NonSyllabic, "", IPASymbol{"N", ""}},
// 		{" ", PhonemeDelimiter, "", IPASymbol{" ", ""}},
// 		{"$", SyllableDelimiter, "", IPASymbol{".", SyllableDelimiter, ""}},
// 		{"\"", Stress, "", IPASymbol{"\"", Stress, ""}},
// 		{"\"\"", Stress, "", IPASymbol{"\"\"", Stress, ""}},
// 	}
// 	ssm, err := NewSymbolSet("test", symbols)
// 	if err != nil {
// 		t.Errorf("MapTranscription() didn't expect error here : %v", err)
// 	}
// 	input := "\"\" b r a N $ k a"
// 	expect := "\"\" b r a N . k a"
// 	result, err := ssm.MapTranscription(input)
// 	if err != nil {
// 		t.Errorf("MapTranscription() didn't expect error here : %v", err)
// 	}
// 	if result != expect {
// 		t.Errorf(fsExpTrans, expect, result)
// 	}
// }

// func Test_MapTranscription_NstXSAMPA_To_WsSAMPA_2(t *testing.T) {
// 	fromName := "NST-XSAMPA"
// 	toName := "WS-SAMPA_IPADUMMY"
// 	symbols := []Symbol{
// 		{"a", Syllabic, "", IPASymbol{"a", ""}},
// 		{"b", NonSyllabic, "", IPASymbol{"b", ""}},
// 		{"r", NonSyllabic, "", IPASymbol{"r", ""}},
// 		{"rs", NonSyllabic, "", IPASymbol{"rs", ""}},
// 		{"s", NonSyllabic, "", IPASymbol{"s", ""}},
// 		{"k", NonSyllabic, "", IPASymbol{"k", ""}},
// 		{"N", NonSyllabic, "", IPASymbol{"N", ""}},
// 		{" ", PhonemeDelimiter, "", IPASymbol{" ", ""}},
// 		{"$", SyllableDelimiter, "", IPASymbol{".", SyllableDelimiter, ""}},
// 		{"\"", Stress, "", IPASymbol{"\"", Stress, ""}},
// 		{"\"\"", Stress, "", IPASymbol{"\"\"", Stress, ""}},
// 	}
// 	ssm, err := NewSymbolSet("test", symbols)
// 	if err != nil {
// 		t.Errorf("MapTranscription() didn't expect error here : %v", err)
// 	}
// 	input := "\"\" b r a $ rs a r s"
// 	expect := "\"\" b r a . rs a r s"
// 	result, err := ssm.MapTranscription(input)
// 	if err != nil {
// 		t.Errorf("MapTranscription() didn't expect error here : %v", err)
// 	}
// 	if result != expect {
// 		t.Errorf(fsExpTrans, expect, result)
// 	}
// }

// func Test_loadSymbolSet_NST2WS(t *testing.T) {
// 	name := "NST-XSAMPA"
// 	fromColumn := "SAMPA"
// 	toColumn := "IPA"
// 	fName := "../test_data/sv-se_nst-xsampa.sym"
// 	ssmNST, err := loadSymbolSet_(name, fName, fromColumn, toColumn)
// 	if err != nil {
// 		t.Errorf("MapTranscription() didn't expect error here : %v", err)
// 	}

// 	name = "WS-SAMPA"
// 	fromColumn = "IPA"
// 	toColumn = "SYMBOL"
// 	fName = "../test_data/sv-se_ws-sampa.sym"
// 	ssmWS, err := loadSymbolSet_(name, fName, fromColumn, toColumn)
// 	if err != nil {
// 		t.Errorf("MapTranscription() didn't expect error here : %v", err)
// 	}

// 	mappers := []SymbolSet{ssmNST, ssmWS}

// 	testMapTranscriptionX(t, mappers, "\"bOt`", "\" b O rt")
// 	testMapTranscriptionX(t, mappers, "\"ku0rd", "\" k u0 r d")
// }

// func Test_NewSymbolSet_FailIfInputContainsDuplicates(t *testing.T) {
// 	fromName := "ssLC"
// 	toName := "ssUC"
// 	symbols := []Symbol{
// 		{"A", NonSyllabic, "", IPASymbol{"a", ""}},
// 		{"A", Syllabic, "", IPASymbol{"A", ""}},
// 		{"p", NonSyllabic, "", IPASymbol{"P", ""}},
// 		{" ", PhonemeDelimiter, "", IPASymbol{" ", ""}},
// 	}
// 	_, err := NewSymbolSet("test", symbols)
// 	if err == nil {
// 		t.Errorf("NewSymbolSet() expected error when input contains duplicates")
// 	}
// }
// func Test_loadSymbolSet_CMU2MARY(t *testing.T) {
// 	name := "CMU2IPA"
// 	fromColumn := "CMU"
// 	toColumn := "IPA"
// 	fName := "../test_data/en-us_cmu.sym"
// 	ssmCMU, err := loadSymbolSet_(name, fName, fromColumn, toColumn)
// 	if err != nil {
// 		t.Errorf("MapTranscription() didn't expect error here : %v", err)
// 		return
// 	}

// 	name = "IPA2MARY"
// 	fromColumn = "IPA"
// 	toColumn = "SYMBOL"
// 	fName = "../test_data/en-us_sampa_mary.sym"
// 	ssmMARY, err := loadSymbolSet_(name, fName, fromColumn, toColumn)
// 	if err != nil {
// 		t.Errorf("MapTranscription() didn't expect error here : %v", err)
// 		return
// 	}

// 	mappers := []SymbolSet{ssmCMU, ssmMARY}

// 	testMapTranscriptionX(t, mappers, "AX $ B AW1 T", "@ - \" b aU t")
// }

// func Test_loadSymbolSet_SAMPA2MARY(t *testing.T) {
// 	name := "SAMPA2IPA"
// 	fromColumn := "SYMBOL"
// 	toColumn := "IPA"
// 	fName := "../test_data/sv-se_ws-sampa.sym"
// 	ssm1, err := loadSymbolSet_(name, fName, fromColumn, toColumn)
// 	if err != nil {
// 		t.Errorf("MapTranscription() didn't expect error here : %v", err)
// 		return
// 	}

// 	name = "IPA2MARY"
// 	fromColumn = "IPA"
// 	toColumn = "SAMPA"
// 	fName = "../test_data/sv-se_sampa_mary.sym"
// 	ssm2, err := loadSymbolSet_(name, fName, fromColumn, toColumn)
// 	if err != nil {
// 		t.Errorf("MapTranscription() didn't expect error here : %v", err)
// 		return
// 	}
// 	mappers := []SymbolSet{ssm1, ssm2}
// 	testMapTranscriptionX(t, mappers, "eu . r \" u: p a", "E*U - r ' u: p a")
// 	testMapTranscriptionX(t, mappers, "@ s . \"\" e", "e s - \" e")
// }

// func Test_loadSymbolSet_MARY2SAMPA(t *testing.T) {
// 	name := "MARY2IPA"
// 	fromColumn := "SAMPA"
// 	toColumn := "IPA"
// 	fName := "../test_data/sv-se_sampa_mary.sym"
// 	ssm1, err := loadSymbolSet_(name, fName, fromColumn, toColumn)
// 	if err != nil {
// 		t.Errorf("MapTranscription() didn't expect error here : %v", err)
// 		return
// 	}

// 	name = "IPA2SAMPA"
// 	fromColumn = "IPA"
// 	toColumn = "SYMBOL"
// 	fName = "../test_data/sv-se_ws-sampa.sym"
// 	ssm2, err := loadSymbolSet_(name, fName, fromColumn, toColumn)
// 	if err != nil {
// 		t.Errorf("MapTranscription() didn't expect error here : %v", err)
// 		return
// 	}
// 	mappers := []SymbolSet{ssm1, ssm2}
// 	testMapTranscriptionX(t, mappers, "E*U - r ' u: p a", "eu . r \" u: p a")
// 	testMapTranscriptionX(t, mappers, "e s - \" e", "e s . \"\" e")
// 	testMapTranscriptionX(t, mappers, "\" e: - p a", "\"\" e: . p a")
// 	testMapTranscriptionX(t, mappers, "\" A: - p a", "\"\" A: . p a")

// 	mapper, err := LoadMapperFromFile("SAMPA", "SYMBOL", "../test_data/sv-se_sampa_mary.sym", "../test_data/sv-se_ws-sampa.sym")
// 	if err != nil {
// 		t.Errorf("Test_LoadMapperFromFile() didn't expect error here : %v", err)
// 		return
// 	}

// 	testMapTranscriptionY(t, mapper, "\" e: - p a", "\"\" e: . p a")
// 	testMapTranscriptionY(t, mapper, "\" A: - p a", "\"\" A: . p a")
// }

// func Test_loadSymbolSet_NST2MARY(t *testing.T) {
// 	name := "NST2IPA"
// 	fromColumn := "SAMPA"
// 	toColumn := "IPA"
// 	fName := "../test_data/sv-se_nst-xsampa.sym"
// 	ssm1, err := loadSymbolSet_(name, fName, fromColumn, toColumn)
// 	if err != nil {
// 		t.Errorf("MapTranscription() didn't expect error here : %v", err)
// 		return
// 	}

// 	name = "IPA2MARY"
// 	fromColumn = "IPA"
// 	toColumn = "SAMPA"
// 	fName = "../test_data/sv-se_sampa_mary.sym"
// 	ssm2, err := loadSymbolSet_(name, fName, fromColumn, toColumn)
// 	if err != nil {
// 		t.Errorf("MapTranscription() didn't expect error here : %v", err)
// 		return
// 	}
// 	mappers := []SymbolSet{ssm1, ssm2}
// 	testMapTranscriptionX(t, mappers, "E*U$r\"u:t`a", "E*U - r ' u: rt a")
// }

// func Test_loadSymbolSet_NST2SAMPA(t *testing.T) {
// 	name := "NST2IPA"
// 	fromColumn := "SAMPA"
// 	toColumn := "IPA"
// 	fName := "../test_data/sv-se_nst-xsampa.sym"
// 	ssm1, err := loadSymbolSet_(name, fName, fromColumn, toColumn)
// 	if err != nil {
// 		t.Errorf("MapTranscription() didn't expect error here : %v", err)
// 		return
// 	}

// 	name = "IPA2SAMPA"
// 	fromColumn = "IPA"
// 	toColumn = "SYMBOL"
// 	fName = "../test_data/sv-se_ws-sampa.sym"
// 	ssm2, err := loadSymbolSet_(name, fName, fromColumn, toColumn)
// 	if err != nil {
// 		t.Errorf("MapTranscription() didn't expect error here : %v", err)
// 		return
// 	}

// 	mappers := []SymbolSet{ssm1, ssm2}
// 	testMapTranscriptionX(t, mappers, "\"kaj$rU", "\" k a j . r U")
// 	testMapTranscriptionX(t, mappers, "E*U$r\"u:t`a", "eu . r \" u: rt a")
// }

// func Test_loadSymbolSet_MARY2CMU(t *testing.T) {
// 	name := "MARY2IPA"
// 	fromColumn := "SYMBOL"
// 	toColumn := "IPA"
// 	fName := "../test_data/en-us_sampa_mary.sym"
// 	ssmMARY, err := loadSymbolSet_(name, fName, fromColumn, toColumn)
// 	if err != nil {
// 		t.Errorf("MapTranscription() didn't expect error here : %v", err)
// 		return
// 	}

// 	name = "IPA2CMU"
// 	fromColumn = "IPA"
// 	toColumn = "CMU"
// 	fName = "../test_data/en-us_cmu.sym"
// 	ssmCMU, err := loadSymbolSet_(name, fName, fromColumn, toColumn)
// 	if err != nil {
// 		t.Errorf("MapTranscription() didn't expect error here : %v", err)
// 		return
// 	}

// 	mappers := []SymbolSet{ssmMARY, ssmCMU}

// 	testMapTranscriptionX(t, mappers, "@ - \" b aU t", "AX $ B AW1 T")
// 	testMapTranscriptionX(t, mappers, "V - \" b aU t", "AH $ B AW1 T")
// }

// func Test_LoadMapperFromFile_MARY2CMU(t *testing.T) {
// 	mappers, err := LoadMapperFromFile("SYMBOL", "CMU", "../test_data/en-us_sampa_mary.sym", "../test_data/en-us_cmu.sym")
// 	if err != nil {
// 		t.Errorf("Test_LoadMapperFromFile() didn't expect error here : %v", err)
// 		return
// 	}

// 	testMapTranscriptionY(t, mappers, "@ - \" b aU t", "AX $ B AW1 T")
// 	testMapTranscriptionY(t, mappers, "V - \" b aU t", "AH $ B AW1 T")
// }

func Test_LoadMapperFromFile_NST2WS(t *testing.T) {
	mapper, err := LoadMapperFromFile("SAMPA", "SYMBOL", "../test_data/nb-no_nst-xsampa.sym", "../test_data/nb-no_ws-sampa.sym")
	if err != nil {
		t.Errorf("Test_LoadMapperFromFile() didn't expect error here : %v", err)
		return
	}

	testMapTranscription(t, mapper, "\"A:$bl@s", "\" A: . b l @ s")
	testMapTranscription(t, mapper, "\"tSE$kIsk", "\" t S e . k i s k")
	testMapTranscription(t, mapper, "\"\"b9$n@r", "\"\" b 2 . n @ r")
	testMapTranscription(t, mapper, "\"b9$n@r", "\" b 2 . n @ r")
	testMapTranscription(t, mapper, "b\"9n", "\" b 2 n")
}

func Test_LoadMapperFromFile_FailIfBothHaveTheSameName(t *testing.T) {
	_, err := LoadMapperFromFile("SAMPA", "SAMPA", "../test_data/nb-no_nst-xsampa.sym", "../test_data/nb-no_ws-sampa.sym")
	if err == nil {
		t.Errorf("LoadMapperFromFile() expected error here")
	}
}

func Test_LoadMapperFromFile_FailIfBothHaveTheSameFile(t *testing.T) {
	_, err := LoadMapperFromFile("XSAMPA", "SAMPA", "../test_data/nb-no_nst-xsampa.sym", "../test_data/nb-no_nst-xsampa.sym")
	if err == nil {
		t.Errorf("LoadMapperFromFile() expected error here")
	}
}

func Test_MapperFromFile_CMU2WS_NoSyllDelim(t *testing.T) {
	mapper, err := LoadMapperFromFile("ENU-CMU", "ENU-WS", "../test_data/en-us_cmu-nosylldelim.sym", "../test_data/en-us_ws-sampa.sym")
	if err != nil {
		t.Errorf("Test_LoadMapperFromFile() didn't expect error here : %v", err)
		return
	}

	testMapTranscription(t, mapper, "P L AE1 T AX P UH2 S", "p l ' { t @ p % U s")

	//_, err = mapper.MapTranscription(input)
	// if err == nil {
	// 	t.Errorf("Expected error here!")
	// }

}

func Test_MapperFromFile_CMU2WS_WithSyllDelim(t *testing.T) {
	mapper, err := LoadMapperFromFile("ENU-CMU", "ENU-WS", "../test_data/en-us_cmu.sym", "../test_data/en-us_ws-sampa.sym")
	if err != nil {
		t.Errorf("Test_LoadMapperFromFile() didn't expect error here : %v", err)
		return
	}

	//testMapTranscription(t, mapper, " ", " ")

	testMapTranscription(t, mapper, "P L AE1 $ T AX $ P UH2 S", "' p l { . t @ . % p U s")

	// input := "P L AE1 T AX P UH2 S"
	// _, err = mapper.MapTranscription(input)
	// if err == nil {
	// 	t.Errorf("Expected error here!")
	// }
}
