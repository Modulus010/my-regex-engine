package main

import (
	"testing"

	"github.com/Modulus010/my-regex-engine/pkg/parse"
)

func TestParseSuccess(t *testing.T) {
	tests := []struct {
		regex string
	}{
		{"^hello$"},
		{"a|b"},
		{"(a|b)c"},
		{"a{3,5}"},
		{"[abc]"},
		{"[^abc]"},
		{".*"},
		{"^a(b|c)*d$"},
	}

	for _, test := range tests {
		t.Run(test.regex, func(t *testing.T) {
			parser := parse.NewParser(test.regex)
			_, err := parser.Parse()
			if err != nil {
				t.Errorf("Expected successful parse, but got error: %v", err)
			}
		})
	}
}

func TestParseErrors(t *testing.T) {
	tests := []struct {
		regex string
	}{
		{"a{"},     // Incomplete quantifier
		{"(a|b"},   // Missing closing parenthesis
		{"[a-c"},   // Missing closing bracket
		{"a{3,2}"}, // Invalid range
		{"*a"},     // Quantifier without preceding element
		{"("},      // Unbalanced parenthesis
		{"[abc"},   // Unbalanced character set
		{"[a-]"},   // Invalid character range
	}

	for _, test := range tests {
		t.Run(test.regex, func(t *testing.T) {
			parser := parse.NewParser(test.regex)
			_, err := parser.Parse()
			if err == nil {
				t.Errorf("Expected parse error, but got success")
			}
		})
	}
}

func TestMatch(t *testing.T) {
	tests := []struct {
		regex    string
		input    string
		expected bool
	}{
		{"^hello$", "hello", true},
		{"^hello$", "hello world", false},
		{"a|b", "a", true},
		{"a|b", "c", false},
		{"a{3}", "aaa", true},
		{"a{3}", "aa", false},
		{"[abc]", "b", true},
		{"[a-c]", "d", false},
		{"(ab)+", "abab", true},
		{".*", "hello world", true},
		{"^a(b|c)*d$", "abcbcd", true},
		{"^a(b|c)*d$", "abce", false},
	}

	for _, test := range tests {
		t.Run(test.regex+"_"+test.input, func(t *testing.T) {
			parser := parse.NewParser(test.regex)
			node, err := parser.Parse()
			if err != nil {
				t.Fatalf("Parse error: %v", err)
			}
			nfa := node.ToNFA()
			result := nfa.Match(test.input)
			if result != test.expected {
				t.Errorf("Expected match: %v, got: %v", test.expected, result)
			}
		})
	}
}

func BenchmarkRegexEngine(b *testing.B) {
	benchmarks := []struct {
		regex string
		input string
	}{
		{"a*", "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"},
		{"(ab)+", "abababababababababababababababababab"},
		{"a{1,100}", "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"},
		{"[a-zA-Z0-9]+", "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"},
		{"^a?{20}a{20}$", "aaaaaaaaaaaaaaaaaaaa"},
		{"^a?{20}a{20}$", "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"},
	}

	for _, bm := range benchmarks {
		b.Run(bm.regex, func(b *testing.B) {
			parser := parse.NewParser(bm.regex)
			node, err := parser.Parse()
			if err != nil {
				b.Fatalf("Parse error: %v", err)
			}
			nfa := node.ToNFA()
			for i := 0; i < b.N; i++ {
				nfa.Match(bm.input)
			}
		})
	}
}
