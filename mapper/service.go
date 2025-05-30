package mapper

import (
	"fmt"
	"log"
	"strings"

	"github.com/stts-se/symbolset"
)

// functions for use by the mapper http service

// Service is a container for maintaining 'cached' mappers and their symbol sets. Please note that currently, MapperService need to be used as mutex, see lexserver/mapper.go
type Service struct {
	SymbolSets map[string]symbolset.SymbolSet
	Mappers    map[string]Mapper
}

// MapperNames lists the names for all loaded mappers
func (s Service) MapperNames() []string {
	var names = make([]string, 0)
	for name := range s.Mappers {
		names = append(names, name)
	}
	return names
}

// DeleteSymbolSet is used to delete a named symbol set from the cache. Deletes the named symbol set, and all mappers using this symbol set.
func (s Service) DeleteSymbolSet(ssName string) error {
	_, ok := s.SymbolSets[ssName]
	if !ok {
		return fmt.Errorf("no existing symbol set named %s", ssName)
	}
	delete(s.SymbolSets, ssName)
	log.Printf("Deleted symbol set %v from cache", ssName)
	for mName := range s.Mappers {
		if strings.HasPrefix(mName, ssName+" ") ||
			strings.HasSuffix(mName, " "+ssName) {
			delete(s.Mappers, mName)
			log.Printf("Deleted mapper %v from cache", mName)
		}
	}
	return nil
}

// DeleteMapper is used to delete a mapper the cache.
func (s Service) DeleteMapper(fromName string, toName string) error {
	name := fromName + " - " + toName
	for mName := range s.Mappers {
		if mName == name {
			delete(s.Mappers, mName)
			log.Printf("Deleted mapper %v from cache", mName)
		}
	}
	return nil
}

// Load is used to load a symbol set from file
func (s Service) Load(symbolSetFile string) error {
	ss, err := symbolset.LoadSymbolSet(symbolSetFile)
	if err != nil {
		return fmt.Errorf("couldn't load symbol set : %v", err)
	}
	s.SymbolSets[ss.Name] = ss
	log.Printf("Loaded symbol set %v into cache", ss.Name)
	return nil
}

// Clear is used to clear the cache (all loaded symbol sets and mappers)
func (s Service) Clear() {
	// TODO: MapperService need to be used as mutex, see lexserver/mapper.go
	//lint:ignore SA4005 faulty warning
	s.SymbolSets = make(map[string]symbolset.SymbolSet)
	//lint:ignore SA4005 faulty warning
	s.Mappers = make(map[string]Mapper)
}

func (s Service) getOrCreateMapper(fromName string, toName string) (Mapper, []symbolset.SymbolSetError, error) {
	var ssErrs []symbolset.SymbolSetError
	name := fromName + " - " + toName
	mapper, ok := s.Mappers[name]
	if ok {
		return mapper, nil, nil
	}

	var from, to symbolset.SymbolSet
	from, okFrom := s.SymbolSets[fromName]
	if !okFrom {
		ssErrs = append(ssErrs, symbolset.SymbolSetError{
			ErrorType: "no such symbol set",
			Values:    []string{fromName},
		})
		//return nilRes, fmt.Errorf("couldn't find left hand symbol set named '%s'", fromName)
	}
	to, okTo := s.SymbolSets[toName]
	if !okTo {
		ssErrs = append(ssErrs, symbolset.SymbolSetError{
			ErrorType: "no such symbol set",
			Values:    []string{toName},
		})
		// return nilRes, fmt.Errorf("couldn't find right hand symbol set named '%s'", toName)
	}
	mapper, ssErrs, err := LoadMapper(from, to)
	if err == nil {
		s.Mappers[name] = mapper
	}
	return mapper, nil, err
}

// Map is used by the server to map a transcription from one symbol set to another
func (s Service) Map(fromName string, toName string, trans string) (string, []symbolset.SymbolSetError, error) {
	var ssErrs []symbolset.SymbolSetError
	if toName == "ipa" {
		ss, ok := s.SymbolSets[fromName]
		if !ok {
			ssErrs = append(ssErrs, symbolset.SymbolSetError{
				ErrorType: "no such symbol set",
				Values:    []string{fromName},
			})
			return "", ssErrs, nil
		}
		return ss.ConvertToIPA(trans)
	} else if fromName == "ipa" {
		ss, ok := s.SymbolSets[toName]
		if !ok {
			ssErrs = append(ssErrs, symbolset.SymbolSetError{
				ErrorType: "no such symbol set",
				Values:    []string{toName},
			})
			return "", ssErrs, nil
		}
		return ss.ConvertFromIPA(trans)
	} else {
		mapper, ssErrs, err := s.getOrCreateMapper(fromName, toName)
		if err != nil {
			return "", nil, fmt.Errorf("couldn't create mapper from %s to %s : %v", fromName, toName, err)
		}
		if len(ssErrs) > 0 {
			return "", ssErrs, nil
		}
		res, ssErrs, err := mapper.MapTranscription(trans)
		if err != nil || len(ssErrs) > 0 {
			return res, ssErrs, err
		}
		return res, nil, nil
	}
}

// GetMapTable is used by the server to show/get a mapping table between two symbol sets
func (s Service) GetMapTable(fromName string, toName string) (Mapper, []symbolset.SymbolSetError, error) {
	mapper, ssErrs, err := s.getOrCreateMapper(fromName, toName)
	if err != nil {
		return Mapper{}, ssErrs, fmt.Errorf("couldn't create mapper from %s to %s : %v", fromName, toName, err)
	}
	return mapper, ssErrs, nil
}
