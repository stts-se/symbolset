package mapper

import (
	"fmt"
	"log"
	"strings"

	"github.com/stts-se/symbolset"
)

// LoadMapper loads a symbol set mapper from two SymbolSet instances
func LoadMapper(s1 symbolset.SymbolSet, s2 symbolset.SymbolSet) (Mapper, []symbolset.SymbolSetError, error) {
	fromName := s1.Name
	toName := s2.Name
	name := fromName + " - " + toName

	mapper := Mapper{name, s1, s2}

	var ssErrs []symbolset.SymbolSetError
	var errs []string

	for _, symbol := range s1.Symbols {
		if len(symbol.String) > 0 {
			mapped, ssErrs0, err := mapper.MapTranscription(symbol.String)
			if len(ssErrs0) > 0 {
				ssErrs = append(ssErrs, ssErrs0...)
			}
			if len(mapped) > 0 {
				if err != nil {
					return mapper, ssErrs, fmt.Errorf("couldn't test mapper: %v", err)
				}
			}
		}
	}
	if len(ssErrs) > 0 {
		return mapper, ssErrs, nil
	}
	if len(errs) > 0 {
		return mapper, nil, fmt.Errorf("mapper initialization tests failed : %v", strings.Join(errs, "; "))
	}

	return mapper, nil, nil
}

// LoadMapperFromFile loads two SymbolSet instances from files.
func LoadMapperFromFile(fromName string, toName string, fName1 string, fName2 string) (Mapper, []symbolset.SymbolSetError, error) {

	if fromName == toName {
		return Mapper{}, nil, fmt.Errorf("should not load symbol sets with the same name: %s", fromName)
	}
	if fName1 == fName2 {
		return Mapper{}, nil, fmt.Errorf("should not load both symbol sets from the same file: %s", fName1)
	}

	m1, err := symbolset.LoadSymbolSetWithName(fromName, fName1)
	if err != nil {
		log.Printf("couldn't load mapper: %v\n", err)
		return Mapper{}, nil, err
	}
	s2, err := symbolset.LoadSymbolSetWithName(toName, fName2)
	if err != nil {
		log.Printf("couldn't load mapper: %v\n", err)
		return Mapper{}, nil, err
	}
	return LoadMapper(m1, s2)
}
