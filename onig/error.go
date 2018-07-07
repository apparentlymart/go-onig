package onig

// onigError is an implementation of error used to return oniguruma errors.
type onigError struct {
	code int
	info *errorInfo
}

func (e onigError) Error() string {
	return errStr(e.code, e.info)
}
