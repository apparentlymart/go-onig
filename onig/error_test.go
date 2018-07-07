package onig

import "testing"

func TestErrorError(t *testing.T) {
	err := onigError{
		code: errCodeEmptyCharClass,
	}

	got := err.Error()
	want := "empty char-class"
	if got != want {
		t.Errorf("wrong error message\ngot:  %s\nwant: %s", got, want)
	}
}
