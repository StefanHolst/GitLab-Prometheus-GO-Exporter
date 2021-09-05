package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	fmt.Println("hej")

	user := User{"stefan", 12}

	fmt.Println(GetIssues(user))

	registerUser(user)
	registerUser(User{"someone", 10})
	startServer()
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
