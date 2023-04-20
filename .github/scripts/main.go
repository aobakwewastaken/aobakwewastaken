package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Story struct {
	ID int `json:"id"`
	Title string `json:"title"`
	URL string `json:"url"`
}
func main() {
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
		fmt.Printf("Title: %s\nURL: %s\n\n", story.Title, story.URL)
	}
}