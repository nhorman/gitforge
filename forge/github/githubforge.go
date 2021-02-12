package githubforge

import (
	//"fmt"
	"context"
	"git-forge/cmds"
	"git-forge/config"
	"git-forge/forge"
	"git-forge/log"

	"github.com/google/go-github/v33/github"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/config"

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
}

func NewGitHubForge() forge.Forge {
	return &GitHubForge{}

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
		transport := &github.BasicAuthTransport{
			Username: opts.Common.User,
			Password: opts.Common.Pass,
		}

		client := github.NewClient(transport.Client())
		ctx := context.Background()
		slug, _ := getRepoSlug(opts.Url)
		repo, _, err := client.Repositories.Get(ctx, opts.Common.User, slug)
		if err != nil {
			return err
		}
		prepo := repo.GetParent()

		rConfig := &config.RemoteConfig{
			Name: "origin-parent",
			URLs: []string{*prepo.CloneURL},
		}

		_, remerr := lrepo.CreateRemote(rConfig)
		if remerr != nil {
			return remerr
		}

	}
	return nil
}

func (f *GitHubForge) Fork(opts forge.ForkOpts) error {

	return nil
}
