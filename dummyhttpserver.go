package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

type Request struct {
	Match    Match    `json:"match"`
	Response Response `json:"response"`
}

type Match struct {
	URL    string `json:"url"`
	Method string `json:"method"`
}

type Response struct {
	Status int    `json:"status"`
	Body   string `json:"body"`
}

var configFile string

func main() {
	configFileParam := flag.String("config", "./config.json", "Configuration file")
	portPtr := flag.Int("port", 8080, "port number")
	flag.Parse()
	configFile = *configFileParam

	listenStr := fmt.Sprintf(":%d", *portPtr)

	fmt.Println(time.Now().Format(time.RFC3339), "Listening On", listenStr)
	http.HandleFunc("/", handleRequest)
	http.ListenAndServe(listenStr, nil)
}

func handleRequest(response http.ResponseWriter, request *http.Request) {
	for _, element := range readConfigFile() {
		if matchRoute(element.Match, request) {
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

func matchRoute(match Match, request *http.Request) bool {
	if !matchURL(match.URL, request.URL.Path) {
		return false
	}

	if match.Method != "" && match.Method != request.Method {
		return false
	}

	return true
}

func matchURL(matchURL string, requestURL string) bool {
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

func readConfigFile() []Request {
	jsonFile, err := os.Open(configFile)

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
