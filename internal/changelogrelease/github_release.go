package changelogrelease

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"hermannm.dev/errclose"
	"hermannm.dev/wrap"
	"hermannm.dev/wrap/ctxwrap"
)

type GitHubApiClient struct {
	httpClient *http.Client
	apiURL     string
}

func (client GitHubApiClient) createRelease(
	ctx context.Context,
	tagName string,
	releaseTitle string,
	changelog string,
	repoName string,
	repoOwner string,
	authToken string,
) (release CreatedRelease, returnedErr error) {
	url := fmt.Sprintf("%s/repos/%s/%s/releases", client.apiURL, repoOwner, repoName)

	requestBody := CreateReleaseRequest{
		TagName: tagName,
		Name:    releaseTitle,
		Body:    changelog,
	}
	requestBodyJson, err := json.Marshal(requestBody)
	if err != nil {
		return CreatedRelease{}, wrap.Error(
			err,
			"Failed to encode GitHub release request body JSON",
		)
	}

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		url,
		bytes.NewBuffer(requestBodyJson),
	)
	if err != nil {
		return CreatedRelease{}, wrap.Error(err, "Failed to create GitHub release HTTP request")
	}
	request.Header.Add("Authorization", "Bearer "+authToken)
	// GitHub API requires sending a User-Agent header, and they recommend setting it to your GitHub
	// username: https://docs.github.com/en/rest/using-the-rest-api/getting-started-with-the-rest-api#user-agent
	// In this case, that will be the repo owner.
	request.Header.Add("User-Agent", repoOwner)
	request.Header.Add("Accept", "application/vnd.github+json")
	request.Header.Add("X-GitHub-Api-Version", "2022-11-28")

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return CreatedRelease{}, wrap.Error(err, "Failed to send create release request to GitHub")
	}
	defer errclose.Close(response.Body, &returnedErr, "GitHub response body")

	if !isSuccessResponse(response) {
		return CreatedRelease{}, ctxwrap.NewErrorWithAttrs(
			ctx,
			"Got unsuccessful response from GitHub when trying to create release",
			"responseStatus", response.StatusCode,
			"responseBody", readErrorResponseBody(response),
		)
	}

	var responseBody CreateReleaseResponse
	if err := json.NewDecoder(response.Body).Decode(&responseBody); err != nil {
		return CreatedRelease{}, wrap.Error(
			err,
			"GitHub create release request succeeded, but failed to get release URL from response body",
		)
	}

	return CreatedRelease{
		Name: releaseTitle,
		URL:  responseBody.HtmlURL,
	}, nil
}

type CreateReleaseRequest struct {
	TagName string `json:"tag_name"`
	Name    string `json:"name"`
	Body    string `json:"body"`
}

type CreateReleaseResponse struct {
	HtmlURL string `json:"html_url"`
}

func isSuccessResponse(response *http.Response) bool {
	return response.StatusCode >= 200 && response.StatusCode <= 299
}

func readErrorResponseBody(response *http.Response) string {
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return "<failed to read>"
	}
	return string(responseBody)
}
