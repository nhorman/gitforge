package githubforge

import (
	//"fmt"
	"context"
	"git-forge/cmds"
	"git-forge/configset"
	"git-forge/forge"
	"git-forge/log"

	"github.com/google/go-github/v33/github"
	//"gopkg.in/src-d/go-git.v4"
	//"gopkg.in/src-d/go-git.v4/config"

	"os"
	"path"
	"path/filepath"
	"strings"
)

func init() {
	err := cmds.RegisterForgeType("github", NewGitHubForge)
	if err != nil {
		logging.Forgelog.Printf("Unable to register: %s\n", err)
	}
}

type GitHubForge struct {
	forge *forge.ForgeObj
	cfg   *forge.ForgeConfig
}

func NewGitHubForge(cfg *forge.ForgeConfig) forge.Forge {
	return &GitHubForge{&forge.ForgeObj{}, cfg}

}

func getRepoSlug(url string) (string, error) {
	base := path.Base(url)
	slug := strings.TrimSuffix(base, filepath.Ext(base))
	return slug, nil
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

func (f *GitHubForge) cleanup(dirname string) error {
	return os.RemoveAll(dirname)
}

func (f *GitHubForge) InitForges() error {

	config, err := gitconfigset.NewForgeConfigSet()
	if err != nil {
		return nil
	}
	defer config.CommitConfig()

	// We want to register the standard forges for bitbucket org as both
	// a git@ prefix and an https prefix
	config.AddForge("github-ssh",
		"github",
		"git@github.com",
		"api.github.com",
		"USERNAMEHERE",
		"PASSWORDHERE")

	config.AddForge("github-https",
		"github",
		"https://github.com",
		"api.github.com",
		"USERNAMEHERE",
		"PASSWORDHERE")

	return nil
}

func (f *GitHubForge) Clone(opts forge.CloneOpts) error {
	logging.Forgelog.Printf("%s appears to be a github forge\n", opts.Url)

	dirname := strings.TrimSuffix(path.Base(opts.Url), filepath.Ext(path.Base(opts.Url)))

	// Start by cloning the repository requested
	_, clonerr := f.forge.CreateLocalRepo(dirname, false, opts.Url)
	if clonerr != nil {
		return clonerr
	}

	if opts.Parentfork == true {
		transport := &github.BasicAuthTransport{
			Username: f.cfg.User,
			Password: f.cfg.Pass,
		}

		client := github.NewClient(transport.Client())
		ctx := context.Background()
		slug, _ := getRepoSlug(opts.Url)
		repo, _, err := client.Repositories.Get(ctx, f.cfg.User, slug)
		if err != nil {
			f.cleanup(dirname)
			return err
		}
		prepo := repo.GetParent()

		_, remerr := f.forge.CreateRemote("origin-parent", *prepo.CloneURL)
		if remerr != nil {
			f.cleanup(dirname)
			return remerr
		}

		cfg, err := gitconfigset.NewForgeConfigSetInDir(dirname)
		if err != nil {
			f.cleanup(dirname)
			return err
		}
		defer cfg.CommitConfig()

		ferr := cfg.AddForgeRemoteSection(f.cfg.Type, "origin", "origin-parent")
		if ferr != nil {
			f.cleanup(dirname)
			return ferr
		}

	}
	return nil
}

func (f *GitHubForge) Fork(opts forge.ForkOpts) error {

	transport := &github.BasicAuthTransport{
		Username: f.cfg.User,
		Password: f.cfg.Pass,
	}

	client := github.NewClient(transport.Client())
	ctx := context.Background()

	slug, _ := getRepoSlug(opts.Url)
	_, resp, err := client.Repositories.CreateFork(ctx, f.cfg.User, slug, &github.RepositoryCreateForkOptions{})
	if err != nil {
		if resp.StatusCode != 202 {
			return err
		}
	}
	logging.Forgelog.Printf("Forked from repo %s to repo %s/%s\n", slug, f.cfg.User, slug)
	return nil
}

func (f *GitHubForge) CreatePr(opts forge.CreatePrOpts) error {
	_, err := f.forge.OpenLocalRepo()
	if err != nil {
		return err
	}

	cfg, err := gitconfigset.NewForgeConfigSet()
	if err != nil {
		return err
	}
	defer cfg.CommitConfig()
	fconfig, err := cfg.GetForgeRemoteSection()
	if err != nil {
		return err
	}

	err = f.forge.Push(fconfig.Child.Name, opts.Sbranch, opts.Tbranch)
	if err != nil {
		return err
	}

	transport := &github.BasicAuthTransport{
		Username: f.cfg.User,
		Password: f.cfg.Pass,
	}

	client := github.NewClient(transport.Client())
	ctx := context.Background()
	_, pslug, powner, _ := getRepoSlugAndOwner(fconfig.Parent.Url)
	_, _, cowner, _ := getRepoSlugAndOwner(fconfig.Child.Url)

	newPR := &github.NewPullRequest{
		Title:               github.String(opts.Title),
		Head:                github.String(cowner + ":" + opts.Sbranch),
		Base:                github.String(opts.Tbranch),
		Body:                github.String(opts.Description),
		MaintainerCanModify: github.Bool(true),
	}

	_, _, err = client.PullRequests.Create(ctx, powner, pslug, newPR)
	if err != nil {
		return err
	}

	return nil
}
