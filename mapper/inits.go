package mapper

import (
	"fmt"
	"log"

	"github.com/stts-se/symbolset"
)

// LoadMapper loads a symbol set mapper from two SymbolSet instances
func LoadMapper(s1 symbolset.SymbolSet, s2 symbolset.SymbolSet) (Mapper, error) {
	fromName := s1.Name
	toName := s2.Name
	name := fromName + " - " + toName

	mapper := Mapper{name, s1, s2}

	var err error

	for _, symbol := range s1.Symbols {
		if len(symbol.String) > 0 {
			mapped, err := mapper.MapTranscription(symbol.String)
			if err != nil {
				break
			}
			if len(mapped) > 0 {
				if err != nil {
					err = fmt.Errorf("couldn't test mapper: %w", err)
					break
				}
			}
		}
	}
	if err != nil {
		return mapper, fmt.Errorf("mapper initialization tests failed : %w", err)
	}

	return mapper, nil
}

// LoadMapperFromFile loads two SymbolSet instances from files.
func LoadMapperFromFile(fromName string, toName string, fName1 string, fName2 string) (Mapper, error) {

	if fromName == toName {
		return Mapper{}, fmt.Errorf("should not load symbol sets with the same name: %s", fromName)
	}
	if fName1 == fName2 {
		return Mapper{}, fmt.Errorf("should not load both symbol sets from the same file: %s", fName1)
	}

	m1, err := symbolset.LoadSymbolSetWithName(fromName, fName1)
	if err != nil {
		log.Printf("couldn't load mapper: %v\n", err)
		return Mapper{}, err
	}
	s2, err := symbolset.LoadSymbolSetWithName(toName, fName2)
	if err != nil {
		log.Printf("couldn't load mapper: %v\n", err)
		return Mapper{}, err
	}
	return LoadMapper(m1, s2)
}
