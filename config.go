package main

import "github.com/prometheus/client_golang/prometheus"

type Config struct {
	Users []User `json:"users"`
	Token string `json:"token"`
}

type User struct {
	UserName                 string `json:"username"`
	Name                     string
	MergeRequestsMetric      *prometheus.GaugeVec
	DraftMergeRequestsMetric *prometheus.GaugeVec
}
