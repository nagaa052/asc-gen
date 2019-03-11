package output

import (
	"bytes"
	"log"
	"os"
	"time"

	"gopkg.in/src-d/go-git.v4/plumbing"

	"github.com/pkg/errors"
	billy "gopkg.in/src-d/go-billy.v4"
	"gopkg.in/src-d/go-billy.v4/memfs"
	gg "gopkg.in/src-d/go-git.v4"
	gconf "gopkg.in/src-d/go-git.v4/config"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
	"gopkg.in/src-d/go-git.v4/storage/memory"
)

// Git ...
type Git struct {
	config  *GitConfig
	repo    *gg.Repository
	fs      billy.Filesystem
	clone   func(g *Git) error
	rewrite func(g *Git, buf *bytes.Buffer) error
	commit  func(g *Git) error
	push    func(g *Git) error
}

// GitConfig ...
type GitConfig struct {
	URL          string `yaml:"url"`
	Branch       string `yaml:"branch"`
	FilePath     string `yaml:"filepath"`
	Password     string
	AuthorUserID string
	AuthorEmail  string
}

// NewGitClient ...
func NewGitClient(config *GitConfig) (*Git, error) {
	if config.Branch == "" {
		config.Branch = "master"
	}
	if config.FilePath == "" {
		config.FilePath = "README.md"
	}
	git := &Git{
		config:  config,
		clone:   clone,
		rewrite: rewrite,
		commit:  commit,
		push:    push,
	}
	if !validGit(git) {
		return nil, errors.New("validation failed git client")
	}
	return git, nil
}

// Publish ...
func (g *Git) Publish(buf *bytes.Buffer) error {
	if err := clone(g); err != nil {
		return err
	}
	if err := rewrite(g, buf); err != nil {
		return err
	}
	if err := commit(g); err != nil {
		return err
	}
	return push(g)
}

func validGit(gh *Git) bool {
	if gh.config.URL == "" {
		return false
	}
	if gh.config.Password == "" {
		return false
	}
	if gh.config.AuthorUserID == "" {
		return false
	}
	if gh.config.AuthorEmail == "" {
		return false
	}
	return true
}

func clone(g *Git) error {
	g.fs = memfs.New()
	log.Printf("clone display repository")
	var err error
	g.repo, err = gg.Clone(memory.NewStorage(), g.fs, &gg.CloneOptions{
		Auth:          &http.BasicAuth{Username: g.config.AuthorEmail, Password: g.config.Password},
		URL:           g.config.URL,
		ReferenceName: plumbing.ReferenceName("refs/heads/master"),
		SingleBranch:  true,
		Progress:      os.Stdout,
	})
	if err != nil {
		return err
	}
	if g.config.Branch != "master" {
		w, err := g.repo.Worktree()
		err = w.Checkout(&gg.CheckoutOptions{
			Create: true,
			Branch: plumbing.ReferenceName(g.config.Branch),
		})
		if err != nil {
			return errors.Wrapf(err, "failed to checkout branch")
		}
	}
	return nil
}

func rewrite(g *Git, buf *bytes.Buffer) error {
	if g.repo == nil {
		g.clone(g)
	}

	log.Printf("rewrite document\n")
	err := g.fs.Remove(g.config.FilePath)
	if err != nil {
		return errors.Wrapf(err, "failed to remove file")
	}

	file, err := g.fs.OpenFile(g.config.FilePath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return errors.Wrapf(err, "failed to open file")
	}
	defer file.Close()

	file.Write(buf.Bytes())
	file.Write([]byte("\n"))

	return nil
}

func commit(g *Git) error {
	if g.repo == nil {
		g.clone(g)
	}
	log.Printf("display add files")
	w, _ := g.repo.Worktree()
	w.Add(".")
	hash, err := w.Commit("reflesh display", &gg.CommitOptions{
		Author: &object.Signature{
			Name:  g.config.AuthorUserID,
			Email: g.config.AuthorEmail,
			When:  time.Now(),
		},
	})
	if err != nil {
		return errors.Wrapf(err, "display commit failed.")
	}
	log.Printf("display commit. hash : %s\n", hash)
	g.repo.Storer.SetReference(plumbing.NewReferenceFromStrings(g.config.Branch, hash.String()))
	return nil
}

func push(g *Git) error {
	if g.repo == nil {
		g.clone(g)
	}
	remote, err := g.repo.Remote("origin")
	if err != nil {
		return errors.Wrapf(err, "display remote add failed.")
	}
	ref := plumbing.ReferenceName(g.config.Branch)
	err = remote.Push(&gg.PushOptions{
		Auth:     &http.BasicAuth{Username: g.config.AuthorEmail, Password: g.config.Password},
		Progress: os.Stdout,
		RefSpecs: []gconf.RefSpec{
			gconf.RefSpec(ref + ":" + plumbing.ReferenceName("refs/heads/"+g.config.Branch)),
		},
	})
	if err != nil {
		return errors.Wrapf(err, "display push failed.")
	}
	log.Printf("display repository push complete.")
	return nil
}
