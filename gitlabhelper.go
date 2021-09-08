package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/tidwall/gjson"
)

func UpdateData(config Config) {
	// Format all username
	var usernameQueries []string
	for i := 0; i < len(config.Users); i++ {
		usernameQueries = append(usernameQueries, "\\\""+config.Users[i].UserName+"\\\"")
	}

	// Download all users mergerequests
	payload := strings.NewReader("{\"query\":\"query {\\n    users(usernames: [" + strings.Join(usernameQueries, ",") + "]) {\\n        nodes{\\n            assignedMergeRequests(first:100, state: opened){\\n                nodes{\\n                    title\\n                    draft\\n                    project{\\n                        id\\n                        name\\n                    }\\n                    milestone{\\n                        title\\n                        dueDate\\n                    }\\n                }\\n            }\\n            name\\n            username\\n        }\\n    }\\n}\",\"variables\":{}}")
	data, err := downloadData(payload)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Get all users
	graphUsers := data.Get("users.nodes").Array()

	// Get a map of all project ids
	projectUpcomingMilestones := getAllProjectFromUsers(data)

	// Format all projects
	var projectQueries []string
	for key := range projectUpcomingMilestones {
		projectQueries = append(projectQueries, "\\\""+key+"\\\"")
	}

	// Download milestones for all projects
	payload2 := strings.NewReader("{\"query\":\"query {\\n    projects(ids:[" + strings.Join(projectQueries, ",") + "]) {\\n        nodes{\\n            id\\n            milestones(state:active, sort: DUE_DATE_DESC,first:1){\\n                nodes{\\n                    id\\n                    title\\n                }\\n            }\\n        }\\n    }\\n}\",\"variables\":{}}")
	data2, err := downloadData(payload2)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Set upcoming milestone for each project
	setUpcomingMilestoneForProjects(data2, projectUpcomingMilestones)

	// filter MRs by milestone == null or milestone == project.upcomingMilestone
	for _, graphUser := range graphUsers {
		fmt.Println(graphUser.Get("name").Value())
		user := getUser(graphUser.Get("username").String())
		if user.UserName == "" {
			continue
		}

		var mergeRequests []gjson.Result
		var draftMergeRequests []gjson.Result

		// Get all mergerequest for user
		graphMergeRequests := graphUser.Get("assignedMergeRequests.nodes").Array()

		for _, graphMergeRequest := range graphMergeRequests {
			if graphMergeRequest.Get("draft").Bool() {
				draftMergeRequests = append(draftMergeRequests, graphMergeRequest)
			} else {
				mergeRequests = append(mergeRequests, graphMergeRequest)
			}
		}

		// group mergerequest by project
		// set metric
		user.MergeRequestsMetric.WithLabelValues(user.Name, "something").Set(float64(len(mergeRequests)))
		user.DraftMergeRequestsMetric.WithLabelValues(user.Name, "something").Set(float64(len(draftMergeRequests)))
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

func setUpcomingMilestoneForProjects(data3 gjson.Result, projectUpcomingMilestones map[string]string) {
	// Update milestone map with milestone id
	projectIds := data3.Get("projects.nodes").Array()
	for _, projectId := range projectIds {
		id := projectId.Get("id").String()
		milestoneId := projectId.Get("milestones.nodes.0.id").String()
		projectUpcomingMilestones[id] = milestoneId
	}
}

func getAllProjectFromUsers(data gjson.Result) map[string]string {
	projectMap := make(map[string]string)

	users := data.Get("users.nodes").Array()
	for _, user := range users {
		projects := user.Get("assignedMergeRequests.nodes").Array()
		for _, project := range projects {
			id := project.Get("project.id").String()
			projectMap[id] = ""
		}
	}

	// var projectIds []string
	// for key := range projectMap {
	// 	projectIds = append(projectIds, key)
	// }

	return projectMap
}

func downloadData(payload *strings.Reader) (gjson.Result, error) {
	client := &http.Client{}
	request, err := http.NewRequest("POST", "https://gitlab.com/api/graphql", payload)
	if err != nil {
		fmt.Println(err)
		return gjson.Result{}, err
	}
	request.Header.Add("Authorization", "Bearer "+config.Token)
	request.Header.Add("Content-Type", "application/json")

	// Send request
	res, err := client.Do(request)
	if err != nil {
		fmt.Println(err)
		return gjson.Result{}, err
	}
	defer res.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return gjson.Result{}, err
	}

	return gjson.GetBytes(body, "data"), nil
}
