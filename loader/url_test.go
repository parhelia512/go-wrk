package loader

import (
	"testing"
)

func TestEscapeUrlStr(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want string
	}{
		{"no_query", "http://x/path", "http://x/path"},
		{"single_param_with_space", "http://x/path?a=hello world", "http://x/path?a=hello+world"},
		{"multi_param", "http://x/p?a=1&b=hello world&c=x", "http://x/p?a=1&b=hello+world&c=x"},
		{"bare_flag_then_param", "http://x/p?flag&a=1", "http://x/p?flag&a=1"},
		{"empty_query", "http://x/p?", "http://x/p?"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := escapeUrlStr(tc.in); got != tc.want {
				t.Fatalf("escapeUrlStr(%q) = %q, want %q", tc.in, got, tc.want)
			}
		})
	}
}
