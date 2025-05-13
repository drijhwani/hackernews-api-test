# Hacker News API Acceptance Tests

This Go test suite verifies key functionality of the [Hacker News API](https://github.com/HackerNews/API). It includes coverage for top stories, item details, comments, and edge cases like deleted or invalid items. The tests are enhanced with built-in **retry logic** to improve resilience against transient network or API errors.

## ✅ Test Coverage

| **Test Function**                         | **Description**                                                               |
| ----------------------------------------- | ----------------------------------------------------------------------------- |
| `TestAcceptance_RetrieveTopStories`       | Verifies the `/topstories.json` endpoint returns a non-empty list of IDs.     |
| `TestAcceptance_TopStoryDetails`          | Fetches the top story and checks for a valid `id`, `title`, and type `story`. |
| `TestAcceptance_TopStoryFirstComment`     | Retrieves the first comment of the top story and validates its structure.     |
| `TestAcceptance_NewStoriesList`           | Validates the `/newstories.json` endpoint returns story IDs.                  |
| `TestAcceptance_BestStoriesList`          | Validates the `/beststories.json` endpoint returns story IDs.                 |
| `TestAcceptance_ItemTypeValidation`       | Confirms that the returned item type is one of the valid types.               |
| `TestAcceptance_DeletedItemReturnsNull`   | Checks that an invalid item ID returns `"null"`.                              |
| `TestAcceptance_UpdatesEndpoint`          | Ensures `/updates.json` returns either updated items or profiles.             |
| `TestAcceptance_UpdatesItemValidation`    | Fetches details of an updated item and checks ID and type fields.             |
| `TestAcceptance_UpdatesProfileValidation` | Validates the structure of an updated user profile.                           |
| `TestAPI_GetTopStories`                   | Duplicate of top stories check; ensures IDs are fetched correctly.            |
| `TestAPI_GetTopStoryItem`                 | Ensures the first top story has a valid `id`, `title`, and type `story`.      |
| `TestAPI_GetFirstCommentOfTopStory`       | Retrieves and validates the first comment of the top story.                   |
| `TestAPI_InvalidStoryID`                  | Ensures item ID `0` returns `"null"` as it is invalid/deleted.                |
| `TestAPI_TopStoryNoComments`              | Scans top stories and finds one without comments to ensure skip logic works.  |

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
go install github.com/go-delve/delve/cmd/dlv@latest
