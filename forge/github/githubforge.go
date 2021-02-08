package githubforge

import (
	//"fmt"
	"git-forge/cmds"
	"git-forge/config"
	"git-forge/forge"
	"git-forge/log"

	//"github.com/google/go-github/v32/github"
	"gopkg.in/src-d/go-git.v4"
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
}

func NewGitHubForge() forge.Forge {
	return &GitHubForge{}

}

//func getRepoSlugAndOwner(url string) (string, string, string, error) {
//	var owner string
//	base := path.Base(url)
//	slug := strings.TrimSuffix(base, filepath.Ext(base))
//	noslugurl := strings.TrimSuffix(url, base)
//	owner = path.Base(noslugurl)
//
//	// Need to see if we need to trim any git@ crap from the owner string
//	if strings.HasPrefix(owner, "git@") == true {
//		// we have to chop off everthing up to the ':'
//		idx := strings.Index(owner, ":")
//		owner = owner[idx+1:]
//	}
//	finalbaseurl := strings.TrimSuffix(noslugurl, owner)
//	return finalbaseurl, slug, owner, nil
//}

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
	_, clonerr := git.PlainClone("./"+dirname, false, &git.CloneOptions{
		URL: opts.Url})
	if clonerr != nil {
		return clonerr
	}

	if opts.Parentfork == true {
	}
	return nil
}

func (f *GitHubForge) Fork(opts forge.ForkOpts) error {

	return nil
}
