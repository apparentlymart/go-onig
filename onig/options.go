package onig

// CompileOptions is a bitmask type used to pass compile-time options to
// NewRegex.
type CompileOptions uint

const (
	NoCompileOpts       CompileOptions = 0
	OptIgnoreCase       CompileOptions = optIgnoreCase
	OptExtend           CompileOptions = optExtend
	OptMultiline        CompileOptions = optMultiline
	OptSingleline       CompileOptions = optSingleline
	OptFindLongest      CompileOptions = optFindLongest
	OptFindNotEmpty     CompileOptions = optFindNotEmpty
	OptNegateSingleline CompileOptions = optNegateSingleline
	OptDontCaptureGroup CompileOptions = optDontCaptureGroup
	OptCaptureGroup     CompileOptions = optCaptureGroup
)

// MatchOptions is a bitmask type used to pass match-time options to the
// various match and search methods on type Regex.
type MatchOptions uint

const (
	NoMatchOpts MatchOptions = 0
	OptNotBOL   MatchOptions = optNotBOL
	OptNotEOL   MatchOptions = optNotEOL
)
