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

var users []User

func main() {
	fmt.Println(os.Getwd())
	loadConfig()

	// Register users
	for i := 0; i < len(users); i++ {
		registerUser(users[i])
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
	json.Unmarshal(data, &users)
}

func startServer() {
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":5000", nil)
}

func registerUser(user User) {
	prometheus.Register(prometheus.NewGaugeFunc(
		prometheus.GaugeOpts{
			Name:        "user_issue_count",
			Help:        "Number of issues assigned to user.",
			ConstLabels: prometheus.Labels{"username": user.Name, "iid": strconv.Itoa(user.Iid)},
		},
		func() float64 { return GetIssues(user) },
	))
}
