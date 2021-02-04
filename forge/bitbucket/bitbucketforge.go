package bitbucketforge

import (
	"fmt"
	"git-forge/cmds"
	"git-forge/config"
	"git-forge/forge"
	"git-forge/log"
	"github.com/ktrysmt/go-bitbucket"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/config"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func init() {
	err := cmds.RegisterForgeType("bitbucket", NewBitBucketForge)
	if err != nil {
		logging.Forgelog.Printf("Unable to register: %s\n", err)
	}
}

type BitBucketForge struct {
}

func NewBitBucketForge() forge.Forge {
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
	lrepo, clonerr := git.PlainClone("./"+dirname, false, &git.CloneOptions{
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

		parentCloneUrl := repo.Parent.Links["html"].(map[string]interface{})["href"].(string)

		rConfig := &config.RemoteConfig{
			Name: "origin-parent",
			URLs: []string{parentCloneUrl},
		}

		logging.Forgelog.Printf("Adding parent of origin as origin-parent\n")
		_, remerr := lrepo.CreateRemote(rConfig)
		if remerr != nil {
			f.cleanup(dirname)
			return remerr
		}

		cfg, ferr := gitconfig.NewForgeConfig("./" + dirname + "/.git/config")
		if ferr != nil {
			f.cleanup(dirname)
			return ferr
		}
		cerr := cfg.AddForgeRemoteSection(opts.ForgeName, "origin", "origin-parent")
		if cerr != nil {
			f.cleanup(dirname)
			return cerr
		}
		defer cfg.CommitConfig()
	}
	return nil
}

func (f *BitBucketForge) Fork(opts forge.ForkOpts) error {

	c := bitbucket.NewBasicAuth(opts.Common.User, opts.Common.Pass)
	// Indicate what repo we want
	_, slug, owner, _ := getRepoSlugAndOwner(opts.Url)
	bopts := &bitbucket.RepositoryOptions{
		RepoSlug: slug,
		Owner:    owner,
	}

	repo, rerr := c.Repositories.Repository.Get(bopts)
	if rerr != nil {
		return rerr
	}
	if repo == nil {
		return fmt.Errorf("Unable to find repository %s/%s\n", owner, slug)
	}

	fConfig := &bitbucket.RepositoryForkOptions{
		FromOwner: owner,
		FromSlug:  slug,
	}

	frepo, ferr := c.Repositories.Repository.Fork(fConfig)
	if ferr != nil {
		return ferr
	}
	logging.Forgelog.Printf("Forked from repo %s to repo %s/%s\n", frepo.Parent.Links["html"].(map[string]interface{})["href"].(string), user, slug)
	return nil
}
