// Package onig provides regular expression handling using the Oniguruma regex
// engine.
//
// Oniguruma is written in C, so this package uses CGo.
//
// Since Go strings are conventionally UTF-8, this library initializes
// Oniguruma only with UTF-8 support.
package onig
