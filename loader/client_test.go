package loader

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/tsliwowicz/go-wrk/util"
)

func TestClient_Default(t *testing.T) {
	c, err := client(false, false, false, 1000, true, "", "", "", false)
	if err != nil {
		t.Fatalf("client() err = %v", err)
	}
	if c == nil {
		t.Fatal("client() returned nil")
	}
	if _, ok := c.Transport.(*http.Transport); !ok {
		t.Fatalf("Transport type = %T, want *http.Transport", c.Transport)
	}
	if c.CheckRedirect != nil {
		t.Errorf("with allowRedirects=true, CheckRedirect should be nil; got non-nil")
	}
}

func TestClient_DisallowRedirects(t *testing.T) {
	c, err := client(false, false, false, 1000, false, "", "", "", false)
	if err != nil {
		t.Fatalf("client() err = %v", err)
	}
	if c.CheckRedirect == nil {
		t.Fatal("with allowRedirects=false, CheckRedirect must be set")
	}

	got := c.CheckRedirect(nil, nil)
	if got == nil {
		t.Fatal("CheckRedirect returned nil error")
	}
	var redirErr *util.RedirectError
	if !errors.As(got, &redirErr) {
		t.Fatalf("CheckRedirect err = %T, want *util.RedirectError", got)
	}
}

func TestClient_SkipVerify(t *testing.T) {
	ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	}))
	t.Cleanup(ts.Close)

	c, err := client(false, false, true, 5000, true, "", "", "", false)
	if err != nil {
		t.Fatalf("client() err = %v", err)
	}
	tr, ok := c.Transport.(*http.Transport)
	if !ok {
		t.Fatalf("Transport type = %T, want *http.Transport", c.Transport)
	}
	if tr.TLSClientConfig == nil || !tr.TLSClientConfig.InsecureSkipVerify {
		t.Fatalf("InsecureSkipVerify should be true, got %+v", tr.TLSClientConfig)
	}

	resp, err := c.Get(ts.URL)
	if err != nil {
		t.Fatalf("GET %s err = %v", ts.URL, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", resp.StatusCode)
	}

	// Sanity: not skipping verification fails against the self-signed cert.
	strict, err := client(false, false, false, 5000, true, "", "", "", false)
	if err != nil {
		t.Fatalf("strict client err = %v", err)
	}
	if _, err := strict.Get(ts.URL); err == nil {
		t.Fatal("strict client should fail on self-signed TLS cert")
	}
}

func TestClient_MissingClientKey(t *testing.T) {
	_, err := client(false, false, false, 1000, true, "cert.pem", "", "ca.pem", false)
	if err == nil {
		t.Fatal("want err for missing client key, got nil")
	}
	if !strings.Contains(err.Error(), "client key") {
		t.Errorf("err = %q, want substring %q", err.Error(), "client key")
	}
}

func TestClient_MissingClientCert(t *testing.T) {
	_, err := client(false, false, false, 1000, true, "", "", "ca.pem", false)
	if err == nil {
		t.Fatal("want err for missing client cert, got nil")
	}
	if !strings.Contains(err.Error(), "client certificate") {
		t.Errorf("err = %q, want substring %q", err.Error(), "client certificate")
	}
}

func TestClient_BadCertFiles(t *testing.T) {
	_, err := client(false, false, false, 1000, true,
		"/nonexistent/cert.pem", "/nonexistent/key.pem", "/nonexistent/ca.pem", false)
	if err == nil {
		t.Fatal("want err for nonexistent cert files, got nil")
	}
	if !strings.Contains(err.Error(), "Unable to load cert") {
		t.Errorf("err = %q, want substring %q", err.Error(), "Unable to load cert")
	}
}
