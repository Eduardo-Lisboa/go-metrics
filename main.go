package main

import (
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	httpRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace:   "",
			Subsystem:   "",
			Name:        "http_requests_total",
			Help:        "Number of requests",
			ConstLabels: nil,
		}, []string{"path"})

	requestsDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "http_requests",
			Name:      "http_requests_duration_seconds",
			Help:      "Duration of requests",
		}, []string{"path"})
)

func init() {
	prometheus.MustRegister(httpRequests)
	prometheus.MustRegister(requestsDuration)
}

func handler(w http.ResponseWriter, r *http.Request) {
	timer := prometheus.NewTimer(requestsDuration.WithLabelValues(r.URL.Path))
	defer timer.ObserveDuration()

	httpRequests.WithLabelValues(r.URL.Path)
	w.Write([]byte("Hello, World!"))
}

func main() {
	http.HandleFunc("/", handler)
	http.Handle("/metrics", promhttp.Handler())

	fmt.Println("Application up and running")
	http.ListenAndServe(":8080", nil)
}
