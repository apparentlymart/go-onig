package onig

// Regex is the main type in this package, representing a compiled regular
// expression.
type Regex struct {
	// c is a buffer into which the oniguruma regex_t data will be placed.
	// It is opaque to Go code. Bindings code (in bindings.go) can access the
	// typed pointer to this via method cPtr.
	c [regexSizeof]byte
}

// NewRegex compiles the given regex pattern using the selected syntax,
// returning a newly-allocated Regex object.
func NewRegex(pattern string, options CompileOptions, syntax Syntax) (*Regex, error) {
	r := new(Regex)
	err := regexInit(r, pattern, options, syntax)
	if err != nil {
		// Don't return our probably-invalid Regex object, since accessing it
		// is likely to cause crashes.
		return nil, err
	}
	return r, nil
}

// Match tests whether the receiver matches a prefix of the given string,
// returning a description of the match if one is found. If no match is found
// then the result is nil.
func (r *Regex) Match(s string, opts MatchOptions) *Match {
	m := new(Match)
	matchInit(m)
	matches := regexMatch(r, s, opts, m)
	if !matches {
		return nil
	}
	return m
}

// MatchBytes tests whether the receiver matches a prefix of the given byte
// slice, returning a description of the match if one is found. If no match is
// found then the result is nil.
func (r *Regex) MatchBytes(b []byte, opts MatchOptions) *Match {
	m := new(Match)
	matchInit(m)
	matches := regexMatchBytes(r, b, opts, m)
	if !matches {
		return nil
	}
	return m
}

// Search tests whether the receiver matches a substring of the given string,
// returning a description of the first match found. If no match is found then
// the result is nil.
func (r *Regex) Search(s string, opts MatchOptions) *Match {
	m := new(Match)
	matchInit(m)
	matches := regexSearch(r, s, opts, false, m)
	if !matches {
		return nil
	}
	return m
}

// SearchBytes tests whether the receiver matches a substring of the given byte
// slice, returning a description of the first match found. If no match is
// found then the result is nil.
func (r *Regex) SearchBytes(b []byte, opts MatchOptions) *Match {
	m := new(Match)
	matchInit(m)
	matches := regexSearchBytes(r, b, opts, false, m)
	if !matches {
		return nil
	}
	return m
}

// Matches tests whether the receiver matches a prefix of the given string,
// returning true if a match is found.
func (r *Regex) Matches(s string, opts MatchOptions) bool {
	return regexMatch(r, s, opts, nil)
}

// MatchesBytes tests whether the receiver matches a prefix of the given byte
// slice, returning true if a match is found.
func (r *Regex) MatchesBytes(b []byte, opts MatchOptions) bool {
	return regexMatchBytes(r, b, opts, nil)
}
