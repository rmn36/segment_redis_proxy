package main

import (
	"log"
	"net/http"
	"fmt"
	"os"
	"strconv"
)

var max_connections = 10
var cache_size = 10
var cache_expiry = 60000
var redis_url = "redis:6379"
var port = "8080"

func aboutHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	
	fmt.Fprintf(w, "Golang Redis proxy for Segment Interview by Ryan Nowacoski\n")
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "https://github.com/rmn36/", 301)
}

func setHandler(w http.ResponseWriter, r *http.Request) {
}

func getHandler(w http.ResponseWriter, r *http.Request) {
}

func readConfigs(){
	var conn_parse_err error
	max_connections, conn_parse_err := strconv.Atoi(os.Getenv("MAX_CONNECTIONS"))
	if conn_parse_err != nil {
		panic("MAX_CONNECTIONS MUST BE A NUMBER")
	}
	log.Println("MAX CONNECTIONS: "+strconv.Itoa(max_connections))

	var cache_parse_err error
	cache_size, cache_parse_err := strconv.Atoi(os.Getenv("CACHE_SIZE"))
	if cache_parse_err != nil {
		panic("CACHE_SIZE MUST BE A NUMBER")
	}
	log.Println("CACHE SIZE: "+strconv.Itoa(cache_size))

	var cache_exp_parse_err error
	cache_expiry, cache_exp_parse_err = strconv.Atoi(os.Getenv("CACHE_EXPIRY_TIME"))
	if cache_exp_parse_err != nil {
		panic("CACHE_EXPIRY_TIME MUST BE A NUMBER")
	}
	log.Println("CACHE EXPIRY TIME (ms): "+strconv.Itoa(cache_expiry))

	redis_url = os.Getenv("REDIS_URL")
	log.Println("REDIS URL: "+redis_url)
	port = os.Getenv("PORT")
	log.Println("LISTENING PORT: "+port)
}

func main() {
	readConfigs()

	http.HandleFunc("/set", setHandler)
	http.HandleFunc("/get", getHandler)
	http.HandleFunc("/about", aboutHandler)
	http.HandleFunc("/", rootHandler)

	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		panic(err)
	}

	log.Println("Server ended on port "+port)
}