package main

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	dto "github.com/prometheus/client_model/go"
)

type GitlabGatherer struct {
}

func (s GitlabGatherer) Gather() ([]*dto.MetricFamily, error) {
	UpdateData(config)
	// config.Users[0].MergeRequestsMetric.WithLabelValues("someone", "something").Set(1)
	return prometheus.DefaultGatherer.Gather()
}

func NewGitlabHandler() http.Handler {
	// Register
	mergeRequestsMetric := register("user_merge_request_count", "Number of merge requests assigned to user.")
	draftMergeRequestsMetric := register("user_draft_merge_request_count", "Number of draft merge requests assigned to user.")

	// Register users
	for i := 0; i < len(config.Users); i++ {
		config.Users[i].MergeRequestsMetric = mergeRequestsMetric
		config.Users[i].DraftMergeRequestsMetric = draftMergeRequestsMetric
	}

	return promhttp.HandlerFor(GitlabGatherer{}, promhttp.HandlerOpts{})
}

func register(name string, help string) *prometheus.GaugeVec {
	gaugeVec := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: name,
			Help: help,
		},
		[]string{"user", "project"},
	)
	prometheus.MustRegister(gaugeVec)
	return gaugeVec

	// user.DraftMergeRequestsMetric = prometheus.NewGaugeVec(
	// 	prometheus.GaugeOpts{
	// 		Name: "user_draft_merge_request_count",
	// 		Help: "Number of draft merge requests assigned to user.",
	// 	},
	// 	[]string{"user", "project"},
	// )
	// prometheus.MustRegister(user.DraftMergeRequestsMetric)

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
