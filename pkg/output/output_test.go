package output

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var defaultConf = &Config{
	Git: &GitConfig{
		URL:      "http://example.com",
		Branch:   "master",
		FilePath: "README.md",
	},
	Local: &LocalConfig{
		FilePath: "./outupt.md",
	},
}

func TestNewClient(t *testing.T) {
	t.Parallel()

	cli, err := NewClient(defaultConf)
	assert.Error(t, err)

	os.Setenv("OUT_GIT_PASSWORD", "aaa")
	os.Setenv("OUT_GIT_AUTHOR_USER", "userA")
	os.Setenv("OUT_GIT_AUTHOR_EMAIL", "user_a@example.com")
	cli, err = NewClient(defaultConf)
	if assert.NoError(t, err) {
		assert.Equal(t, len(cli.Clients), 2)

		_, ok := cli.Clients[TypeLocal]
		assert.True(t, ok)
	}
}
