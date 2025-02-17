package handlers_test

import (
	"bytes"
	"itu-minitwit/setup"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPingHandler(t *testing.T) {
	r := setup.SetupTest()

	req, _ := http.NewRequest("GET", "/ping", bytes.NewBufferString(""))
	req.Header.Set("Content-Type", "text/plain")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK || w.Body.String() != "pong" {
		t.Fatalf("Test failed with status code %v and body %v", w.Code, w.Body.String())
	}
}
