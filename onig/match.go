package onig

import "fmt"

// Match describes how a pattern matched a particular string or byte array.
type Match struct {
	// c is a buffer into which the OnigRegion data will be placed.
	// It is opaque to Go code. Bindings code (in bindings.go) can access the
	// typed pointer to this via method cPtr.
	c [regionSizeof]byte
}

// Bounds returns a span describing the whole match.
func (m *Match) Bounds() Span {
	return m.Capture(0)
}

// Capture returns the number of captures.
func (m *Match) CaptureCount() int {
	return matchCaptureCount(m)
}

// Capture returns a span describing the capture with the given index, or
// panics if the index is out of bounds.
//
// Captures are 1-indexed, so the valid range for captures is from 1 to
// m.CaptureCount() inclusive. If index is zero, this method has the same
// result as method Bounds.
func (m *Match) Capture(index int) Span {
	return matchCapture(m, index)
}

// Equal returns true if the receiver and the other given match both describe
// the same bounds and captures.
func (m *Match) Equal(other *Match) bool {
	return matchEqual(m, other)
}

// GoString returns a Go-syntax-like representation of the receiver, which is
// primarily useful for debugging.
func (m *Match) GoString() string {
	if m == nil {
		return "(*onig.Match)(nil)"
	}
	l := m.CaptureCount()
	capts := make([]Span, l)
	for i := 0; i < l; i++ {
		capts[i] = m.Capture(i + 1)
	}

	return fmt.Sprintf("&onig.Match{Bounds:%#v,Captures:%#v}", m.Bounds(), capts)
}

// mustFakeMatch is a helper for testing which allocates and initializes a
// Match object populated with the given spans.
//
// If the match allocation fails (due to running out of memory) then this
// function will panic.
func mustFakeMatch(spans []Span) *Match {
	m := new(Match)
	err := matchInitFake(m, spans)
	if err != nil {
		panic(err.Error())
	}
	return m
}
