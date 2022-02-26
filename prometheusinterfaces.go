package main

import (
	"bytes"
	"fmt"
	"text/template"
)

// Basic prometheus Flag
type PrometheusFlag struct {
	Key   string
	Value string
}

// Render the flag array, there's probably a better way to do this.
func (f *PrometheusFlag) RenderFlag() string {
	return fmt.Sprintf("%s=\"%s\"", f.Key, f.Value)
}

// Basic Prometheus Metric structure.
type PrometheusMetric struct {
	Name  string
	Flags []PrometheusFlag
	Value int
}

// Render a prometheus-ready metric.
func (p *PrometheusMetric) RenderMetric() (string, error) {
	// Register metric template.
	metricTemplate, metricTemplateError := template.New("metric").Parse("{{ .Name }}{ {{ range $index, $flag := .Flags }}{{if $index}},{{end}}{{ $flag.RenderFlag }}{{ end}} } {{ .Value }}")
	if metricTemplateError != nil {
		return "", fmt.Errorf("error registering metricTemplate: %s", metricTemplateError)
	}

	// Init template buffer
	var templateBuffer bytes.Buffer
	// Fill the buffer then return the resulting string
	metricTemplate.Execute(&templateBuffer, p)
	return templateBuffer.String(), nil
}
