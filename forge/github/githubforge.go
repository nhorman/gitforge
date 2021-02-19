package githubforge

import (
	//"fmt"
	"context"
	"git-forge/cmds"
	"git-forge/config"
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
}

func NewGitHubForge() forge.Forge {
	return &GitHubForge{&forge.ForgeObj{}}

}

func getRepoSlug(url string) (string, error) {
	base := path.Base(url)
	slug := strings.TrimSuffix(base, filepath.Ext(base))
	return slug, nil
}

func (f *GitHubForge) cleanup(dirname string) error {
	return os.RemoveAll(dirname)
}

func (f *GitHubForge) InitForges(config *gitconfig.ForgeConfig) error {

	// We want to register the standard forges for bitbucket org as both
	// a git@ prefix and an https prefix
	config.AddForge("github-ssh",
		"github",
		"git@github.com",
		"api.github.com",
		"USERNAMEHERE",
		"PASSWORDHERE")

	config.AddForge("github-https",
		"bitbucket",
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
			Username: opts.Common.User,
			Password: opts.Common.Pass,
		}

		client := github.NewClient(transport.Client())
		ctx := context.Background()
		slug, _ := getRepoSlug(opts.Url)
		repo, _, err := client.Repositories.Get(ctx, opts.Common.User, slug)
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

		ferr := f.forge.CreateForgeConfig(opts.ForgeName, "origin", "origin-parent")
		if ferr != nil {
			f.cleanup(dirname)
			return ferr
		}
	}
	return nil
}

func (f *GitHubForge) Fork(opts forge.ForkOpts) error {

	transport := &github.BasicAuthTransport{
		Username: opts.Common.User,
		Password: opts.Common.Pass,
	}

	client := github.NewClient(transport.Client())
	ctx := context.Background()

	slug, _ := getRepoSlug(opts.Url)
	_, resp, err := client.Repositories.CreateFork(ctx, opts.Common.User, slug, &github.RepositoryCreateForkOptions{})
	if err != nil {
		if resp.StatusCode != 202 {
			return err
		}
	}
	logging.Forgelog.Printf("Forked from repo %s to repo %s/%s\n", slug, opts.Common.User, slug)
	return nil
}

func (f *GitHubForge) CreatePr(opts forge.CreatePrOpts) error {
	_, err := f.forge.OpenLocalRepo()
	if err != nil {
		return err
	}

	cfg, err := f.forge.GetForgeConfig()
	if err != nil {
		return err
	}

	err = f.forge.Push(cfg.Child.Name, opts.Sbranch, opts.Tbranch)
	if err != nil {
		return err
	}

	return nil
}
