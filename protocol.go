package symbolset

import (
	"errors"
	"fmt"
	"strings"
)

// MapRequest : container
type MapRequest struct {
	From  string
	To    string
	Input string
}

// MapError : container
type MapError struct {
	Type      string     `json:"type"`       // error/result
	ErrorType string     `json:"error_type"` // examples: unknown phoneme(s), etc
	ErrorCode int        `json:"error_code"` //
	Values    []string   `json:"values"`
	Request   MapRequest `json:"request"`
}

func UnknownMapError() MapError {
	return MapError{
		ErrorType: "unknown",
		ErrorCode: 99,
	}
}

// SymbolSetError : container
type SymbolSetError struct {
	ErrorType string   `json:"error_type"` // examples: unknown phoneme(s), etc
	ErrorCode int      `json:"error_code"` //
	Values    []string `json:"values"`
}

func UnknownInputSymbol() SymbolSetError {
	return SymbolSetError{
		ErrorType: "Unknown input symbol",
		ErrorCode: 25,
	}
}

func (ss SymbolSetError) String() string {
	return fmt.Sprintf("[%s]: %s", ss.ErrorType, strings.Join(ss.Values, ", "))
}

func (ss SymbolSetError) Error() error {
	return fmt.Errorf("%s: %s", ss.ErrorType, strings.Join(ss.Values, ", "))
}

func SymbolSetErrors2Error(ssErrs []SymbolSetError) error {
	var errs []string
	for _, ssErr := range ssErrs {
		errs = append(errs, ssErr.String())
	}
	return errors.New(strings.Join(errs, "; "))
}
