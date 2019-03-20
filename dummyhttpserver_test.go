package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMatchRoute(t *testing.T) {
	truthy := map[Match]*http.Request{
		Match{"/", "GET"}:               httptest.NewRequest("GET", "/", nil),
		Match{"/single", "GET"}:         httptest.NewRequest("GET", "/single", nil),
		Match{"/single", "POST"}:        httptest.NewRequest("POST", "/single", nil),
		Match{"/first/second", "GET"}:   httptest.NewRequest("GET", "/first/second", nil),
		Match{"/first/:second", "GET"}:  httptest.NewRequest("GET", "/first/nope", nil),
		Match{"/first/second?", "GET"}:  httptest.NewRequest("GET", "/first", nil),
		Match{"/first/secondu?", "GET"}: httptest.NewRequest("GET", "/first/secondu", nil),
		Match{"/first/:any?", "GET"}:    httptest.NewRequest("GET", "/first/nope", nil),
	}
	falsy := map[Match]*http.Request{
		Match{"/single", "GET"}:        httptest.NewRequest("GET", "/relationship", nil),
		Match{"/single", "GET"}:        httptest.NewRequest("POST", "/relationship", nil),
		Match{"/single", "POST"}:       httptest.NewRequest("GET", "/relationship", nil),
		Match{"/first/second?", "GET"}: httptest.NewRequest("GET", "/first/nope", nil),
	}

	for match, request := range truthy {

		if !matchRoute(match, request) {
			t.Errorf("Expected %s to match %s", match.URL, request.URL.Path)
		}
	}

	for match, request := range falsy {
		if matchRoute(match, request) {
			t.Errorf("Expected %s to not be matching %s", match.URL, request.URL.Path)
		}
	}
}
