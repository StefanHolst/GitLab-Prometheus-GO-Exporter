package main

import (
	"math/rand"

	"github.com/prometheus/client_golang/prometheus"
)

type GitlabCollector struct {
	MergeRequests     *prometheus.Desc
	DraftMergeRequest *prometheus.Desc
}

func newGitlabCollector() *GitlabCollector {
	return &GitlabCollector{
		MergeRequests: prometheus.NewDesc("user_merge_requests_count",
			"Number of merge requests assigned to user",
			nil, nil,
		),
		DraftMergeRequest: prometheus.NewDesc("user_draft_merge_requests_count",
			"Number of merge requests assigned to user",
			nil, nil,
		),
	}
}

func (collector *GitlabCollector) Describe(ch chan<- *prometheus.Desc) {

	//Update this section with the each metric you create for a given collector
	// ch <- collector.MergeRequests
	// ch <- collector.DraftMergeRequest
}

func (collector *GitlabCollector) Collect(ch chan<- prometheus.Metric) {

	//Implement logic here to determine proper metric value to return to prometheus
	//for each descriptor or call other functions that do so.
	var metricValue float64
	metricValue = rand.Float64()

	//Write latest value for each metric in the prometheus metric channel.
	//Note that you can pass CounterValue, GaugeValue, or UntypedValue types here.
	ch <- prometheus.MustNewConstMetric(collector.MergeRequests, prometheus.CounterValue, metricValue)
	ch <- prometheus.MustNewConstMetric(collector.DraftMergeRequest, prometheus.CounterValue, metricValue)
}
