package converter

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"

	"github.com/dlclark/regexp2"
	"github.com/stts-se/symbolset"
)

func isTest(s string) bool {
	return strings.HasPrefix(s, "TEST\t")
}

var testRe = regexp.MustCompile("^TEST\t([^\t]+)\t([^\t]+)$")

func parseTest(s string) (test, error) {
	var matchRes []string = testRe.FindStringSubmatch(s)
	if matchRes == nil {
		return test{}, fmt.Errorf("invalid symbol set definition: %s", s)
	}
	return test{from: matchRes[1], to: matchRes[2]}, nil
}

var commentAtEndRe = regexp.MustCompile("^(.*[^/]+)//+.*$")

func trimComment(s string) string {
	return strings.TrimSpace(commentAtEndRe.ReplaceAllString(s, "$1"))
}

func isComment(s string) bool {
	return strings.HasPrefix(s, "//")
}

func isFrom(s string) bool {
	return strings.HasPrefix(s, "FROM\t")
}

func isTo(s string) bool {
	return strings.HasPrefix(s, "TO\t")
}

func isRegexpRule(s string) bool {
	return strings.HasPrefix(s, "RE\t")
}

var regexpRuleRe = regexp.MustCompile("^RE\t([^\t]+)\t([^\t]+)$")

func parseRegexpRule(s string) (Rule, error) {
	var matchRes []string = regexpRuleRe.FindStringSubmatch(s)
	if matchRes == nil {
		return RegexpRule{}, fmt.Errorf("invalid regexp rule definition: %s", s)
	}
	from, err := regexp2.Compile(matchRes[1], regexp2.None)
	if err != nil {
		return RegexpRule{}, err
	}
	to := matchRes[2]
	return RegexpRule{From: from, To: to}, nil
}

func isSymbolRule(s string) bool {
	return strings.HasPrefix(s, "SYMBOL\t")
}

var symbolRuleRe = regexp.MustCompile("^SYMBOL\t([^\t]+)\t([^\t]+)$")

func parseSymbolRule(s string) (Rule, error) {
	var matchRes []string = symbolRuleRe.FindStringSubmatch(s)
	if matchRes == nil {
		return SymbolRule{}, fmt.Errorf("invalid symbol rule definition: %s", s)
	}
	from := matchRes[1]
	to := matchRes[2]
	return SymbolRule{From: from, To: to}, nil
}

func isBlankLine(s string) bool {
	return len(s) == 0
}

var symbolSetRe = regexp.MustCompile("^(FROM|TO)\t([^\t]+)$")

func parseSymbolSet(s string) (string, error) {
	var matchRes []string = symbolSetRe.FindStringSubmatch(s)
	if matchRes == nil {
		return "", fmt.Errorf("invalid symbol set definition: %s", s)
	}
	return matchRes[2], nil
}

//var fileSuffix = regexp.MustCompile(".[^.]+$")

// LoadFile loads a converter file and runs the specified tests
func LoadFile(symbolSets map[string]symbolset.SymbolSet, fName string) (Converter, TestResult, error) {
	name := filepath.Base(fName)
	var extension = filepath.Ext(name)
	name = name[0 : len(name)-len(extension)]
	var converter = Converter{Name: name}
	var err error
	fh, err := os.Open(filepath.Clean(fName))
	if err != nil {
		return Converter{}, TestResult{}, err
	}
	/* #nosec G307 */
	defer fh.Close()
	n := 0
	s := bufio.NewScanner(fh)
	var testLines []test
	for s.Scan() {
		if err := s.Err(); err != nil {
			return Converter{}, TestResult{}, err
		}
		n++
		l := trimComment(strings.TrimSpace(s.Text()))
		if isBlankLine(l) || isComment(l) {
		} else if isFrom(l) {
			ss, err := parseSymbolSet(l)
			if err != nil {
				return Converter{}, TestResult{}, err
			}
			if val, ok := symbolSets[ss]; ok {
				converter.From = val
			} else {
				return Converter{}, TestResult{}, fmt.Errorf("symbolset not defined: %s", ss)
			}
		} else if isTo(l) {
			ss, err := parseSymbolSet(l)
			if err != nil {
				return Converter{}, TestResult{}, err
			}
			if val, ok := symbolSets[ss]; ok {
				converter.To = val
			} else {
				return Converter{}, TestResult{}, fmt.Errorf("symbolset not defined: %s", ss)
			}
		} else if isSymbolRule(l) {
			rule, err := parseSymbolRule(l)
			if err != nil {
				return Converter{}, TestResult{}, err
			}
			converter.Rules = append(converter.Rules, rule)
		} else if isRegexpRule(l) {
			rule, err := parseRegexpRule(l)
			if err != nil {
				return Converter{}, TestResult{}, err
			}
			converter.Rules = append(converter.Rules, rule)
		} else if isTest(l) {
			test, err := parseTest(l)
			if err != nil {
				return Converter{}, TestResult{}, err
			}
			testLines = append(testLines, test)
		}
	}
	testRes, err := converter.Test(testLines)
	if err != nil {
		return Converter{}, TestResult{}, err
	}
	return converter, testRes, nil
}

// Suffix defines the suffix string for converter files (.cnv)
var Suffix = ".cnv"

// LoadFromDir loads a converters from the specified folder (all files with .cnv extension)
func LoadFromDir(symbolSets map[string]symbolset.SymbolSet, dirName string) (map[string]Converter, map[string]TestResult, error) {
	// list files in dir
	fileInfos, err := ioutil.ReadDir(dirName)
	if err != nil {
		return nil, nil, fmt.Errorf("failed reading symbol set dir : %w", err)
	}
	var fErrs error
	var convs = make(map[string]Converter)
	var res = make(map[string]TestResult)
	for _, fi := range fileInfos {
		var testResult = TestResult{OK: true}
		if strings.HasSuffix(fi.Name(), Suffix) {
			conv, testRes, err := LoadFile(symbolSets, filepath.Join(dirName, fi.Name()))
			if err != nil {
				thisErr := fmt.Errorf("could't load converter from file %s : %w", fi.Name(), err)
				if fErrs != nil {
					fErrs = fmt.Errorf("%v : %w", fErrs, thisErr)
				} else {
					fErrs = thisErr
				}
			} else {
				if !testRes.OK {
					testResult.OK = false
				}
				testResult.Errors = append(testResult.Errors, testRes.Errors...)
				res[conv.Name] = testResult
				// TODO check that x.Name doesn't already exist ?
				convs[conv.Name] = conv
			}
		}
	}

	if fErrs != nil {
		return nil, nil, fErrs
	}

	return convs, res, nil
}

// Test runs the input tests and returns a test result
func (c Converter) Test(tests []test) (TestResult, error) {
	res1, err := c.testExamples(tests)
	if err != nil {
		return TestResult{}, err
	}
	res2, err := c.testInternals()
	if err != nil {
		return TestResult{}, err
	}
	if res1.OK && res2.OK {
		return TestResult{OK: true}, nil
	}
	return TestResult{OK: false, Errors: append(res1.Errors, res2.Errors...)}, nil
}

// runs pre-defined tests (defined in the input file)
func (c Converter) testExamples(tests []test) (TestResult, error) {
	errors := []string{}
	for _, test := range tests {
		result, err := c.Convert(test.from)
		if err != nil {
			errors = append(errors, fmt.Sprintf("%s", err))
			//return TestResult{}, err
		}
		if result != test.to {
			msg := fmt.Sprintf("From /%s/ expected /%s/, but got /%s/", test.from, test.to, result)
			errors = append(errors, msg)
		}
		invalid, err := c.getInvalidSymbols(result, c.To)
		if err != nil {
			return TestResult{}, err
		}
		if len(invalid) > 0 {
			errors = append(errors, fmt.Sprintf("Invalid symbol(s) in output transcription for test /%s/: %v", test, invalid))
		}
	}
	ok := (len(errors) == 0)
	return TestResult{OK: ok, Errors: errors}, nil
}

// runs internal tests
func (c Converter) testInternals() (TestResult, error) {
	errors := []string{}
	var symbolsThatNeedARule []string
	for _, phn := range c.From.Symbols {
		// check that all input symbols can be converted without errors
		res, err := c.Convert(phn.String)
		if err != nil {
			errors = append(errors, fmt.Sprintf("%s", err))
			//return TestResult{}, err
		}
		// check that all output symbols are valid as defined in c.To
		invalid, err := c.getInvalidSymbols(res, c.To)
		if err != nil {
			return TestResult{}, err
		}
		if len(invalid) > 0 {
			errors = append(errors, fmt.Sprintf("Invalid symbol(s) in output transcription /%s/: %v", res, invalid))
		}
		if !c.To.ValidSymbol(phn.String) {
			symbolsThatNeedARule = append(symbolsThatNeedARule, phn.String)
		}
	}

	// check that all input symbols that are not also part of the output symbol set, have a fallback rule
	for _, symbol := range symbolsThatNeedARule {
		var hasSymbolRule = false
		for _, rule := range c.Rules {
			if reflect.TypeOf(rule).Name() == "SymbolRule" {
				var sr = rule.(SymbolRule)
				if sr.From == symbol {
					hasSymbolRule = true
				}
			}
		}
		if !hasSymbolRule {
			errors = append(errors, fmt.Sprintf("Symbol rule needed for input phoneme /%s/", symbol))
		}
	}

	// for each symbol rule, check that input is defined in c.From, and output is defined in c.To
	for _, rule := range c.Rules {
		if reflect.TypeOf(rule).Name() == "SymbolRule" {
			var sr = rule.(SymbolRule)
			invalid, err := c.getInvalidSymbols(sr.From, c.From)
			if err != nil {
				return TestResult{}, err
			}
			if len(invalid) > 0 {
				errors = append(errors, fmt.Sprintf("Invalid symbol(s) in input transcription for rule %s: %v", rule, invalid))
			}
			invalid, err = c.getInvalidSymbols(sr.To, c.To)
			if err != nil {
				return TestResult{}, err
			}
			if len(invalid) > 0 {
				errors = append(errors, fmt.Sprintf("Invalid symbol(s) in output transcription for rule %s: %v", rule, invalid))
			}
		} else if reflect.TypeOf(rule).Name() == "RegexpRule" {
			var rr = rule.(RegexpRule)
			invalid, err := c.getInvalidSymbols(rr.To, c.To)
			if err != nil {
				return TestResult{}, err
			}
			if len(invalid) > 0 {
				errors = append(errors, fmt.Sprintf("Invalid symbol(s) in output transcription for rule %s: %v", rule, invalid))
			}
		}
	}
	ok := (len(errors) == 0)
	return TestResult{OK: ok, Errors: errors}, nil
}
