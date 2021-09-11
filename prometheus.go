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
	mergeRequestsMetric := registerGauge("user_merge_request_count", "Number of merge requests assigned to user.", []string{"user", "project"})
	draftMergeRequestsMetric := registerGauge("user_draft_merge_request_count", "Number of draft merge requests assigned to user.", []string{"user", "project"})
	projectMetric := registerGauge("project_board", "Project Board", []string{"project", "label", "order"})

	// Add metric to users
	for i := 0; i < len(config.Users); i++ {
		config.Users[i].MergeRequestsMetric = mergeRequestsMetric
		config.Users[i].DraftMergeRequestsMetric = draftMergeRequestsMetric
	}

	// Add metric to projects
	for i := 0; i < len(config.Projects); i++ {
		config.Projects[i].Metric = projectMetric
	}

	return promhttp.HandlerFor(GitlabGatherer{}, promhttp.HandlerOpts{})
}

func registerGauge(name string, help string, labels []string) *prometheus.GaugeVec {
	gaugeVec := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: name,
			Help: help,
		},
		labels,
	)
	prometheus.MustRegister(gaugeVec)
	return gaugeVec
}
