package asc

import (
	in "github.com/nagaa052/asc-gen/pkg/input"
	out "github.com/nagaa052/asc-gen/pkg/output"
)

// Asc ...
type Asc struct {
	c *AssConfig
}

// New ...
func New(config *AssConfig) (*Asc, error) {
	return &Asc{
		c: config,
	}, nil
}

// Run ...
func (asc *Asc) Run() error {
	in, err := in.NewClient(asc.c.Input)
	if err != nil {
		return err
	}

	items, err := in.GetItems()
	if err != nil {
		return err
	}

	buf, err := BindDocument(asc.c.Output.Title, items, asc.c.TemplatePath, asc.c.Output.Labels)
	if err != nil {
		return err
	}

	dc, err := out.NewClient(asc.c.Output.Out)
	if err != nil {
		return err
	}

	return dc.Publishs(buf)
}
