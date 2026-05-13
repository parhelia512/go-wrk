package loader

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

const sessionDeadline = 5 * time.Second

func runSession(t *testing.T, cfg *LoadCfg, ch <-chan *RequesterStats) *RequesterStats {
	t.Helper()
	go cfg.RunSingleLoadSession()
	select {
	case s := <-ch:
		return s
	case <-time.After(sessionDeadline):
		t.Fatalf("RunSingleLoadSession did not return within %v", sessionDeadline)
		return nil
	}
}

func TestRunSingleLoadSession_HappyPath(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	}))
	t.Cleanup(ts.Close)

	ch := make(chan *RequesterStats, 1)
	cfg := NewLoadCfg(1, 1, ts.URL, "", "GET", "", nil, ch, 1000, true, false, false, false, "", "", "", false)

	stats := runSession(t, cfg, ch)

	if stats.NumRequests == 0 {
		t.Fatal("NumRequests = 0, want > 0")
	}
	if stats.NumErrs != 0 {
		t.Errorf("NumErrs = %d, want 0; ErrMap=%v", stats.NumErrs, stats.ErrMap)
	}
	if stats.TotRespSize <= 0 {
		t.Errorf("TotRespSize = %d, want > 0", stats.TotRespSize)
	}
	if stats.TotDuration <= 0 {
		t.Errorf("TotDuration = %v, want > 0", stats.TotDuration)
	}
	if stats.Histogram == nil {
		t.Fatal("Histogram = nil")
	}
	if got := stats.Histogram.TotalCount(); got != int64(stats.NumRequests) {
		t.Errorf("Histogram.TotalCount() = %d, NumRequests = %d", got, stats.NumRequests)
	}
	if len(stats.ErrMap) != 0 {
		t.Errorf("ErrMap not empty: %v", stats.ErrMap)
	}
}
