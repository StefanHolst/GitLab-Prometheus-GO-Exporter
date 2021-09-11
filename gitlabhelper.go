package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/tidwall/gjson"
)

func UpdateData(config Config) {
	fmt.Println("Updating " + time.Now().Format("Mon Jan _2 2006 15:04:05"))

	// Format all username
	var usernameQueries []string
	for i := 0; i < len(config.Users); i++ {
		usernameQueries = append(usernameQueries, "\\\""+config.Users[i].UserName+"\\\"")
	}

	// Download all users mergerequests
	payload := strings.NewReader("{\"query\":\"query {\\n    users(usernames: [" + strings.Join(usernameQueries, ",") + "]) {\\n        nodes{\\n            assignedMergeRequests(first:100, state: opened){\\n                nodes{\\n                    title\\n                    draft\\n                    project{\\n                        id\\n                        name\\n                    }\\n                    milestone{\\n                        id\\n                    }\\n                }\\n            }\\n            name\\n            username\\n        }\\n    }\\n}\",\"variables\":{}}")
	data, err := downloadData(payload)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Get all users
	graphUsers := data.Get("users.nodes").Array()

	// Get a map of all project ids
	projectUpcomingMilestones := getAllProjectFromUsers(data)

	// Format all user project ids
	var userProjectQueries []string
	for key := range projectUpcomingMilestones {
		userProjectQueries = append(userProjectQueries, "\\\""+key+"\\\"")
	}

	// Download milestones for all project ids
	payload = strings.NewReader("{\"query\":\"query {\\n    projects(ids:[" + strings.Join(userProjectQueries, ",") + "]) {\\n        nodes{\\n            id\\n            milestones(state:active, sort: DUE_DATE_DESC,first:1){\\n                nodes{\\n                    id\\n                    title\\n                }\\n            }\\n        }\\n    }\\n}\",\"variables\":{}}")
	data, err = downloadData(payload)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Set upcoming milestone for each project
	setUpcomingMilestoneForProjects(data, projectUpcomingMilestones)

	// filter MRs by milestone == null or milestone == project.upcomingMilestone
	for _, graphUser := range graphUsers {
		fmt.Println("- " + graphUser.Get("name").String())
		user, err := getUser(graphUser.Get("username").String())
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		user.Name = graphUser.Get("name").String()

		mergeRequests := make(map[string]int)
		draftMergeRequests := make(map[string]int)

		// Get all mergerequest for user
		graphMergeRequests := graphUser.Get("assignedMergeRequests.nodes").Array()

		for _, graphMergeRequest := range graphMergeRequests {
			projectName := graphMergeRequest.Get("project.name").String()
			projectId := graphMergeRequest.Get("project.id").String()
			milestone := graphMergeRequest.Get("milestone.id").String()
			if milestone != "" {
				_, ok := projectUpcomingMilestones[projectId]
				if ok == false { // If MR is not in upcoming milestone we ignore it
					continue
				}
			}

			if graphMergeRequest.Get("draft").Bool() {
				count, ok := draftMergeRequests[projectName]
				if ok == false {
					draftMergeRequests[projectName] = 1
				} else {
					draftMergeRequests[projectName] = count + 1
				}
			} else {
				count, ok := mergeRequests[projectName]
				if ok == false {
					mergeRequests[projectName] = 1
				} else {
					mergeRequests[projectName] = count + 1
				}
			}
		}

		// group mergerequest by project
		// set metric per project
		for mr := range mergeRequests {
			user.MergeRequestsMetric.WithLabelValues(user.Name, mr).Set(float64(mergeRequests[mr]))
		}
		for mr := range draftMergeRequests {
			user.DraftMergeRequestsMetric.WithLabelValues(user.Name, mr).Set(float64(draftMergeRequests[mr]))
		}
	}

	// Format all projects
	fmt.Println("Updating projects")
	var projectQueries []string
	for _, project := range config.Projects {
		projectQueries = append(projectQueries, "\\\""+project.Id+"\\\"")
	}

	// Download issues for all projects
	payload = strings.NewReader("{\"query\":\"query {\\n    projects(ids:[" + strings.Join(projectQueries, ",") + "]) {\\n        nodes{\\n            id\\n            name\\n            issues(milestoneWildcardId: UPCOMING){\\n                nodes{\\n                    state\\n                    labels{\\n                        nodes{\\n                            title\\n                        }\\n                    }\\n                }\\n            }\\n        }\\n    }\\n}\",\"variables\":{}}")
	data, err = downloadData(payload)

	// Get all projects
	graphProjects := data.Get("projects.nodes").Array()
	for _, graphProject := range graphProjects {
		fmt.Println("- " + graphProject.Get("name").String())
		project, err := getProject(graphProject.Get("id").String())
		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		projectLabelCount := make(map[string]int)
		for _, label := range project.Labels {
			projectLabelCount[label] = 0
		}

		// Count issues in project
		issues := graphProject.Get("issues.nodes").Array()
		for _, issue := range issues {
			// If state is closed that takes priority over labels
			issueState := issue.Get("state").String()
			if issueState == "closed" {
				projectLabelCount[issueState] = projectLabelCount[issueState] + 1
				continue
			}

			// Try matching issues labels to project labels
			foundLabel := false
			for _, label := range issue.Get("labels.nodes").Array() {
				labelTitle := label.Get("title").String()
				_, ok := projectLabelCount[labelTitle]
				if ok { // if label exists
					projectLabelCount[labelTitle] = projectLabelCount[labelTitle] + 1
					foundLabel = true
					break
				}
			}
			if foundLabel {
				continue
			}

			// If no label match was found, add as an opened issue
			projectLabelCount[issueState] = projectLabelCount[issueState] + 1
		}

		// Update metrics
		index := 0
		for label := range projectLabelCount {
			fmt.Println(label + " " + strconv.Itoa(index))
			project.Metric.WithLabelValues(graphProject.Get("name").String(), label, strconv.Itoa(index)).Set(float64(projectLabelCount[label]))
			index++
		}
	}

	fmt.Println("Update Complete")
}

func getUser(username string) (User, error) {
	for _, user := range config.Users {
		if user.UserName == username {
			return user, nil
		}
	}
	return User{}, errors.New("Could not find user: " + username)
}

func getProject(id string) (Project, error) {
	for _, project := range config.Projects {
		if project.Id == id {
			return project, nil
		}
	}
	return Project{}, errors.New("Could not find project: " + id)
}

func setUpcomingMilestoneForProjects(data gjson.Result, projectUpcomingMilestones map[string]string) {
	// Update milestone map with milestone id
	projectIds := data.Get("projects.nodes").Array()
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
