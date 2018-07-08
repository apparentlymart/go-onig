package onig

import (
	"fmt"
)

// Span describes a region of an input string or byte array by the indices
// of its bounds.
//
// Start is inclusive and End is exclusive, like a slice.
type Span struct {
	Start, End int
}

// Substr is the same as str[s.Start:s.End], offered for convenience.
func (s Span) Substr(str string) string {
	return str[s.Start:s.End]
}

// Slice is the same as b[s.Start:s.End], offered for convenience.
func (s Span) Slice(b []byte) []byte {
	return b[s.Start:s.End]
}

// Len is the same as s.End-s.Start, offered for convenience.
func (s Span) Len() int {
	return s.End - s.Start
}

func (s Span) GoString() string {
	return fmt.Sprintf("&onig.Span{%d, %d}", s.Start, s.End)
}
