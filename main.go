package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"gopkg.in/yaml.v2"
)

var config Config

func main() {
	fmt.Println(os.Getwd())
	loadConfig()

	// UpdateData(config)

	startServer()
}

func loadConfig() {
	// Open file
	file, err := os.Open("config.yaml")
	if err != nil {
		fmt.Println("Could not open file 'config.json'.")
	}
	defer file.Close()

	// Parse json
	data, _ := ioutil.ReadAll(file)
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		fmt.Println(err)
	}
}

func startServer() {
	http.Handle("/metrics", NewGitlabHandler())
	err := http.ListenAndServe(":5000", nil)
	if err != nil {
		fmt.Println(err)
	}
}
