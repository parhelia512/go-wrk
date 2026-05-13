package util

import (
	"testing"
)

func TestHeaderList_SetAndString(t *testing.T) {
	var hl HeaderList

	if got := hl.String(); got != "" {
		t.Fatalf("empty HeaderList.String() = %q, want %q", got, "")
	}

	if err := hl.Set("A: 1"); err != nil {
		t.Fatalf("Set returned err: %v", err)
	}
	if got := hl.String(); got != "A: 1" {
		t.Fatalf("after one Set, String() = %q, want %q", got, "A: 1")
	}

	_ = hl.Set("B: 2")
	_ = hl.Set("C: 3")
	if got := hl.String(); got != "A: 1, B: 2, C: 3" {
		t.Fatalf("after three Sets, String() = %q, want %q", got, "A: 1, B: 2, C: 3")
	}

	if len(hl) != 3 {
		t.Fatalf("len(hl) = %d, want 3", len(hl))
	}
}

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
