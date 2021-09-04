package main

import (
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	fmt.Println("hej")

	user := User{"stefan", 12}

	fmt.Println(GetIssues(user))

	registerUser(user)
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
			ConstLabels: prometheus.Labels{"username": user.Name, "iid": string(rune(user.Iid))},
		},
		func() float64 { return GetIssues(user) },
	))
}
