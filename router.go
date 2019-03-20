package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

// Request is a single object in the config file
// Contain route to match and response if matches
type Request struct {
	Match    Match    `json:"match"`
	Response Response `json:"response"`
}

// Match is compared against a HTTP request
// It'll match on URL where wildcards and an optional last route is possible
// It'll also match on Method which should be uppercase
type Match struct {
	URL    string `json:"url"`
	Method string `json:"method"`
}

// Response is what will be returned if a Match is hit
type Response struct {
	Status int    `json:"status"`
	Body   string `json:"body"`
}

// Router will respond to HTTP requests
type Router struct {
	ConfigFile string
}

// HandleRequest compares the given request against route matches in the configuration file
// If one matches it stops and responds to the request with the response from config match
// If no match is found 404 response is returned
func (router Router) HandleRequest(response http.ResponseWriter, request *http.Request) {
	for _, element := range router._ReadConfigFile() {
		if router._MatchRoute(element.Match, request) {
			response.WriteHeader(element.Response.Status)
			fmt.Fprint(response, element.Response.Body)
			fmt.Println(element.Response.Body, request.Method, element.Response.Status, request.URL, request.UserAgent())
			return
		}
	}

	response.WriteHeader(404)
	fmt.Fprint(response, "No request matching")
	fmt.Println(time.Now().Format(time.RFC3339), request.Method, 404, request.URL, request.UserAgent())
}

func (router Router) _ReadConfigFile() []Request {
	jsonFile, err := os.Open(router.ConfigFile)

	if err != nil {
		fmt.Println("Could not open configuration file.")
		os.Exit(0)
	}

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	var requests []Request
	json.Unmarshal(byteValue, &requests)

	return requests
}

// _MatchRoute will check multiple route parameters from match against the incoming request
// It returns true if everything matches
func (router Router) _MatchRoute(match Match, request *http.Request) bool {
	if !router._MatchURL(match.URL, request.URL.Path) {
		return false
	}

	if match.Method != "" && match.Method != request.Method {
		return false
	}

	return true
}

// _MatchURL will check if the matchURL matches requestURL
// It handles wildcards in case of : and optional last routes with ?
func (router Router) _MatchURL(matchURL string, requestURL string) bool {
	matchParts := strings.Split(strings.TrimRight(matchURL, "/"), "/")
	requestParts := strings.Split(strings.TrimRight(requestURL, "/"), "/")
	for index, part := range matchParts {
		if strings.HasSuffix(part, "?") && index == len(requestParts) {
			continue
		}

		if index == len(matchParts)-1 && len(matchParts) != len(requestParts) {
			return false
		}

		if index >= len(requestParts) {
			return false
		}

		if strings.HasPrefix(part, ":") {
			continue
		}

		if strings.Replace(part, "?", "", -1) != requestParts[index] {
			return false
		}
	}

	return true
}
