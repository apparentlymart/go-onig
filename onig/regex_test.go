package onig

import (
	"fmt"
	"reflect"
	"testing"
)

func TestRegexMatch(t *testing.T) {
	tests := []struct {
		Pattern string
		Str     string
		Want    *Match
	}{
		{
			`hello`,
			`hello world`,
			mustFakeMatch([]Span{
				{0, 5},
			}),
		},
		{
			`hel*o`,
			`helllllo world`,
			mustFakeMatch([]Span{
				{0, 8},
			}),
		},
		{
			`he(l*)o`,
			`helllllo world`,
			mustFakeMatch([]Span{
				{0, 8},
				{2, 7},
			}),
		},
		{
			`he(l*)o`,
			`goodbye world`,
			nil,
		},
		{
			`hello`,
			`why hello, world`,
			nil,
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%q in %q", test.Pattern, test.Str), func(t *testing.T) {
			r, err := NewRegex(test.Pattern, NoCompileOpts, SyntaxRuby)
			if err != nil {
				t.Fatal(err)
			}

			got := r.Match(test.Str, NoMatchOpts)
			if !got.Equal(test.Want) {
				t.Errorf(
					"wrong Match result\npattern: %s\nstring:  %s\ngot:     %#v\nwant:    %#v",
					test.Pattern, test.Str, got, test.Want,
				)
			}

			got = r.MatchBytes([]byte(test.Str), NoMatchOpts)
			if !got.Equal(test.Want) {
				t.Errorf(
					"wrong MatchBytes result\npattern: %s\nstring:  %s\ngot:     %#v\nwant:    %#v",
					test.Pattern, test.Str, got, test.Want,
				)
			}
		})
	}
}

func TestRegexSearch(t *testing.T) {
	tests := []struct {
		Pattern string
		Str     string
		Want    *Match
	}{
		{
			`hello`,
			`hello world`,
			mustFakeMatch([]Span{
				{0, 5},
			}),
		},
		{
			`hel*o`,
			`helllllo world`,
			mustFakeMatch([]Span{
				{0, 8},
			}),
		},
		{
			`he(l*)o`,
			`helllllo world`,
			mustFakeMatch([]Span{
				{0, 8},
				{2, 7},
			}),
		},
		{
			`he(l*)o`,
			`goodbye world`,
			nil,
		},
		{
			`hello`,
			`why hello, world`,
			mustFakeMatch([]Span{
				{4, 9},
			}),
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%q in %q", test.Pattern, test.Str), func(t *testing.T) {
			r, err := NewRegex(test.Pattern, NoCompileOpts, SyntaxRuby)
			if err != nil {
				t.Fatal(err)
			}

			got := r.Search(test.Str, NoMatchOpts)
			if !got.Equal(test.Want) {
				t.Errorf(
					"wrong Search result\npattern: %s\nstring:  %s\ngot:     %#v\nwant:    %#v",
					test.Pattern, test.Str, got, test.Want,
				)
			}

			got = r.SearchBytes([]byte(test.Str), NoMatchOpts)
			if !got.Equal(test.Want) {
				t.Errorf(
					"wrong SearchBytes result\npattern: %s\nstring:  %s\ngot:     %#v\nwant:    %#v",
					test.Pattern, test.Str, got, test.Want,
				)
			}
		})
	}
}

func TestRegexSearchAround(t *testing.T) {
	tests := []struct {
		Pattern    string
		Str        string
		WantBefore string
		WantInside string
		WantAfter  string
	}{
		{
			`hello`,
			`hello world`,
			``, `hello`, ` world`,
		},
		{
			`hel*o`,
			`helllllo world`,
			``, `helllllo`, ` world`,
		},
		{
			`he(l*)o`,
			`helllllo world`,
			``, `helllllo`, ` world`,
		},
		{
			`he(l*)o`,
			`goodbye world`,
			``, ``, `goodbye world`,
		},
		{
			`hello`,
			`why hello, world`,
			`why `, `hello`, `, world`,
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%q in %q", test.Pattern, test.Str), func(t *testing.T) {
			r, err := NewRegex(test.Pattern, NoCompileOpts, SyntaxRuby)
			if err != nil {
				t.Fatal(err)
			}

			{
				gotBefore, gotInside, gotAfter := r.SearchAround(test.Str, NoMatchOpts)
				if gotBefore != test.WantBefore || gotInside != test.WantInside || gotAfter != test.WantAfter {
					t.Errorf(
						"wrong SearchAround result\npattern: %s\nstring:  %s\ngot:     %q %q %q\nwant:    %q %q %q",
						test.Pattern, test.Str,
						gotBefore, gotInside, gotAfter,
						test.WantBefore, test.WantInside, test.WantAfter,
					)
				}
			}
		})
	}
}

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
		{
			`hello`,
			`why hello, world`,
			false,
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%q in %q", test.Pattern, test.Str), func(t *testing.T) {
			r, err := NewRegex(test.Pattern, NoCompileOpts, SyntaxRuby)
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

func TestRegexCaptureCount(t *testing.T) {
	tests := []struct {
		Pattern string
		Want    int
	}{
		{
			`hello`,
			0,
		},
		{
			`hel*o`,
			0,
		},
		{
			`he(l*)o`,
			1,
		},
		{
			`he((l)*)o`,
			2,
		},
	}

	for _, test := range tests {
		t.Run(test.Pattern, func(t *testing.T) {
			r, err := NewRegex(test.Pattern, NoCompileOpts, SyntaxRuby)
			if err != nil {
				t.Fatal(err)
			}

			got := r.CaptureCount()
			if got != test.Want {
				t.Errorf(
					"wrong CaptureCount result\npattern: %s\ngot:     %#v\nwant:    %#v",
					test.Pattern, got, test.Want,
				)
			}
		})
	}
}

func TestRegexNamedCaptures(t *testing.T) {
	tests := []struct {
		Pattern string
		Want    map[string][]int
	}{
		{
			`hello`,
			nil,
		},
		{
			`hel*o`,
			nil,
		},
		{
			`he(l*)o`,
			nil,
		},
		{
			`he((l)*)o`,
			nil,
		},
		{
			`he(?<els>(l)*)o`,
			map[string][]int{
				"els": {1},
			},
		},
		{
			`he((?<el>l)*)o`,
			map[string][]int{
				"el": {1},
			},
		},
		{
			`he(?<els>(?<el>l)*)o`,
			map[string][]int{
				"els": {1},
				"el":  {2},
			},
		},
		{
			`he(?<foo>(?<foo>l)*)o`,
			map[string][]int{
				"foo": {1, 2},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.Pattern, func(t *testing.T) {
			r, err := NewRegex(test.Pattern, NoCompileOpts, SyntaxRuby)
			if err != nil {
				t.Fatal(err)
			}

			got := r.NamedCaptures()
			if !reflect.DeepEqual(got, test.Want) {
				t.Errorf(
					"wrong CaptureCount result\npattern: %s\ngot:     %#v\nwant:    %#v",
					test.Pattern, got, test.Want,
				)
			}
		})
	}
}
