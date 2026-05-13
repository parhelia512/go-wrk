package util

import (
	"net/http"
	"strings"
	"testing"
	"time"
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

func TestMaxDuration(t *testing.T) {
	a := 100 * time.Millisecond
	b := 200 * time.Millisecond

	if got := MaxDuration(a, b); got != b {
		t.Errorf("MaxDuration(a,b) = %v, want %v", got, b)
	}
	if got := MaxDuration(b, a); got != b {
		t.Errorf("MaxDuration(b,a) = %v, want %v", got, b)
	}
	if got := MaxDuration(a, a); got != a {
		t.Errorf("MaxDuration(a,a) = %v, want %v", got, a)
	}
}

func TestMinDuration(t *testing.T) {
	a := 100 * time.Millisecond
	b := 200 * time.Millisecond

	if got := MinDuration(a, b); got != a {
		t.Errorf("MinDuration(a,b) = %v, want %v", got, a)
	}
	if got := MinDuration(b, a); got != a {
		t.Errorf("MinDuration(b,a) = %v, want %v", got, a)
	}
	if got := MinDuration(a, a); got != a {
		t.Errorf("MinDuration(a,a) = %v, want %v", got, a)
	}
}

func TestEstimateHttpHeadersSize(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		got := EstimateHttpHeadersSize(http.Header{})
		want := int64(len("\r\n"))
		if got != want {
			t.Fatalf("empty headers size = %d, want %d", got, want)
		}
	})

	t.Run("single_header_single_value", func(t *testing.T) {
		h := http.Header{}
		h.Add("X-Foo", "bar")
		// per implementation: len(k) + len(": \r\n") + len(v) for each value, + trailing "\r\n"
		want := int64(len("X-Foo") + len(": \r\n") + len("bar") + len("\r\n"))
		if got := EstimateHttpHeadersSize(h); got != want {
			t.Fatalf("size = %d, want %d", got, want)
		}
	})

	t.Run("single_header_multi_value", func(t *testing.T) {
		h := http.Header{}
		h.Add("X-Foo", "a")
		h.Add("X-Foo", "bb")
		want := int64(len("X-Foo") + len(": \r\n") + len("a") + len("bb") + len("\r\n"))
		if got := EstimateHttpHeadersSize(h); got != want {
			t.Fatalf("size = %d, want %d", got, want)
		}
	})

	t.Run("multiple_headers", func(t *testing.T) {
		h := http.Header{}
		h.Add("A", "1")
		h.Add("BB", "22")
		// Map iteration order doesn't matter — addition is commutative.
		want := int64(
			len("A") + len(": \r\n") + len("1") +
				len("Bb") + len(": \r\n") + len("22") + // canonical form: "Bb"
				len("\r\n"))
		if got := EstimateHttpHeadersSize(h); got != want {
			t.Fatalf("size = %d, want %d", got, want)
		}
	})
}

func TestRedirectError(t *testing.T) {
	const msg = "no redirects please"
	err := NewRedirectError(msg)
	if err == nil {
		t.Fatal("NewRedirectError returned nil")
	}
	if err.Error() != msg {
		t.Errorf("err.Error() = %q, want %q", err.Error(), msg)
	}

	// Confirm it satisfies the error interface.
	var _ error = err

	// Avoid the unused import lint if strings becomes unused.
	_ = strings.Repeat
}
