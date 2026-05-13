package util

import (
	"testing"
)

func TestByteSizeString(t *testing.T) {
	cases := []struct {
		name string
		size float64
		want string
	}{
		{"bytes_small", 0, "0.00bytes"},
		{"bytes_just_under_kb", 1024, "1024.00bytes"},
		{"kb", 1025, "1.00KB"},
		{"kb_mid", 2048, "2.00KB"},
		{"mb", 2 * 1024 * 1024, "2.00MB"},
		{"gb", 3 * 1024 * 1024 * 1024, "3.00GB"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := ByteSize{Size: tc.size}.String()
			if got != tc.want {
				t.Fatalf("ByteSize{%v}.String() = %q, want %q", tc.size, got, tc.want)
			}
		})
	}
}
