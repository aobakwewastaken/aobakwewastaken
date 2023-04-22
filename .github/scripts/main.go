package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

type Story struct {
	ID int `json:"id"`
	Title string `json:"title"`
	URL string `json:"url"`
}
func main() {
	var readMeContent string
	accessToken := os.Getenv("MY_ACCESS_TOKEN")
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{ AccessToken: accessToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	owner := "aobakwewastaken"
	repo := "aobakwewastaken"
	readme, _, err := client.Repositories.GetReadme(ctx, owner, repo, &github.RepositoryContentGetOptions{})
	if err != nil {
		fmt.Println("Failed to fetch readme: ", err)
	}

	resp, err := http.Get("https://hacker-news.firebaseio.com/v0/topstories.json")
	if err != nil {
		fmt.Println("Failed to fetch top stories: ", err)
	}
	defer resp.Body.Close()

	var storyIDs []int

	err = json.NewDecoder(resp.Body).Decode(&storyIDs)

	if err != nil {
		fmt.Println("Failed to decode response: ", err)
		return
	}
	readMeContent = "# Top Stories on hackernews <br />"
	for i := 0; i < 10 && i < len(storyIDs); i++ {
		storyID := storyIDs[i]
		storyResp, err := http.Get(fmt.Sprintf("https://hacker-news.firebaseio.com/v0/item/%d.json", storyID))
		if err != nil {
			fmt.Println("Failed to fetch story: ", err)
			continue
		}
		defer storyResp.Body.Close()

		var story Story
		err = json.NewDecoder(storyResp.Body).Decode(&story)
		if err != nil {
			fmt.Println("Failed to decode story:", err)
			continue
		}
		readMeContent += fmt.Sprintf("\n[%s](%s)\n", story.Title, story.URL)
	}
	currentReadMeContent, err := readme.GetContent()

	if err != nil {
		fmt.Println("Failed to get content: ", err)
	}
	currentReadMeContent = readMeContent
	opts := &github.RepositoryContentFileOptions{
		Message: github.String("Update README"),
		Content: []byte(currentReadMeContent),
		SHA:     readme.SHA,
	}
	_, _, err = client.Repositories.UpdateFile(ctx, owner, repo, *readme.Path, opts)
	if err != nil {
		fmt.Println("Failed to update Readme: ", err)
	}
	
}