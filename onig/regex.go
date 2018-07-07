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
// returning true if a match is found.
func (r *Regex) Match(s string, opts MatchOptions) bool {
	return regexMatch(r, s, opts)
}
