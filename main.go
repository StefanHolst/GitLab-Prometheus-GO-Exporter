package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"encoding/json"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var config Config

func main() {
	fmt.Println(os.Getwd())
	loadConfig()

	// Register users
	for i := 0; i < len(config.Users); i++ {
		registerUser(config.Users[i])
	}

	// UpdateData()

	startServer()
}

func loadConfig() {
	// Open file
	file, err := os.Open("config.json")
	if err != nil {
		fmt.Println("Could not open file 'config.json'.")
	}
	defer file.Close()

	// Parse json
	data, _ := ioutil.ReadAll(file)
	err = json.Unmarshal(data, &config)
	if err != nil {
		fmt.Println(err)
	}
}

func startServer() {
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":5000", nil)
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
