package asc

import (
	"github.com/nagaa052/asc-gen/pkg/input"
	"github.com/nagaa052/asc-gen/pkg/output"
)

// AssConfig ...
type AssConfig struct {
	TemplatePath string        `yaml:"template_path"`
	Input        *input.Config `yaml:"input"`
	Output       struct {
		Labels []string       `yaml:"labels"`
		Title  string         `yaml:"title"`
		Out    *output.Config `yaml:"out"`
	} `yaml:"output"`
}
