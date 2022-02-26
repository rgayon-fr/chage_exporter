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
	// Init array of prometheus metric strings
	var prometheusMetrics []PrometheusMetric
	// Get data from chage for every user in relevant users.
	for _, user := range relevantUsers {
		chageData, chageError := GetChage(user)
		// An error here shouldn't crash the software, notify user before proceeding with the next loop.
		if chageError != nil {
			log.Printf("error getting chage for %s: %s", user, chageError)
			continue
		}

		// Return how much time before password change.
		var daysBeforePwExpire int
		timeBeforePwExpire := chageData.UntilExpiryDate()

		// If the password doesn't expire, just return a big negavive number.
		// If it does, just return the actual value.
		if timeBeforePwExpire == 0 {
			daysBeforePwExpire = -999
		} else {
			daysBeforePwExpire = int(chageData.UntilExpiryDate().Hours() / 24)
		}
		prometheusMetrics = append(prometheusMetrics, PrometheusMetric{
			Name:  "chage_user_pw_expire_days",
			Flags: []PrometheusFlag{{Key: "user", Value: user}},
			Value: daysBeforePwExpire,
		})

		// Return the password age
		prometheusMetrics = append(prometheusMetrics, PrometheusMetric{
			Name:  "chage_user_pw_age_days",
			Flags: []PrometheusFlag{{Key: "user", Value: user}},
			Value: int(chageData.PasswordAge().Hours() / 24),
		})
	}

	// For each prometheus metric, inject response with its rendered format.
	for _, metric := range prometheusMetrics {
		renderedMetric, renderError := metric.RenderMetric()
		// Again, an error here shoudn't really hinder the result of the whole thing... warn user and move on.
		if renderError != nil {
			log.Printf("error rendering prometheus metric: %s", renderError)
			continue
		}
		fmt.Fprintln(w, renderedMetric)
	}

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
