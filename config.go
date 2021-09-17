package main

import "github.com/prometheus/client_golang/prometheus"

type Config struct {
	Users    []User    `yaml:"users"`
	Projects []Project `yaml:"projects"`
	Servers  []string  `yaml:"servers"`
	Token    string    `yaml:"token"`
}

type User struct {
	UserName                 string `yaml:"username"`
	Name                     string
	MergeRequestsMetric      *prometheus.GaugeVec
	DraftMergeRequestsMetric *prometheus.GaugeVec
}

type Project struct {
	Id     string  `yaml:"id"`
	Labels []Label `yaml:"labels"`
	Metric *prometheus.GaugeVec
}
type Label struct {
	Text  string `yaml:"text"`
	Label string `yaml:"label"`
	Order int    `yaml:"order"`
}
