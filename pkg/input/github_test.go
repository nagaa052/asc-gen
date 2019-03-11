package input

import (
	"testing"
	"time"

	ggh "github.com/google/go-github/github"
	"github.com/stretchr/testify/assert"
)

var (
	token       = "token"
	org         = "org"
	defaultConf = &GithubConfig{
		Organization: org,
		Token:        token,
	}
)

func TestNewGithubClient(t *testing.T) {
	t.Parallel()

	_, err := NewGithubClient(&GithubConfig{
		Organization: org,
	})
	assert.Error(t, err)

	cli, err := NewGithubClient(defaultConf)
	if assert.NoError(t, err) {
		assert.Equal(t, cli.conf.Token, token)
		assert.Equal(t, cli.conf.Organization, org)
		assert.NotNil(t, cli.client)
	}
}

func TestGitHubGetItems(t *testing.T) {
	t.Parallel()

	cli, err := NewGithubClient(defaultConf)

	if assert.NoError(t, err) {

		now := ggh.Timestamp{Time: time.Now()}
		mockRepos := []struct {
			name      string
			desc      string
			url       string
			updatedAt ggh.Timestamp
			topcis    []string
		}{
			{"repo1", "", "https://example.com/apc/example1.git", now, []string{"aaa", "bbb", "ccc"}},
			{"repo2", "repo2 desc", "https://example.com/apc/example2.git", now, []string{"ddd"}},
			{"repo3", "", "https://example.com/apc/example3.git", now, []string{}},
		}

		cli.getRepos = func(cli *Github) ([]*ggh.Repository, error) {
			repos := make([]*ggh.Repository, 0)
			for _, r := range mockRepos {
				repos = append(repos, &ggh.Repository{
					Name:        sp(r.name),
					Description: sp(r.desc),
					HTMLURL:     sp(r.url),
					UpdatedAt:   &r.updatedAt,
					Topics:      r.topcis,
				})
			}
			return repos, nil
		}

		items, err := cli.GetItems()
		if assert.NoError(t, err) {
			assert.Equal(t, len(items), 3)
			assert.Equal(t, items[0].Name, "repo1")
			assert.Equal(t, items[0].Description, "")
			assert.Equal(t, items[1].Description, "repo2 desc")
			assert.Equal(t, len(items[1].Labels), 1)
			assert.Equal(t, items[2].URL, "https://example.com/apc/example3.git")
		}

		// caching item test
		items2, err := cli.GetItems()
		if assert.NoError(t, err) {
			assert.Equal(t, items2[0], items[0])
		}
	}
}

func sp(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
