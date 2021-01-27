package bitbucketforge

import (
	"git-forge/forge"
	"git-forge/log"
	"github.com/ktrysmt/go-bitbucket"
	"gopkg.in/src-d/go-git.v4"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type BitBucketForge struct {
}

func NewBitBucketForge() *BitBucketForge {
	return &BitBucketForge{}

}

func getRepoSlugAndOwner(url string) (string, string, error) {
	base := path.Base(url)
	slug := strings.TrimSuffix(base, filepath.Ext(base))
	noslugurl := strings.TrimSuffix(url, base)
	owner := path.Base(noslugurl)
	// Need to see if we need to trim any git@ crap from the owner string
	return slug, owner, nil
}

func (f *BitBucketForge) Clone(opts forge.CloneOpts) error {
	logging.Forgelog.Printf("%s appears to be a bitbucket forge\n", opts.Url)

	dirname := strings.TrimSuffix(path.Base(opts.Url), filepath.Ext(path.Base(opts.Url)))
	err := os.Mkdir("./"+dirname, 0755)
	if err != nil {
		return err
	}

	// Start by cloning the repository requested
	_, clonerr := git.PlainClone("./"+dirname, false, &git.CloneOptions{
		URL: opts.Url})
	if clonerr != nil {
		return clonerr
	}

	// now get us our auth token for the bitbucket api
	c := bitbucket.NewBasicAuth(opts.Common.User, opts.Common.Pass)
	// Indicate what repo we want
	slug, owner, _ := getRepoSlugAndOwner(opts.Url)
	bopts := &bitbucket.RepositoryOptions{
		RepoSlug: slug,
		Owner:    owner,
	}

	_, rerr := c.Repositories.Repository.Get(bopts)
	if rerr != nil {
		return rerr
	}
	return nil
}

func (f *BitBucketForge) Fork(opts forge.ForkOpts) error {
	return nil
}
