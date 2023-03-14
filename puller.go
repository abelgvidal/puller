package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/slack-go/slack"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Repos        []string `yaml:"repos"`
	Slacktoken   string
	Githubtoken  string
	Slackchannel string
}

type PullRequest struct {
	Title string `json:"title"`
	URL   string `json:"html_url"`
}

func sendToSlack(texto string, config Config) error {
	api := slack.New(config.Slacktoken)
	attachment := slack.Attachment{
		Text: texto,
	}

	_, timestamp, err := api.PostMessage(
		config.Slackchannel,
		slack.MsgOptionAttachments(attachment),
		slack.MsgOptionAsUser(true),
	)
	if err != nil {
		fmt.Printf("%s\n", err)
		return err
	}
	fmt.Printf("Message successfully sent to channel %s at %s", config.Slackchannel, timestamp)

	return nil
}

func getPullRequestsFromOneRepo(repo string, Githubtoken string) ([]PullRequest, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("https://api.github.com/repos/%s/pulls", repo), nil)
	if err != nil {
		panic(err)
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", Githubtoken))

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
		pullRequests, _ := getPullRequestsFromOneRepo(repo, config.Githubtoken)
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

	if len(pullRequests) > 0 {
		sendToSlack("¡Qué momento más maravilloso para revisar los pull requests pendientes!", config)
	}
	for _, pr := range pullRequests {
		err := sendToSlack(fmt.Sprintf("Title: %s\nURL: %s\n\n", pr.Title, pr.URL), config)
		if err != nil {
			fmt.Print("Error sending to slack:", err)
		}
		fmt.Printf("Title: %s\nURL: %s\n\n", pr.Title, pr.URL)
	}
}
