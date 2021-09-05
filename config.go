package main

type Config struct {
	Users []User `json:"users"`
}

type User struct {
	Name             string `json:"name"`
	Iid              int    `json:"iid"`
	MergeRequests    int    `json:"mergeRequests"`
	WipMergeRequests int    `json:"wipMergeRequests"`
}
