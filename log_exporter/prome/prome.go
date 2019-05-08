package prome

import (
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	HttpResponseStatus = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "ngxlog_response_status_total",
			Help: "url response status from nginx log.",
		},
		[]string{"url", "code"},
	)

	HttpResponseDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "ngxlog_response_duration_microseconds",
			Help:    "url response duration from nging log",
			Buckets: prometheus.LinearBuckets(100000, 100000, 10), //10 buckets,each 100 millisecond
		},
		[]string{"url"},
	)
)

func init() {
	prometheus.MustRegister(HttpResponseStatus)
	prometheus.MustRegister(HttpResponseDuration)
}

func Start() {
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":9101", nil))
}
