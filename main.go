package main

import (
	"fmt"
	"log"
	"net/http"

	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	// Port number to listen on
	port = kingpin.Flag("port", "Port for chage_exporter to listen on.").Short('p').Default("9200").Int()
)

func handleMetricsRequest(w http.ResponseWriter, r *http.Request) {
	// Get data from chage
	// Use template to send data back to Prometheus
}

func main() {
	// Parse args
	kingpin.Parse()
	// Define HTTP Requests
	http.HandleFunc("/metrics", handleMetricsRequest)

	// Listen for HTTP request
	http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)

	log.Println("chage_exporter")
}
