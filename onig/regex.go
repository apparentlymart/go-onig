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

// SearchAround is equivalent to Search followed by slicing the string
// around the first match, if any.
//
// If there is no match, both the "before" and "inside" strings are empty and
// the "after" string exactly matches the input.
func (r *Regex) SearchAround(s string, opts MatchOptions) (before, inside, after string) {
	match := r.Search(s, opts)
	return match.SubstrAround(s)
}

// SearchAroundBytes is equivalent to SearchBytes followed by slicing the
// given slice around the first match, if any.
//
// If there is no match, both the "before" and "inside" slices are nil and
// the "after" slice exactly matches the input.
func (r *Regex) SearchAroundBytes(b []byte, opts MatchOptions) (before, inside, after []byte) {
	match := r.SearchBytes(b, opts)
	return match.SliceAround(b)
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

// CaptureCount returns the number of capture sequences present in the
// receiver. This is the highest number that can be passed to method Capture
// on any match returned from this regex.
func (r *Regex) CaptureCount() int {
	return regexCaptureCount(r)
}

// NamedCaptures returns a map from capture names in the regex to the numbers
// of those captures as could be passed to method Capture on any match returned
// from this regex.
//
// A particular name can appear multiple times, so this function returns all
// of them. For a helper that returns only the first number for each, see
// NamedCapturesSingle.
func (r *Regex) NamedCaptures() map[string][]int {
	entries := regexNameTable(r)
	if len(entries) == 0 {
		return nil
	}
	ret := make(map[string][]int, len(entries))
	for _, entry := range entries {
		ret[entry.Name] = append(ret[entry.Name], entry.Num)
	}
	return ret
}

// NamedCapturesFirst returns a map from capture names in the regex to the
// numbers of those captures as could be passed to method Capture on any match
// returned from this regex.
//
// A particular name can appear multiple times, so this function discards
// any subsequent re-definitions of a name to return only a single number for
// each.
func (r *Regex) NamedCapturesFirst() map[string]int {
	entries := regexNameTable(r)
	if len(entries) == 0 {
		return nil
	}
	ret := make(map[string]int, len(entries))
	for _, entry := range entries {
		if _, exists := ret[entry.Name]; exists {
			continue
		}
		ret[entry.Name] = entry.Num
	}
	return ret
}
