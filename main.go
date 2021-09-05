package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

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
	prometheus.Register(prometheus.NewCounterFunc(
		prometheus.CounterOpts{
			Name:        "user_issue_count",
			Help:        "Number of issues assigned to user.",
			ConstLabels: prometheus.Labels{"username": user.Name, "iid": strconv.Itoa(user.Iid)},
		},
		func() float64 { return GetIssues(user) },
	))
}
