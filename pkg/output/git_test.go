package output

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var gitConf = &GitConfig{
	Password:     "123123",
	URL:          "http://example.com",
	AuthorUserID: "userA",
	AuthorEmail:  "user_a@example.com",
}

func TestNewGitClient(t *testing.T) {
	t.Parallel()

	g, err := NewGitClient(&GitConfig{
		URL:      "http://example.com",
		Branch:   "master",
		FilePath: "README.md",
	})
	assert.Error(t, err)

	g, err = NewGitClient(gitConf)
	if assert.NoError(t, err) {
		assert.Equal(t, g.config.Password, gitConf.Password)
		assert.Equal(t, g.config.AuthorUserID, gitConf.AuthorUserID)
		assert.Equal(t, g.config.AuthorEmail, gitConf.AuthorEmail)
		assert.Equal(t, g.config.Branch, "master")
		assert.Equal(t, g.config.FilePath, "README.md")
	}
}
