package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"
)

var configFile string

func main() {
	configFileParam := flag.String("config", "./config.json", "Configuration file")
	portPtr := flag.Int("port", 8080, "port number")
	flag.Parse()
	configFile = *configFileParam
	router := Router{configFile}

	listenStr := fmt.Sprintf(":%d", *portPtr)

	log.Println(time.Now().Format(time.RFC3339), "Listening On", listenStr)
	http.HandleFunc("/", router.HandleRequest)
	http.ListenAndServe(listenStr, nil)
}
