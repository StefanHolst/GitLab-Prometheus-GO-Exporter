package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
)

func GetIssues(user User) float64 { //gid://gitlab/User/1880311
	query := `{"query": 
	"query {
		user(username: "StefanHolst0") {
			name
		}
	}
	"}`
	query = strings.Replace(query, "\n", "", -1)
	query = strings.Replace(query, "\t", "", -1)

	data := bytes.NewBufferString(query)
	client := &http.Client{}
	request, err := http.NewRequest("POST", "https://gitlab.com/api/graphql", data)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Authorization", "Bearer 1VzXksMHksrANMueszWT")

	response, err := client.Do(request)
	if err != nil {
		fmt.Println(err)
	}

	var values map[string]interface{}
	json.NewDecoder(response.Body).Decode(&values)

	fmt.Println(values)

	return rand.Float64()
}
