package output

import (
	"bytes"
	"os"
)

type Output struct {
	Clients map[Type]outputClient
}

type outputClient interface {
	Publish(buf *bytes.Buffer) error
}

type Type string

const (
	TypeGit   Type = "git"
	TypeEsa   Type = "esa"
	TypeLocal Type = "local"
)

// Config ...
type Config struct {
	Git   *GitConfig   `yaml:"git"`
	Local *LocalConfig `yaml:"local"`
}

// NewClient ...
func NewClient(config *Config) (*Output, error) {
	clients := map[Type]outputClient{}
	if config.Git != nil {
		config.Git.Password = getGitPassword()
		config.Git.AuthorUserID = getGitAuthorUser()
		config.Git.AuthorEmail = getGitAuthorEmail()
		gc, err := NewGitClient(config.Git)
		if err != nil {
			return nil, err
		}
		clients[TypeGit] = gc
	}
	if config.Local != nil {
		lc, err := NewLocalClient(config.Local)
		if err != nil {
			return nil, err
		}
		clients[TypeLocal] = lc
	}
	return &Output{
		Clients: clients,
	}, nil
}

// Publishs is prints documents to multiple clients
func (out *Output) Publishs(buf *bytes.Buffer) error {
	for _, cli := range out.Clients {
		if err := cli.Publish(buf); err != nil {
			return err
		}
	}
	return nil
}

func getGitPassword() string {
	return os.Getenv("OUT_GIT_PASSWORD")
}

func getGitAuthorUser() string {
	return os.Getenv("OUT_GIT_AUTHOR_USER")
}

func getGitAuthorEmail() string {
	return os.Getenv("OUT_GIT_AUTHOR_EMAIL")
}
