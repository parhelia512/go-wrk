package loader

import (
	"errors"
	"net/http"
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
