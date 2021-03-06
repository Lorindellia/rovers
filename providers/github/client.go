package github

import (
	"encoding/json"
	"fmt"
	"golang.org/x/oauth2"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/src-d/rovers/providers/github/model"
)

const (
	httpTimeout  = 30 * time.Second
	githubApiURL = "https://api.github.com/repositories?since=%d"

	rateLimitLimitHeader     = "X-RateLimit-Limit"
	rateLimitRemainingHeader = "X-RateLimit-Remaining"
)

type response struct {
	Next         int
	Repositories []*model.Repository

	Total     int
	Remaining int
}

type client struct {
	c *http.Client
}

func newClient(token string) *client {
	c := &http.Client{}

	if token != "" {
		t := &oauth2.Token{AccessToken: token}
		c = oauth2.NewClient(oauth2.NoContext, oauth2.StaticTokenSource(t))
	}

	c.Timeout = httpTimeout

	return &client{c}
}

// Repositories returns a response with the next page id and a list of Repositories.
// It automatically slow down if we are doing requests too fast.
func (c *client) Repositories(since int) (*response, error) {
	start := time.Now()

	u := fmt.Sprintf(githubApiURL, since)
	res, err := c.c.Get(u)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	if res.StatusCode >= 400 {
		return nil, fmt.Errorf("request error. Status code %s, %s", res.Status, res.Status)
	}

	repositories, err := c.decode(res.Body)

	total := c.toInt(res.Header.Get(rateLimitLimitHeader))
	remaining := c.toInt(res.Header.Get(rateLimitRemainingHeader))

	minRequestDuration := time.Hour / time.Duration(total)
	defer func() {
		needsWait := minRequestDuration - time.Since(start)
		if needsWait > 0 {
			time.Sleep(needsWait)
		}
	}()

	next := 0
	if len(repositories) != 0 {
		next = repositories[len(repositories)-1].GithubID
	}

	return &response{
		Next:         next,
		Repositories: repositories,
		Total:        total,
		Remaining:    remaining,
	}, nil
}

func (c *client) toInt(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}

func (c *client) decode(body io.Reader) ([]*model.Repository, error) {
	var record []*model.Repository
	if err := json.NewDecoder(body).Decode(&record); err != nil {
		return nil, err
	}

	return record, nil
}
