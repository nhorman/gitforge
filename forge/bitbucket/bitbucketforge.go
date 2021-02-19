package bitbucketforge

import (
	"fmt"
	"git-forge/cmds"
	"git-forge/config"
	"git-forge/forge"
	"git-forge/log"
	"github.com/ktrysmt/go-bitbucket"
	//"gopkg.in/src-d/go-git.v4"
	//"gopkg.in/src-d/go-git.v4/config"
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
	forge *forge.ForgeObj
}

func NewBitBucketForge() forge.Forge {
	return &BitBucketForge{&forge.ForgeObj{}}

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

func (f *BitBucketForge) InitForges(config *gitconfig.ForgeConfig) error {

	// We want to register the standard forges for bitbucket org as both
	// a git@ prefix and an https prefix
	config.AddForge("bitbucket-ssh",
		"bitbucket",
		"git@bitbucket.org",
		"api.bitbucket.org/2.0",
		"USERNAMEHERE",
		"PASSWORDHERE")

	config.AddForge("bitbucket-https",
		"bitbucket",
		"https://USERNAME@bitbucket.org",
		"api.bitbucket.org/2.0",
		"USERNAMEHERE",
		"PASSWORDHERE")

	return nil
}

func (f *BitBucketForge) Clone(opts forge.CloneOpts) error {
	logging.Forgelog.Printf("%s appears to be a bitbucket forge\n", opts.Url)

	dirname := strings.TrimSuffix(path.Base(opts.Url), filepath.Ext(path.Base(opts.Url)))

	// Start by cloning the repository requested
	_, clonerr := f.forge.CreateLocalRepo(dirname, false, opts.Url)
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

		_, remerr := f.forge.CreateRemote("origin-remote", parentCloneUrl)
		if remerr != nil {
			f.cleanup(dirname)
			return remerr
		}

		ferr := f.forge.CreateForgeConfig(opts.ForgeName, "origin", "origin-parent")
		if ferr != nil {
			f.cleanup(dirname)
			return ferr
		}
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
	logging.Forgelog.Printf("Forked from repo %s to repo %s/%s\n", frepo.Parent.Links["html"].(map[string]interface{})["href"].(string), owner, slug)
	return nil
}

func (f *BitBucketForge) CreatePr(opts forge.CreatePrOpts) error {

	_, err := f.forge.OpenLocalRepo()
	if err != nil {
		return err
	}

	logging.Forgelog.Printf("Synchronizing branches\n")

	err = f.forge.Push(opts.Remote, opts.Sbranch, opts.Tbranch)
	if err != nil {
		return fmt.Errorf("Push Failed: %s\n", err)
	}

	fcfg, ferr := f.forge.GetForgeConfig()
	if ferr != nil {
		return fmt.Errorf("Forge config is busted: %s\n", ferr)
	}

	c := bitbucket.NewBasicAuth(opts.Common.User, opts.Common.Pass)
	// Indicate what repo we want
	_, slug, _, _ := getRepoSlugAndOwner(fcfg.Child.Url)

	propts := &bitbucket.PullRequestsOptions{
		Owner:             opts.Common.User,
		RepoSlug:          slug,
		SourceBranch:      opts.Sbranch,
		DestinationBranch: opts.Tbranch,
	}

	_, cerr := c.Repositories.PullRequests.Create(propts)
	if cerr != nil {
		return fmt.Errorf("Unable to create Pull Request: %s\n", cerr)
	}

	return nil
}
