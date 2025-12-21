package main

import (
	"io"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"
)

func TestRootHandler(t *testing.T) {
	// Save current mode and restore it after test
	originalMode := mode
	defer func() { mode = originalMode }()

	// Set mode to native
	mode = modeNative

	// Initialize the regex pattern if not already done
	if re == nil {
		re = regexp.MustCompile(`((?:\d{1,3}.){3}\d{1,3}):\d+|\[(.+)\]:\d+`)
	}

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

func TestRootHandlerProxyMode(t *testing.T) {
	// Save current mode and restore it after test
	originalMode := mode
	defer func() { mode = originalMode }()

	// Set mode to proxy
	mode = modeProxy

	tests := []struct {
		name           string
		xForwardedFor  string
		expectedIP     string
	}{
		{
			name:          "Single IP",
			xForwardedFor: "203.0.113.195",
			expectedIP:    "203.0.113.195",
		},
		{
			name:          "Multiple IPs - should return first",
			xForwardedFor: "203.0.113.195, 70.41.3.18, 150.172.238.178",
			expectedIP:    "203.0.113.195",
		},
		{
			name:          "IP with whitespace",
			xForwardedFor: "  203.0.113.195  ",
			expectedIP:    "203.0.113.195",
		},
		{
			name:          "Multiple IPs with extra whitespace",
			xForwardedFor: "203.0.113.195 ,  70.41.3.18 , 150.172.238.178",
			expectedIP:    "203.0.113.195",
		},
		{
			name:          "Empty header",
			xForwardedFor: "",
			expectedIP:    "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest("GET", "/", nil)
			if tt.xForwardedFor != "" {
				r.Header.Set("X-Forwarded-For", tt.xForwardedFor)
			}
			w := httptest.NewRecorder()

			rootHandler(w, r)

			resp := w.Result()
			body, _ := io.ReadAll(resp.Body)
			respIP := strings.TrimSuffix(string(body), "\n")

			if respIP != tt.expectedIP {
				t.Errorf("expected %q got %q", tt.expectedIP, respIP)
			}
		})
	}
}
