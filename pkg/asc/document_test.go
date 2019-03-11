package asc

import (
	"testing"

	in "github.com/nagaa052/asc-gen/pkg/input"
	"github.com/stretchr/testify/assert"
)

func TestBindDocument(t *testing.T) {
	title := "Test List"
	items := []*in.Item{
		&in.Item{
			Name:        "ItemA",
			Description: "is ItemA",
			URL:         "https://example.com/apc/example1.git",
			UpdatedAt:   "2019-01-01",
			Labels:      []string{"aaa", "bbb", "ccc"},
		},
		&in.Item{
			Name:        "ItemB",
			Description: "is ItemA",
			URL:         "https://example.com/apc/example1.git",
			UpdatedAt:   "2019-01-01",
			Labels:      []string{"ccc", "ddd"},
		},
		&in.Item{
			Name:        "ItemB",
			Description: "is ItemA",
			URL:         "https://example.com/apc/example1.git",
			UpdatedAt:   "2019-01-01",
			Labels:      []string{"eee"},
		},
	}
	labels := []string{"aaa", "ccc"}
	buf, err := BindDocument(title, items, "../../templates/output.md.tpl", labels)
	assert.NoError(t, err)
	assert.NotNil(t, buf)
}
