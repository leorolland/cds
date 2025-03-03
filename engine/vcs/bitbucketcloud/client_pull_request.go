package bitbucketcloud

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/url"

	"github.com/rockbears/log"

	"github.com/ovh/cds/sdk"
)

func (client *bitbucketcloudClient) PullRequest(ctx context.Context, fullname string, id string) (sdk.VCSPullRequest, error) {
	url := fmt.Sprintf("/repositories/%s/pullrequests/%s", fullname, id)
	status, body, _, err := client.get(ctx, url)
	if err != nil {
		log.Warn(ctx, "bitbucketcloudClient.Pullrequest> Error %s", err)
		return sdk.VCSPullRequest{}, err
	}
	if status >= 400 {
		return sdk.VCSPullRequest{}, sdk.NewError(sdk.ErrRepoNotFound, errorAPI(body))
	}
	var pullrequest PullRequest
	if err := sdk.JSONUnmarshal(body, &pullrequest); err != nil {
		log.Warn(ctx, "bitbucketcloudClient.PullRequest> Unable to parse bitbucket cloud commit: %s", err)
		return sdk.VCSPullRequest{}, err
	}

	return pullrequest.ToVCSPullRequest(), nil
}

// PullRequests fetch all the pull request for a repository
func (client *bitbucketcloudClient) PullRequests(ctx context.Context, fullname string, opts sdk.VCSPullRequestOptions) ([]sdk.VCSPullRequest, error) {
	var pullrequests []PullRequest
	path := fmt.Sprintf("/repositories/%s/pullrequests", fullname)
	params := url.Values{}

	switch opts.State {
	case sdk.VCSPullRequestStateOpen:
		params.Set("state", "OPEN")
	case sdk.VCSPullRequestStateMerged:
		params.Set("state", "MERGED")
	case sdk.VCSPullRequestStateClosed:
		params.Set("state", "DECLINED")
	}

	params.Set("pagelen", "50")
	nextPage := 1
	for {
		if ctx.Err() != nil {
			break
		}

		if nextPage != 1 {
			params.Set("page", fmt.Sprintf("%d", nextPage))
		}

		var response PullRequests
		if err := client.do(ctx, "GET", "core", path, params, nil, &response); err != nil {
			return nil, sdk.WrapError(err, "Unable to get pull requests")
		}
		if cap(pullrequests) == 0 {
			pullrequests = make([]PullRequest, 0, response.Size)
		}

		pullrequests = append(pullrequests, response.Values...)

		if response.Next == "" {
			break
		} else {
			nextPage++
		}
	}

	responsePullRequest := make([]sdk.VCSPullRequest, 0, len(pullrequests))
	for _, pr := range pullrequests {
		responsePullRequest = append(responsePullRequest, pr.ToVCSPullRequest())
	}

	return responsePullRequest, nil
}

type BitbucketCloudPullRequestComment struct {
	Content struct {
		Raw    string `json:"raw,omitempty"`
		Markup string `json:"markup,omitempty"`
		HTML   string `json:"html,omitempty"`
	} `json:"content,omitempty"`
}

// PullRequestComment push a new comment on a pull request
func (client *bitbucketcloudClient) PullRequestComment(ctx context.Context, repo string, prRequest sdk.VCSPullRequestCommentRequest) error {
	if client.DisableStatus {
		log.Warn(ctx, "bitbucketcloud.PullRequestComment>  ⚠ bitbucketcloud statuses are disabled")
		return nil
	}

	path := fmt.Sprintf("/repositories/%s/pullrequests/%d/comments", repo, prRequest.ID)
	payload := BitbucketCloudPullRequestComment{}
	payload.Content.Raw = prRequest.Message

	values, _ := json.Marshal(payload)
	res, err := client.post(ctx, path, "application/json", bytes.NewReader(values), &postOptions{skipDefaultBaseURL: false, asUser: true})
	if err != nil {
		return sdk.WrapError(err, "Unable to post comment")
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return sdk.WrapError(err, "Unable to read body")
	}

	log.Debug(ctx, "%v", string(body))

	if res.StatusCode != 201 {
		return sdk.WrapError(err, "Unable to add comment. Status code : %d - Body: %s", res.StatusCode, body)
	}

	return nil
}

func (client *bitbucketcloudClient) PullRequestCreate(ctx context.Context, repo string, pr sdk.VCSPullRequest) (sdk.VCSPullRequest, error) {
	path := fmt.Sprintf("/repos/%s/pulls", repo)
	payload := map[string]string{
		"title": pr.Title,
		"head":  pr.Head.Branch.DisplayID,
		"base":  pr.Base.Branch.DisplayID,
	}
	values, _ := json.Marshal(payload)
	res, err := client.post(ctx, path, "application/json", bytes.NewReader(values), &postOptions{skipDefaultBaseURL: false, asUser: true})
	if err != nil {
		return sdk.VCSPullRequest{}, sdk.WrapError(err, "Unable to post status")
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return sdk.VCSPullRequest{}, sdk.WrapError(err, "Unable to read body")
	}

	var prResponse PullRequest
	if err := sdk.JSONUnmarshal(body, &prResponse); err != nil {
		return sdk.VCSPullRequest{}, sdk.WrapError(err, "Unable to unmarshal pullrequest %s", string(body))
	}

	return prResponse.ToVCSPullRequest(), nil
}

func (pullr PullRequest) ToVCSPullRequest() sdk.VCSPullRequest {
	return sdk.VCSPullRequest{
		ID: pullr.ID,
		Base: sdk.VCSPushEvent{
			Repo: pullr.Destination.Repository.FullName,
			Branch: sdk.VCSBranch{
				ID:           pullr.Destination.Branch.Name,
				DisplayID:    pullr.Destination.Branch.Name,
				LatestCommit: pullr.Destination.Commit.Hash,
			},
			Commit: sdk.VCSCommit{
				Hash: pullr.Destination.Commit.Hash,
			},
		},
		Head: sdk.VCSPushEvent{
			Repo: pullr.Source.Repository.FullName,
			Branch: sdk.VCSBranch{
				ID:           pullr.Source.Branch.Name,
				DisplayID:    pullr.Source.Branch.Name,
				LatestCommit: pullr.Source.Commit.Hash,
			},
			Commit: sdk.VCSCommit{
				Hash: pullr.Source.Commit.Hash,
			},
		},
		URL: pullr.Links.HTML.Href,
		User: sdk.VCSAuthor{
			Avatar:      pullr.Author.Links.Avatar.Href,
			DisplayName: pullr.Author.DisplayName,
			Name:        pullr.Author.Username,
			ID:          pullr.Author.UUID,
		},
		Closed:  pullr.State == "SUPERSEDED",
		Merged:  pullr.State == "MERGED",
		Updated: pullr.UpdatedOn,
	}
}
