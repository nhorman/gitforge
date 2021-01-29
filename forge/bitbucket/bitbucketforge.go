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

func getRepoSlugAndOwner(url string) (string, string, string, error) {
	var owner string
	base := path.Base(url)
	slug := strings.TrimSuffix(base, filepath.Ext(base))
	noslugurl := strings.TrimSuffix(url, base)
	owner = path.Base(noslugurl)

	// Need to see if we need to trim any git@ crap from the owner string
	if strings.HasPrefix(owner, "git@") == true {
		// we have to chop off everthing up to the ':'
		idx := strings.Index(owner, ":")
		owner = owner[idx+1:]
	}
	finalbaseurl := strings.TrimSuffix(noslugurl, owner)
	return finalbaseurl, slug, owner, nil
}

func (f *BitBucketForge) cleanup(dirname string) error {
	return os.RemoveAll(dirname)
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

	if opts.Parentfork == true {
		// now get us our auth token for the bitbucket api
		c := bitbucket.NewBasicAuth(opts.Common.User, opts.Common.Pass)
		// Indicate what repo we want
		_, slug, owner, _ := getRepoSlugAndOwner(opts.Url)
		bopts := &bitbucket.RepositoryOptions{
			RepoSlug: slug,
			Owner:    owner,
		}

		repo, rerr := c.Repositories.Repository.Get(bopts)
		if rerr != nil {
			f.cleanup(dirname)
			return rerr
		}

		parentFullName := repo.Parent.Full_name
		parentSlug := path.Base(parentFullName)
		parentOwner := strings.TrimSuffix(parentFullName, "/"+parentSlug)
		logging.Forgelog.Printf("%s %s\n", parentSlug, parentOwner)
	}
	return nil
}

func (f *BitBucketForge) Fork(opts forge.ForkOpts) error {
	return nil
}
