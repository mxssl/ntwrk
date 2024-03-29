package main

import (
	"io"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRootHandler(t *testing.T) {
	ips := make(map[string]string)
	ips["1.1.1.1"] = "1.1.1.1:65535"
	ips["::1"] = "[::1]:65535"

	for ip, remaddr := range ips {
		r := httptest.NewRequest("GET", "/", nil)
		r.RemoteAddr = remaddr
		w := httptest.NewRecorder()
		rootHandler(w, r)
		resp := w.Result()
		body, _ := io.ReadAll(resp.Body)
		respIP := strings.TrimSuffix(string(body), "\n")
		if respIP != ip {
			t.Errorf("expected \"%v\" got \"%v\"", ip, string(body))
		}
	}
}
