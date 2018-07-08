package onig

import (
	"fmt"
	"testing"
)

func TestRegexMatches(t *testing.T) {
	tests := []struct {
		Pattern string
		Str     string
		Want    bool
	}{
		{
			`hello`,
			`hello world`,
			true,
		},
		{
			`he`,
			`hello world`,
			true,
		},
		{
			`hel*o`,
			`hellllllllllo world`,
			true,
		},
		{
			`world`,
			`hello world`,
			false, // Prefix must match
		},
		{
			`hello`,
			`goodbye world`,
			false,
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%q in %q", test.Pattern, test.Str), func(t *testing.T) {
			r, err := NewRegex(test.Pattern, NoCompileOpts, SyntaxDefault)
			if err != nil {
				t.Fatal(err)
			}

			got := r.Matches(test.Str, NoMatchOpts)
			if got != test.Want {
				t.Errorf(
					"wrong Matches result\npattern: %s\nstring:  %s\ngot:     %#v\nwant:    %#v",
					test.Pattern, test.Str, got, test.Want,
				)
			}

			got = r.MatchesBytes([]byte(test.Str), NoMatchOpts)
			if got != test.Want {
				t.Errorf(
					"wrong MatchesBytes result\npattern: %s\nstring:  %s\ngot:     %#v\nwant:    %#v",
					test.Pattern, test.Str, got, test.Want,
				)
			}
		})
	}
}
