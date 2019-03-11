package input

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

type cliMock1 struct{}

func (cm1 *cliMock1) GetItems() ([]*Item, error) {
	items := []*Item{
		&Item{Name: "item1-1"},
		&Item{Name: "item1-2"},
		&Item{Name: "item1-3"},
	}
	return items, nil
}

type cliMock2 struct{}

func (cm2 *cliMock2) GetItems() ([]*Item, error) {
	items := []*Item{
		&Item{Name: "item2-1"},
		&Item{Name: "item2-2"},
	}
	return items, nil
}

func TestNewClient(t *testing.T) {
	t.Parallel()

	os.Setenv("IN_GITHUB_TOKEN", "aaa")
	conf := &Config{
		Github: &GithubConfig{
			Organization: "org",
		},
	}

	cli, err := NewClient(conf)
	if assert.NoError(t, err) {
		assert.Equal(t, len(cli.Clients), 1)

		_, ok := cli.Clients[TypeGithub]
		assert.True(t, ok)
	}
}

func TestGetItems(t *testing.T) {
	t.Parallel()

	clis := map[Type]inputClient{}
	clis[TypeGithub] = &cliMock1{}
	clis[TypeEsa] = &cliMock2{}
	in := &Input{
		Clients: clis,
	}

	items, err := in.GetItems()
	if assert.NoError(t, err) {
		assert.Equal(t, len(items), 5)
	}
}
