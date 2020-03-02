package converter

import (
	"strings"
	"testing"

	"github.com/stts-se/symbolset"
)

func TestLoadFromDir(t *testing.T) {
	sSets, err := symbolset.LoadSymbolSetsFromDir("../test_data")
	if err != nil {
		t.Errorf("LoadSymbolSetsFromDir() didn't expect error here : %v", err)
		return
	}
	_, testRes, err := LoadFromDir(sSets, "../test_data")
	if err != nil {
		t.Errorf("LoadSymbolSetsFromDir() didn't expect error here : %v", err)
		return
	}
	for name, res := range testRes {
		if !res.OK && !strings.Contains(name, "FAIL") {
			for _, err := range res.Errors {
				t.Errorf("%s: %s", name, err)
			}
		} else if strings.Contains(name, "FAIL") && res.OK {
			t.Errorf("EXPECTED FAIL FOR CONVERTER %s", name)
		}
	}
}
