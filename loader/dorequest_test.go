package loader

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newTestServer(t *testing.T, h http.HandlerFunc) *httptest.Server {
	t.Helper()
	ts := httptest.NewServer(h)
	t.Cleanup(ts.Close)
	return ts
}

func defaultTestClient(t *testing.T) *http.Client {
	t.Helper()
	c, err := client(false, false, false, 5000, true, "", "", "", false)
	if err != nil {
		t.Fatalf("client() err = %v", err)
	}
	return c
}

func TestDoRequest_Success200(t *testing.T) {
	ts := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("User-Agent"); got != USER_AGENT {
			t.Errorf("User-Agent = %q, want %q", got, USER_AGENT)
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("hello"))
	})

	respSize, dur, err := DoRequest(defaultTestClient(t), nil, "GET", "", ts.URL, "")
	if err != nil {
		t.Fatalf("DoRequest err = %v", err)
	}
	if respSize <= 0 {
		t.Errorf("respSize = %d, want > 0", respSize)
	}
	if dur <= 0 {
		t.Errorf("duration = %v, want > 0", dur)
	}
}

func TestDoRequest_CustomHeadersAndHost(t *testing.T) {
	ts := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Host != "myvhost" {
			t.Errorf("r.Host = %q, want %q", r.Host, "myvhost")
		}
		if got := r.Header.Get("X-Foo"); got != "bar" {
			t.Errorf("X-Foo = %q, want %q", got, "bar")
		}
		w.WriteHeader(http.StatusOK)
	})

	hdr := map[string]string{"X-Foo": "bar"}
	_, _, err := DoRequest(defaultTestClient(t), hdr, "GET", "myvhost", ts.URL, "")
	if err != nil {
		t.Fatalf("DoRequest err = %v", err)
	}
}

func TestDoRequest_Body(t *testing.T) {
	ts := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		if string(body) != "payload" {
			t.Errorf("server got body %q, want %q", string(body), "payload")
		}
		if r.Method != "POST" {
			t.Errorf("r.Method = %q, want POST", r.Method)
		}
		w.WriteHeader(http.StatusOK)
	})

	_, _, err := DoRequest(defaultTestClient(t), nil, "POST", "", ts.URL, "payload")
	if err != nil {
		t.Fatalf("DoRequest err = %v", err)
	}
}

func TestDoRequest_QueryEscaping(t *testing.T) {
	ts := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.Query().Get("q"); got != "hello world" {
			t.Errorf("q = %q, want %q", got, "hello world")
		}
		w.WriteHeader(http.StatusOK)
	})

	_, _, err := DoRequest(defaultTestClient(t), nil, "GET", "", ts.URL+"/?q=hello world", "")
	if err != nil {
		t.Fatalf("DoRequest err = %v", err)
	}
}

func TestDoRequest_301RedirectBlocked(t *testing.T) {
	ts := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Location", "http://example.invalid/")
		w.WriteHeader(http.StatusMovedPermanently)
	})

	// Default client: allowRedirects=false → CheckRedirect returns RedirectError.
	c, err := client(false, false, false, 5000, false, "", "", "", false)
	if err != nil {
		t.Fatalf("client() err = %v", err)
	}

	respSize, dur, err := DoRequest(c, nil, "GET", "", ts.URL, "")
	if err == nil {
		t.Fatalf("want err for blocked redirect, got nil (respSize=%d dur=%v)", respSize, dur)
	}
	if respSize != 0 {
		t.Errorf("respSize = %d, want 0", respSize)
	}
	if dur != 0 {
		t.Errorf("dur = %v, want 0", dur)
	}
}

// keepLastResponseClient returns an *http.Client that surfaces redirect responses
// without following them, so DoRequest's 301/307 branch is reachable.
func keepLastResponseClient() *http.Client {
	return &http.Client{
		Transport: &http.Transport{},
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
}

func TestDoRequest_301RedirectAsResponse(t *testing.T) {
	ts := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Location", "http://example.invalid/")
		w.WriteHeader(http.StatusMovedPermanently)
	})

	respSize, dur, err := DoRequest(keepLastResponseClient(), nil, "GET", "", ts.URL, "")
	if err != nil {
		t.Fatalf("DoRequest err = %v", err)
	}
	if respSize <= 0 {
		t.Errorf("respSize = %d, want > 0", respSize)
	}
	if dur <= 0 {
		t.Errorf("dur = %v, want > 0", dur)
	}
}

func TestDoRequest_307RedirectAsResponse(t *testing.T) {
	ts := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Location", "http://example.invalid/")
		w.WriteHeader(http.StatusTemporaryRedirect)
	})

	respSize, dur, err := DoRequest(keepLastResponseClient(), nil, "GET", "", ts.URL, "")
	if err != nil {
		t.Fatalf("DoRequest err = %v", err)
	}
	if respSize <= 0 {
		t.Errorf("respSize = %d, want > 0", respSize)
	}
	if dur <= 0 {
		t.Errorf("dur = %v, want > 0", dur)
	}
}
