package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Repos []string `yaml:"repos"`
	Token string
}

type PullRequest struct {
	Title string `json:"title"`
	URL   string `json:"html_url"`
}

func getPullRequestsFromOneRepo(repo string, token string) ([]PullRequest, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("https://api.github.com/repos/%s/pulls", repo), nil)
	if err != nil {
		panic(err)
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var pullRequests []PullRequest
	err = json.NewDecoder(resp.Body).Decode(&pullRequests)
	if err != nil {
		panic(err)
	}

	return pullRequests, nil
}

func getPullRequests(config Config) ([]PullRequest, error) {
	var totalPullRequests []PullRequest

	for _, repo := range config.Repos {
		pullRequests, _ := getPullRequestsFromOneRepo(repo, config.Token)
		totalPullRequests = append(totalPullRequests, pullRequests...)
	}
	return totalPullRequests, nil
}

func getConfig() Config {
	f, err := os.Open("config.yml")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	var cfg Config
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		panic(err)
	}
	return cfg
}

func main() {
	config := getConfig()
	pullRequests, _ := getPullRequests(config)

	for _, pr := range pullRequests {
		fmt.Printf("Title: %s\nURL: %s\n\n", pr.Title, pr.URL)
	}
}
