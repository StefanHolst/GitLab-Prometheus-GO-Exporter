package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"

	"github.com/tidwall/gjson"
)

func UpdateData() {
	// TODO: Get upcoming milestones for configured projects

	// Format all usernames
	var usernames []string
	for i := 0; i < len(config.Users); i++ {
		usernames = append(usernames, "\\\""+config.Users[i].UserName+"\\\"")
	}

	// Prepare request
	payload := strings.NewReader("{\"query\":\"query {\\n    users(usernames: [" + strings.Join(usernames, ",") + "]) {\\n        nodes{\\n            assignedMergeRequests(first:100, state: opened){\\n                nodes{\\n                    title\\n                    draft\\n                    project{\\n                        name\\n                    }\\n                    milestone{\\n                        title\\n                        dueDate\\n                    }\\n                }\\n            }\\n            name\\n        }\\n    }\\n}\",\"variables\":{}}")
	client := &http.Client{}
	request, err := http.NewRequest("POST", "https://gitlab.com/api/graphql", payload)
	if err != nil {
		fmt.Println(err)
		return
	}
	request.Header.Add("Authorization", "Bearer "+config.Token)
	request.Header.Add("Content-Type", "application/json")

	// Send request
	res, err := client.Do(request)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	graphUsers := gjson.GetBytes(body, "data.users.nodes").Array()

	for _, graphUser := range graphUsers {
		fmt.Println(graphUser.Get("name").Value())
		user := getUser(graphUser.Get("name").String())
		if user.UserName == "" {
			continue
		}

		var mergeRequests []gjson.Result
		var draftMergeRequests []gjson.Result

		// Count all mergerequest for user
		graphMergeRequests := graphUser.Get("assignedMergeRequests").Array()
		for _, graphMergeRequest := range graphMergeRequests {
			if graphMergeRequest.Get("draft").Bool() {
				draftMergeRequests = append(draftMergeRequests, graphMergeRequest)
			} else {
				mergeRequests = append(mergeRequests, graphMergeRequest)
			}
		}
	}
}

func getUser(username string) User {
	for _, user := range config.Users {
		if user.UserName == username {
			return user
		}
	}
	return User{}
}

func GetIssues(user User) float64 {
	query := ""

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
