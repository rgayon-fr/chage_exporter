package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"gopkg.in/alecthomas/kingpin.v2"
	"gopkg.in/yaml.v2"
)

var (
	// Port number to listen on
	port       = kingpin.Flag("port", "Port for chage_exporter to listen on.").Short('p').Default("9200").Int()
	configPath = kingpin.Flag("config", "Path to chage_exporter config.").Short('c').Required().ExistingFile()

	relevantUsers []string
)

func handleMetricsRequest(w http.ResponseWriter, r *http.Request) {
	// Get data from chage
	// Use template to send data back to Prometheus
}

func main() {
	// Parse args
	kingpin.Parse()

	// Parse configuration file.

	// Init configuration structure
	var config Config
	// Read file
	configFileData, configFileError := ioutil.ReadFile(*configPath)
	if configFileError != nil {
		log.Fatalf("error reading the config file: %s", configFileError)
	}
	// Unmarshall YAML data
	parseError := yaml.Unmarshal(configFileData, &config)
	if parseError != nil {
		log.Fatalf("error parsing config file: %s", parseError)
	}

	relevantUsers = config.Users

	// Define HTTP Requests
	http.HandleFunc("/metrics", handleMetricsRequest)

	// Listen for HTTP request
	http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)

	log.Println("chage_exporter")
}
