package input

import "os"

type Input struct {
	Clients map[Type]inputClient
}

type inputClient interface {
	GetItems() ([]*Item, error)
}

// Item ...
type Item struct {
	Name        string
	Description string
	URL         string
	UpdatedAt   string
	Labels      []string
}

type Config struct {
	Github *GithubConfig `yaml:"github"`
}

type Type string

const (
	TypeGithub    Type = "github"
	TypeGitlab    Type = "gitlab"
	TypeBitbucket Type = "bitbucket"
	TypeEsa       Type = "esa"
)

// NewClient ...
func NewClient(config *Config) (*Input, error) {
	clients := map[Type]inputClient{}
	if config.Github != nil {
		config.Github.Token = getGithubToken()
		gc, err := NewGithubClient(config.Github)
		if err != nil {
			return nil, err
		}
		clients[TypeGithub] = gc
	}
	return &Input{
		Clients: clients,
	}, nil
}

// GetItems ...
func (m *Input) GetItems() ([]*Item, error) {
	items := make([]*Item, 0)
	for _, cl := range m.Clients {
		item, err := cl.GetItems()
		if err != nil {
			return nil, err
		}
		items = append(items, item...)
	}
	return items, nil
}

func getGithubToken() string {
	return os.Getenv("IN_GITHUB_TOKEN")
}
