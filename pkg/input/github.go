package input

import (
	"context"
	"log"
	"time"

	"github.com/pkg/errors"

	ggh "github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

// GithubConfig ...
type GithubConfig struct {
	Token        string
	Organization string `yaml:"organization"`
}

// Github ...
type Github struct {
	conf     *GithubConfig
	client   *ggh.Client
	repos    []*ggh.Repository
	getRepos func(gh *Github) ([]*ggh.Repository, error)
}

// NewGithubClient ...
func NewGithubClient(config *GithubConfig) (*Github, error) {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: config.Token},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	gh := &Github{
		conf:     config,
		client:   ggh.NewClient(tc),
		repos:    make([]*ggh.Repository, 0),
		getRepos: getAllRepository,
	}
	if !validGithub(gh) {
		return nil, errors.New("validation failed github client")
	}
	return gh, nil
}

// GetItems ...
func (gh *Github) GetItems() ([]*Item, error) {
	items := make([]*Item, 0)
	repos, err := gh.getRepos(gh)
	if err != nil {
		return nil, err
	}
	for _, repo := range repos {
		items = append(items, &Item{
			Name:        ps(repo.Name),
			Description: ps(repo.Description),
			URL:         ps(repo.HTMLURL),
			UpdatedAt:   ts(repo.UpdatedAt),
			Labels:      repo.Topics,
		})
	}
	return items, nil
}

func validGithub(gh *Github) bool {
	if gh.conf.Token == "" {
		return false
	}
	if gh.conf.Organization == "" {
		return false
	}
	return true
}

func getAllRepository(gh *Github) ([]*ggh.Repository, error) {
	if len(gh.repos) > 0 {
		return gh.repos, nil
	}

	log.Printf("call api to get repository of %s ...", gh.conf.Organization)
	opt := &ggh.RepositoryListByOrgOptions{
		Type:        "all",
		ListOptions: ggh.ListOptions{PerPage: 100},
	}
	for {
		repos, res, err := gh.client.Repositories.ListByOrg(context.Background(), gh.conf.Organization, opt)
		if err != nil {
			return nil, err
		}
		for _, repo := range repos {
			gh.repos = append(gh.repos, repo)
		}

		if res.NextPage == 0 {
			break
		}
		opt.Page = res.NextPage
		time.Sleep(5 * time.Second)
	}
	log.Printf("get repository data is complete. The number of repositories is %d", len(gh.repos))
	return gh.repos, nil
}

func ps(p *string) string {
	if p != nil {
		return *p
	}
	return ""
}

func ts(t *ggh.Timestamp) string {
	if t != nil {
		return t.Format("2006-01-02")
	}
	return ""
}
