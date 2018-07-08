package onig

// #cgo pkg-config: oniguruma
// #include <bindings.h>
import "C"

import (
	"runtime"
	"unsafe"
)

// NOTE WELL: This is written against Oniguruma 5.9.6. At the time of writing,
// latest master contains some API changes that may impact the behavior of
// these bindings.
//
// This is the only file in this package that is allowed to import "C". All
// other files must access C objects via unexported Go symbols defined in this
// file. No references to "C" may be visible in the package godoc.

type errorInfo [C.sizeof_OnigErrorInfo]byte

const regexSizeof = C.sizeof_regex_t

// For error_test.go only
const errCodeEmptyCharClass = C.ONIGERR_EMPTY_CHAR_CLASS

const (
	optIgnoreCase       CompileOptions = C.ONIG_OPTION_IGNORECASE
	optExtend           CompileOptions = C.ONIG_OPTION_EXTEND
	optMultiline        CompileOptions = C.ONIG_OPTION_MULTILINE
	optSingleline       CompileOptions = C.ONIG_OPTION_SINGLELINE
	optFindLongest      CompileOptions = C.ONIG_OPTION_FIND_LONGEST
	optFindNotEmpty     CompileOptions = C.ONIG_OPTION_FIND_NOT_EMPTY
	optNegateSingleline CompileOptions = C.ONIG_OPTION_NEGATE_SINGLELINE
	optDontCaptureGroup CompileOptions = C.ONIG_OPTION_DONT_CAPTURE_GROUP
	optCaptureGroup     CompileOptions = C.ONIG_OPTION_CAPTURE_GROUP
	optNotBOL           MatchOptions   = C.ONIG_OPTION_NOTBOL
	optNotEOL           MatchOptions   = C.ONIG_OPTION_NOTEOL
)

func init() {
	C.onig_init()

	SyntaxDefault = Syntax(unsafe.Pointer(C.OnigDefaultSyntax))
	SyntaxAsIs = Syntax(unsafe.Pointer(&C.OnigSyntaxASIS))
	SyntaxPosixBasic = Syntax(unsafe.Pointer(&C.OnigSyntaxPosixBasic))
	SyntaxPosixExtended = Syntax(unsafe.Pointer(&C.OnigSyntaxPosixExtended))
	SyntaxEmacs = Syntax(unsafe.Pointer(&C.OnigSyntaxEmacs))
	SyntaxGrep = Syntax(unsafe.Pointer(&C.OnigSyntaxGrep))
	SyntaxGNU = Syntax(unsafe.Pointer(&C.OnigSyntaxGnuRegex))
	SyntaxJava = Syntax(unsafe.Pointer(&C.OnigSyntaxJava))
	SyntaxPerl = Syntax(unsafe.Pointer(&C.OnigSyntaxPerl))
	SyntaxPerlNG = Syntax(unsafe.Pointer(&C.OnigSyntaxPerl_NG))
	SyntaxRuby = Syntax(unsafe.Pointer(&C.OnigSyntaxRuby))
}

func errStr(code int, info *errorInfo) string {
	buf := make([]byte, C.ONIG_MAX_ERROR_MESSAGE_LEN)
	l := C.goonig_error_code_to_str((*C.OnigUChar)(unsafe.Pointer(&buf[0])), C.int(code), info.cPtr())
	buf = buf[:l]
	return string(buf)
}

func regexInit(r *Regex, pattern string, options CompileOptions, syntax Syntax) error {
	patternC := C.CString(pattern)
	var errInfo errorInfo
	errInfoPtr := &errInfo
	errCode := C.goonig_init_regex(
		r.cPtr(),
		patternC,
		C.int(len(pattern)),
		options.cVal(),
		syntax.cPtr(),
		errInfoPtr.cPtr(),
	)
	if errCode != 0 {
		return onigError{
			code: int(errCode),
			info: errInfoPtr,
		}
	}
	runtime.SetFinalizer(r, func(r *Regex) {
		// Free any ancillary objects associated with the regex.
		C.goonig_free_regex(r.cPtr())
	})
	return nil
}

func regexMatch(r *Regex, s string, options MatchOptions) bool {
	sC := C.CString(s)
	result := C.goonig_regex_match(
		r.cPtr(),
		sC,
		C.int(len(s)),
		nil,
		options.cVal(),
	)
	return result >= 0
}

func regexMatchBytes(r *Regex, b []byte, options MatchOptions) bool {
	sC := (*C.char)(unsafe.Pointer(&b[0]))
	result := C.goonig_regex_match(
		r.cPtr(),
		sC,
		C.int(len(b)),
		nil,
		options.cVal(),
	)
	return result >= 0
}

func (r *Regex) cPtr() *C.regex_t {
	return (*C.regex_t)(unsafe.Pointer(&r.c[0]))
}

func (r *errorInfo) cPtr() *C.OnigErrorInfo {
	if r == nil {
		return nil
	}

	return (*C.OnigErrorInfo)(unsafe.Pointer(&r[0]))
}

func (s Syntax) cPtr() *C.OnigSyntaxType {
	return (*C.OnigSyntaxType)(unsafe.Pointer(uintptr(s)))
}

func (o CompileOptions) cVal() C.OnigOptionType {
	return C.OnigOptionType(o)
}

func (o MatchOptions) cVal() C.OnigOptionType {
	return C.OnigOptionType(o)
}
