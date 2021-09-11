package main

import "github.com/prometheus/client_golang/prometheus"

type Config struct {
	Users    []User    `json:"users"`
	Projects []Project `json:"projects"`
	Token    string    `json:"token"`
}

type User struct {
	UserName                 string `json:"username"`
	Name                     string
	MergeRequestsMetric      *prometheus.GaugeVec
	DraftMergeRequestsMetric *prometheus.GaugeVec
}

type Project struct {
	Id     string   `json:"Id"`
	Labels []string `json:"labels"`
	Metric *prometheus.GaugeVec
}
