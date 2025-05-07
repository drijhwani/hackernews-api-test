# Hacker News API Acceptance Tests

This Go test suite verifies key functionality of the [Hacker News API](https://github.com/HackerNews/API). It includes coverage for top stories, item details, comments, and edge cases like deleted or invalid items. The tests are enhanced with built-in **retry logic** to improve resilience against transient network or API errors.

## ✅ Test Coverage

| Test Name                                 | Description                                                                 |
|-------------------------------------------|-----------------------------------------------------------------------------|
| `TestAcceptance_RetrieveTopStories`       | Checks if top stories are returned by the API                              |
| `TestAcceptance_TopStoryDetails`          | Validates details (ID, title, type) of the first top story                 |
| `TestAcceptance_TopStoryFirstComment`     | Fetches and verifies the first comment of the top story                   |
| `TestAcceptance_NewStoriesList`           | Confirms the new stories endpoint returns data                            |
| `TestAcceptance_BestStoriesList`          | Verifies the best stories endpoint works                                  |
| `TestAcceptance_ItemTypeValidation`       | Asserts item types are one of `story`, `comment`, `poll`, etc.            |
| `TestAcceptance_DeletedItemReturnsNull`   | Validates deleted items return `null`                                     |
| `TestAPI_GetTopStories`                   | Tests basic retrieval of top stories list                                 |
| `TestAPI_GetTopStoryItem`                 | Fetches the full story data for the top story                             |
| `TestAPI_GetFirstCommentOfTopStory`       | Retrieves and validates the first comment for the top story               |
| `TestAPI_InvalidStoryID`                  | Checks that a bad item ID returns `null`                                  |
| `TestAPI_TopStoryNoComments`             | Skips if all stories have comments; logs the first story without comments |

## 🔁 Retry Logic

All test cases use a shared `getJSON` helper with automatic retry logic. This helps mitigate:

- Temporary network issues
- Transient API downtimes
- Intermittent JSON parse failures

**Retry settings:**
- Maximum attempts: 3
- Backoff delay: 2 seconds per retry

## Getting Started

### Prerequisites

- Go 1.20+
- Internet connection

### Run Tests

```bash
go test -v

## Run single test
go test -run '{test name}$'

## debug test
# install dlv 
go install -v github.com/go-delve/delve/cmd/dlv@latest 

## Run debugger
dlv test-- -test.run '{test name}$'

## Install dependencies
go get github.com/stretchr/testify/assert
