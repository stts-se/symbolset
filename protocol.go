package symbolset

import (
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

const (
	ErrCodeFileNotFound       = 24
	ErrCodeUnknownInputSymbol = 25
	ErrCodeUnknownSymbolType  = 26
	ErrCodeUnknownSymbolSet   = 27
)

// SymbolSetError : container
type SymbolSetError struct {
	ErrorType string   `json:"error_type"` // examples: unknown phoneme(s), etc
	ErrorCode int      `json:"error_code"` //
	Values    []string `json:"values"`
}

func (e *SymbolSetError) Error() string {
	return fmt.Sprintf("[%d]: %s: %v", e.ErrorCode, e.ErrorType, e.Values)
}

func UnknownInputSymbol(values []string) *SymbolSetError {
	return &SymbolSetError{
		ErrorType: "Unknown input symbol",
		ErrorCode: ErrCodeUnknownInputSymbol,
		Values:    values,
	}
}

func UnknownSymbolType(values []string) *SymbolSetError {
	return &SymbolSetError{
		ErrorType: "Unknown symbol type",
		ErrorCode: ErrCodeUnknownSymbolType,
		Values:    values,
	}
}

func UnknownSymbolSet(values []string) *SymbolSetError {
	return &SymbolSetError{
		ErrorType: "Unknown symbol set",
		ErrorCode: ErrCodeUnknownSymbolSet,
		Values:    values,
	}
}

func (ss SymbolSetError) String() string {
	return fmt.Sprintf("[%s]: %s", ss.ErrorType, strings.Join(ss.Values, ", "))
}
