package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	yconf = `
input: 
  github:
    organization: org
output:
  labels: 
  - label1
  - label2
  title: My Repository List
  out:
    git:
      url: https://example.com/org/my-repos.git
      branch: master
      filepath: README.md
    local:
      filepath: README.md

`
)

func setup()    {}
func teardown() {}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	if code == 0 {
		teardown()
	}
	os.Exit(code)
}

func TestGetConfig(t *testing.T) {
	t.Parallel()

	defaultConfigPath = "../../config/config.yaml"

	// test import config file
	config, err := getConfig()
	if assert.NoError(t, err) {
		assert.Equal(t, config.Input.Github.Organization, "OrgA")
		assert.Equal(t, config.Output.Title, "My Repository List")
		assert.Equal(t, config.Output.Out.Local.FilePath, "./output.md")
	}

	// test import env config
	os.Setenv("ASC_CONFIG", yconf)
	config, err = getConfig()
	if assert.NoError(t, err) {
		assert.Equal(t, config.Input.Github.Organization, "org")
		assert.Equal(t, len(config.Output.Labels), 2)
		assert.Equal(t, config.Output.Out.Git.URL, "https://example.com/org/my-repos.git")
		assert.Equal(t, config.Output.Out.Git.Branch, "master")
		assert.Equal(t, config.Output.Out.Git.FilePath, "README.md")
	}
}
