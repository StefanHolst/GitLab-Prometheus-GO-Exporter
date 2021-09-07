package main

import (
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	dto "github.com/prometheus/client_model/go"
)

var config Config

type GitlabGatherer struct {
}

func (s GitlabGatherer) Gather() ([]*dto.MetricFamily, error) {
	// is this the event we are looking for?
	fmt.Println("hej")
	UpdateData(config)
	return prometheus.DefaultGatherer.Gather()
}

func NewGitlabHandler() http.Handler {

	// Register users
	for i := 0; i < len(config.Users); i++ {
		registerUser(config.Users[i])
	}

	return promhttp.HandlerFor(GitlabGatherer{}, promhttp.HandlerOpts{})
}

func registerUser(user User) {
	user.MergeRequestsMetric = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "user_merge_request_count",
			Help: "Number of merge requests assigned to user.",
		},
		[]string{"user", "project"},
	)
	prometheus.Register(user.MergeRequestsMetric)

	user.DraftMergeRequestsMetric = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "user_draft_merge_request_count",
			Help: "Number of draft merge requests assigned to user.",
		},
		[]string{"user", "project"},
	)
	prometheus.Register(user.DraftMergeRequestsMetric)

	// prometheus.Register(prometheus.NewCounterFunc(
	// 	prometheus.CounterOpts{
	// 		Name:        "user_merge_request_count",
	// 		Help:        "Number of merge requests assigned to user.",
	// 		ConstLabels: prometheus.Labels{"username": user.Name, "name": user.Name},
	// 	},
	// 	func() float64 { return float64(user.MergeRequests) },
	// ))

	// prometheus.Register(prometheus.NewCounterFunc(
	// 	prometheus.CounterOpts{
	// 		Name:        "user_draft_merge_request_count",
	// 		Help:        "Number of draft merge requests assigned to user.",
	// 		ConstLabels: prometheus.Labels{"username": user.Name, "name": user.Name},
	// 	},
	// 	func() float64 { return float64(user.DraftMergeRequests) },
	// ))
}
