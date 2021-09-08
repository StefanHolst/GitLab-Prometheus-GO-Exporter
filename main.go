package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"encoding/json"
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
	http.Handle("/metrics", NewGitlabHandler())
	http.ListenAndServe(":5000", nil)
}
