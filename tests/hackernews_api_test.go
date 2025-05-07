package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const baseURL = "https://hacker-news.firebaseio.com/v0"

func stringifyID(id int) string {
	return fmt.Sprintf("%d", id)
}

func getJSON(t *testing.T, url string, target interface{}) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("GET failed for %s: %v", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Unexpected HTTP status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("Error reading body: %v", err)
	}

	if err := json.Unmarshal(body, target); err != nil {
		return fmt.Errorf("Unmarshal failed: %v", err)
	}
	return nil
}

func retryRequest(fn func() error) error {
	const maxRetries = 3
	for i := 0; i < maxRetries; i++ {
		if err := fn(); err != nil {
			log.Printf("Attempt %d failed: %v", i+1, err)
			time.Sleep(time.Second * time.Duration(i+1))
			continue
		}
		return nil
	}
	return fmt.Errorf("all retry attempts failed")
}

func TestAcceptance_RetrieveTopStories(t *testing.T) {
	log.Println("TestAcceptance_RetrieveTopStories started")
	var storyIDs []int
	err := retryRequest(func() error {
		return getJSON(t, baseURL+"/topstories.json", &storyIDs)
	})
	assert.NoError(t, err)
	assert.Greater(t, len(storyIDs), 0)
}

func TestAcceptance_TopStoryDetails(t *testing.T) {
	log.Println("TestAcceptance_TopStoryDetails started")
	var storyIDs []int
	err := retryRequest(func() error {
		return getJSON(t, baseURL+"/topstories.json", &storyIDs)
	})
	assert.NoError(t, err)
	assert.Greater(t, len(storyIDs), 0)

	var story struct {
		ID    int    `json:"id"`
		Title string `json:"title"`
		Type  string `json:"type"`
	}
	err = retryRequest(func() error {
		return getJSON(t, baseURL+"/item/"+stringifyID(storyIDs[0])+".json", &story)
	})
	assert.NoError(t, err)
	assert.Equal(t, storyIDs[0], story.ID)
	assert.NotEmpty(t, story.Title)
	assert.Equal(t, "story", story.Type)
}

func TestAcceptance_TopStoryFirstComment(t *testing.T) {
	log.Println("TestAcceptance_TopStoryFirstComment started")
	var storyIDs []int
	err := retryRequest(func() error {
		return getJSON(t, baseURL+"/topstories.json", &storyIDs)
	})
	assert.NoError(t, err)
	assert.Greater(t, len(storyIDs), 0)

	var story struct {
		ID   int   `json:"id"`
		Kids []int `json:"kids"`
	}
	err = retryRequest(func() error {
		return getJSON(t, baseURL+"/item/"+stringifyID(storyIDs[0])+".json", &story)
	})
	assert.NoError(t, err)

	if len(story.Kids) == 0 {
		t.Skip("Top story has no comments")
	}

	var comment struct {
		ID   int    `json:"id"`
		Text string `json:"text"`
		Type string `json:"type"`
	}
	err = retryRequest(func() error {
		return getJSON(t, baseURL+"/item/"+stringifyID(story.Kids[0])+".json", &comment)
	})
	assert.NoError(t, err)
	assert.Equal(t, "comment", comment.Type)
	assert.NotEmpty(t, comment.Text)
}

func TestAcceptance_NewStoriesList(t *testing.T) {
	log.Println("TestAcceptance_NewStoriesList started")
	var storyIDs []int
	err := retryRequest(func() error {
		return getJSON(t, baseURL+"/newstories.json", &storyIDs)
	})
	assert.NoError(t, err)
	assert.Greater(t, len(storyIDs), 0)
}

func TestAcceptance_BestStoriesList(t *testing.T) {
	log.Println("TestAcceptance_BestStoriesList started")
	var storyIDs []int
	err := retryRequest(func() error {
		return getJSON(t, baseURL+"/beststories.json", &storyIDs)
	})
	assert.NoError(t, err)
	assert.Greater(t, len(storyIDs), 0)
}

func TestAcceptance_ItemTypeValidation(t *testing.T) {
	log.Println("TestAcceptance_ItemTypeValidation started")
	var storyIDs []int
	err := retryRequest(func() error {
		return getJSON(t, baseURL+"/topstories.json", &storyIDs)
	})
	assert.NoError(t, err)
	assert.Greater(t, len(storyIDs), 0)

	var item struct {
		ID   int    `json:"id"`
		Type string `json:"type"`
	}
	err = retryRequest(func() error {
		return getJSON(t, baseURL+"/item/"+stringifyID(storyIDs[0])+".json", &item)
	})
	assert.NoError(t, err)
	validTypes := map[string]bool{"story": true, "comment": true, "poll": true, "job": true, "pollopt": true}
	assert.True(t, validTypes[item.Type])
}

func TestAcceptance_DeletedItemReturnsNull(t *testing.T) {
	log.Println("TestAcceptance_DeletedItemReturnsNull started")
	err := retryRequest(func() error {
		resp, err := http.Get(baseURL + "/item/0.json")
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		assert.Equal(t, "null", string(body))
		return nil
	})
	assert.NoError(t, err)
}

func TestAPI_GetTopStories(t *testing.T) {
	log.Println("TestAPI_GetTopStories started")
	var storyIDs []int
	err := retryRequest(func() error {
		return getJSON(t, baseURL+"/topstories.json", &storyIDs)
	})
	assert.NoError(t, err)
	assert.Greater(t, len(storyIDs), 0)
}

func TestAPI_GetTopStoryItem(t *testing.T) {
	log.Println("TestAPI_GetTopStoryItem started")
	var storyIDs []int
	err := retryRequest(func() error {
		return getJSON(t, baseURL+"/topstories.json", &storyIDs)
	})
	assert.NoError(t, err)
	assert.Greater(t, len(storyIDs), 0)

	var story struct {
		ID    int    `json:"id"`
		Title string `json:"title"`
		Type  string `json:"type"`
	}
	err = retryRequest(func() error {
		return getJSON(t, baseURL+"/item/"+stringifyID(storyIDs[0])+".json", &story)
	})
	assert.NoError(t, err)
	assert.Equal(t, storyIDs[0], story.ID)
	assert.NotEmpty(t, story.Title)
	assert.Equal(t, "story", story.Type)
}

func TestAPI_GetFirstCommentOfTopStory(t *testing.T) {
	log.Println("TestAPI_GetFirstCommentOfTopStory started")
	var storyIDs []int
	err := retryRequest(func() error {
		return getJSON(t, baseURL+"/topstories.json", &storyIDs)
	})
	assert.NoError(t, err)
	assert.Greater(t, len(storyIDs), 0)

	var story struct {
		ID   int   `json:"id"`
		Kids []int `json:"kids"`
	}
	err = retryRequest(func() error {
		return getJSON(t, baseURL+"/item/"+stringifyID(storyIDs[0])+".json", &story)
	})
	assert.NoError(t, err)
	if len(story.Kids) == 0 {
		t.Skip("No comments on top story")
	}

	var comment struct {
		ID   int    `json:"id"`
		Text string `json:"text"`
		Type string `json:"type"`
	}
	err = retryRequest(func() error {
		return getJSON(t, baseURL+"/item/"+stringifyID(story.Kids[0])+".json", &comment)
	})
	assert.NoError(t, err)
	assert.Equal(t, "comment", comment.Type)
	assert.NotEmpty(t, comment.Text)
}

func TestAPI_InvalidStoryID(t *testing.T) {
	log.Println("TestAPI_InvalidStoryID started")
	err := retryRequest(func() error {
		resp, err := http.Get(baseURL + "/item/0.json")
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		assert.Equal(t, "null", string(body))
		return nil
	})
	assert.NoError(t, err)
}

func TestAPI_TopStoryNoComments(t *testing.T) {
	log.Println("TestAPI_TopStoryNoComments started")
	var storyIDs []int
	err := retryRequest(func() error {
		return getJSON(t, baseURL+"/topstories.json", &storyIDs)
	})
	assert.NoError(t, err)
	assert.Greater(t, len(storyIDs), 0)

	found := false
	for _, id := range storyIDs {
		var story struct {
			ID   int   `json:"id"`
			Kids []int `json:"kids"`
		}
		err := retryRequest(func() error {
			return getJSON(t, baseURL+"/item/"+stringifyID(id)+".json", &story)
		})
		assert.NoError(t, err)
		if len(story.Kids) == 0 {
			log.Printf("Found story with no comments: %d", id)
			found = true
			break
		}
	}
	if !found {
		t.Skip("All top stories have comments")
	}
}
