package asc

import (
	"bytes"
	"text/template"

	in "github.com/nagaa052/asc-gen/pkg/input"
)

type Label string

type Document struct {
	Title   string
	Content map[Label][]in.Item
}

// BindDocument ...
func BindDocument(title string, items []*in.Item, tplPath string, labels []string) (*bytes.Buffer, error) {
	doc := inItemToDocument(title, items, labels)
	buf := new(bytes.Buffer)
	t, err := template.ParseFiles(tplPath)
	if err != nil {
		return nil, err
	}

	if err := t.Execute(buf, doc); err != nil {
		return nil, err
	}

	return buf, nil
}

func inItemToDocument(title string, items []*in.Item, labels []string) *Document {
	cl := make(map[Label][]in.Item)
	set := make(map[string]struct{}, len(labels))
	for _, s := range labels {
		set[s] = struct{}{}
	}

	for _, item := range items {
		if len(item.Labels) <= 0 {
			continue
		}

		for _, label := range item.Labels {
			if _, ok := set[label]; !ok {
				continue
			}

			key := Label(label)
			if _, ok := cl[key]; !ok {
				cl[key] = make([]in.Item, 0)
			}
			cl[key] = append(cl[key], *item)
		}
	}
	return &Document{
		Title:   title,
		Content: cl,
	}
}
