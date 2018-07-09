package onig

// #cgo pkg-config: oniguruma
// #include <bindings.h>
import "C"

import (
	"reflect"
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

const (
	regexSizeof  = C.sizeof_regex_t
	regionSizeof = C.sizeof_OnigRegion
)

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

type nameTableEntry struct {
	Name string
	Num  int
}

func init() {
	C.onig_init()

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

func regexMatch(r *Regex, s string, options MatchOptions, m *Match) bool {
	sC := C.CString(s)
	result := C.goonig_regex_match(
		r.cPtr(),
		sC,
		C.int(len(s)),
		m.cPtr(),
		options.cVal(),
	)
	return result >= 0
}

func regexMatchBytes(r *Regex, b []byte, options MatchOptions, m *Match) bool {
	sC := (*C.char)(unsafe.Pointer(&b[0]))
	result := C.goonig_regex_match(
		r.cPtr(),
		sC,
		C.int(len(b)),
		m.cPtr(),
		options.cVal(),
	)
	return result >= 0
}

func regexSearch(r *Regex, s string, options MatchOptions, rev bool, m *Match) bool {
	sC := C.CString(s)
	revC := C.int(0)
	if rev {
		revC = C.int(1)
	}
	result := C.goonig_regex_search(
		r.cPtr(),
		sC,
		C.int(len(s)),
		revC,
		m.cPtr(),
		options.cVal(),
	)
	return result >= 0
}

func regexSearchBytes(r *Regex, b []byte, options MatchOptions, rev bool, m *Match) bool {
	sC := (*C.char)(unsafe.Pointer(&b[0]))
	revC := C.int(0)
	if rev {
		revC = C.int(1)
	}
	result := C.goonig_regex_search(
		r.cPtr(),
		sC,
		C.int(len(b)),
		revC,
		m.cPtr(),
		options.cVal(),
	)
	return result >= 0
}

func regexCaptureCount(r *Regex) int {
	result := C.goonig_regex_capture_count(r.cPtr())
	return int(result)
}

func regexNameTable(r *Regex) []nameTableEntry {
	l := regexCaptureCount(r)
	if l == 0 {
		return nil
	}
	table := make([]C.goonig_name_table_entry, l)
	realLen := C.goonig_regex_name_table(r.cPtr(), &table[0])
	if realLen == 0 {
		return nil
	}
	ret := make([]nameTableEntry, realLen)
	for i, raw := range table {
		nameLen := int(raw.len)
		nameBytes := *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
			uintptr(unsafe.Pointer(raw.start)), nameLen, nameLen,
		}))

		ret[i] = nameTableEntry{
			Name: string(nameBytes),
			Num:  int(raw.idx),
		}
	}
	return ret
}

func matchInit(m *Match) {
	C.goonig_init_region(m.cPtr())
	runtime.SetFinalizer(m, func(m *Match) {
		// Free any buffers associated with the match.
		C.goonig_free_region(m.cPtr())
	})
}

func matchInitFake(m *Match, spans []Span) error {
	matchInit(m)
	c := m.cPtr()
	l := len(spans)
	errCode := C.goonig_region_resize(c, C.int(l))
	if errCode != 0 {
		return onigError{
			code: int(errCode),
		}
	}
	begs := *(*[]C.int)(unsafe.Pointer(&reflect.SliceHeader{
		uintptr(unsafe.Pointer(c.beg)), l, l,
	}))
	ends := *(*[]C.int)(unsafe.Pointer(&reflect.SliceHeader{
		uintptr(unsafe.Pointer(c.end)), l, l,
	}))
	for i, span := range spans {
		begs[i] = C.int(span.Start)
		ends[i] = C.int(span.End)
	}
	return nil
}

func matchCaptureCount(m *Match) int {
	c := m.cPtr()
	return int(c.num_regs) - 1
}

func matchCapture(m *Match, idx int) Span {
	max := matchCaptureCount(m)
	if idx > max || idx < 0 {
		panic("capture index out of range")
	}

	c := m.cPtr()
	begPU := uintptr(unsafe.Pointer(c.beg))
	endPU := uintptr(unsafe.Pointer(c.end))
	begPU += uintptr(idx)
	endPU += uintptr(idx)
	begP := (*C.int)(unsafe.Pointer(begPU))
	endP := (*C.int)(unsafe.Pointer(endPU))

	return Span{
		Start: int(*begP),
		End:   int(*endP),
	}
}

func matchEqual(a *Match, b *Match) bool {
	if (a == nil) != (b == nil) {
		return false
	}
	if a == nil {
		return true
	}

	aC := a.cPtr()
	bC := b.cPtr()
	if aC.num_regs != bC.num_regs {
		return false
	}
	num := int(aC.num_regs)

	aBegs := *(*[]int)(unsafe.Pointer(&reflect.SliceHeader{
		uintptr(unsafe.Pointer(aC.beg)), num, num,
	}))
	bBegs := *(*[]int)(unsafe.Pointer(&reflect.SliceHeader{
		uintptr(unsafe.Pointer(bC.beg)), num, num,
	}))
	aEnds := *(*[]int)(unsafe.Pointer(&reflect.SliceHeader{
		uintptr(unsafe.Pointer(aC.end)), num, num,
	}))
	bEnds := *(*[]int)(unsafe.Pointer(&reflect.SliceHeader{
		uintptr(unsafe.Pointer(bC.end)), num, num,
	}))

	for i := range aBegs {
		if aBegs[i] != bBegs[i] {
			return false
		}
	}
	for i := range aEnds {
		if aEnds[i] != bEnds[i] {
			return false
		}
	}

	return true
}

func (r *Regex) cPtr() *C.regex_t {
	if r == nil {
		return nil
	}
	return (*C.regex_t)(unsafe.Pointer(&r.c[0]))
}

func (m *Match) cPtr() *C.OnigRegion {
	if m == nil {
		return nil
	}
	return (*C.OnigRegion)(unsafe.Pointer(&m.c[0]))
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
