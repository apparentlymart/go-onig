package onig

// Syntax is an enumeration of regex syntaxes that can be passed to NewRegex.
type Syntax uintptr

var (
	SyntaxOniguruma     Syntax
	SyntaxAsIs          Syntax
	SyntaxPosixBasic    Syntax
	SyntaxPosixExtended Syntax
	SyntaxEmacs         Syntax
	SyntaxGrep          Syntax
	SyntaxGNU           Syntax
	SyntaxJava          Syntax
	SyntaxPerl          Syntax
	SyntaxPerlNG        Syntax
	SyntaxRuby          Syntax
)
