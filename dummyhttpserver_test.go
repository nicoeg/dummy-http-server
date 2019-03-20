package main

import (
	"net/http"
	"net/url"
	"testing"
)

func TestMatchRoute(t *testing.T) {
	truthy := map[Match]*http.Request{
		Match{"/", "GET"}:               &http.Request{Method: "GET", URL: &url.URL{Path: "/"}},
		Match{"/single", "GET"}:         &http.Request{Method: "GET", URL: &url.URL{Path: "/single"}},
		Match{"/single", "POST"}:        &http.Request{Method: "POST", URL: &url.URL{Path: "/single"}},
		Match{"/first/second", "GET"}:   &http.Request{Method: "GET", URL: &url.URL{Path: "/first/second"}},
		Match{"/first/:second", "GET"}:  &http.Request{Method: "GET", URL: &url.URL{Path: "/first/nope"}},
		Match{"/first/second?", "GET"}:  &http.Request{Method: "GET", URL: &url.URL{Path: "/first"}},
		Match{"/first/secondu?", "GET"}: &http.Request{Method: "GET", URL: &url.URL{Path: "/first/secondu"}},
		Match{"/first/:any?", "GET"}:    &http.Request{Method: "GET", URL: &url.URL{Path: "/first/nope"}},
	}
	falsy := map[Match]*http.Request{
		Match{"/single", "GET"}:        &http.Request{Method: "GET", URL: &url.URL{Path: "/relationship"}},
		Match{"/single", "GET"}:        &http.Request{Method: "POST", URL: &url.URL{Path: "/relationship"}},
		Match{"/single", "POST"}:       &http.Request{Method: "GET", URL: &url.URL{Path: "/relationship"}},
		Match{"/first/second?", "GET"}: &http.Request{Method: "GET", URL: &url.URL{Path: "/first/nope"}},
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
