package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMain(m *testing.M) {
	log.SetOutput(ioutil.Discard)
	m.Run()
}

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

	router := Router{""}
	for match, request := range truthy {
		if !router._MatchRoute(match, request) {
			t.Errorf("Expected %s to match %s", match.URL, request.URL.Path)
		}
	}

	for match, request := range falsy {
		if router._MatchRoute(match, request) {
			t.Errorf("Expected %s to not be matching %s", match.URL, request.URL.Path)
		}
	}
}

func TestHandleRequest(t *testing.T) {
	t.Run("Handle request success", testHandleRequestSuccess)
	t.Run("Handle request fail", testHandleRequestFail)
}

func testHandleRequestSuccess(t *testing.T) {
	router := Router{"config_test.json"}
	responseWriter := httptest.NewRecorder()
	router.HandleRequest(responseWriter, httptest.NewRequest("GET", "/", nil))

	response := responseWriter.Result()

	if response.StatusCode != 200 {
		t.Errorf("Expected status code to be 200, was %d", response.StatusCode)
	}

	body, _ := ioutil.ReadAll(response.Body)
	if string(body) != "Success" {
		t.Errorf("Expected body to equal 'Success', was %s", string(body))
	}
}

func testHandleRequestFail(t *testing.T) {
	router := Router{"config_test.json"}
	responseWriter := httptest.NewRecorder()
	router.HandleRequest(responseWriter, httptest.NewRequest("GET", "/test", nil))

	response := responseWriter.Result()

	if response.StatusCode != 404 {
		t.Errorf("Expected status code to be 200, was %d", response.StatusCode)
	}

	body, _ := ioutil.ReadAll(response.Body)
	if string(body) != "No request matching" {
		t.Errorf("Expected body to equal 'Success', was %s", string(body))
	}
}
