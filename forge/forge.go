package forge

import (
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/config"
	"os"

	"git-forge/config"
	"git-forge/log"
)

type CommonOpts struct {
	User string
	Pass string
}

type CloneOpts struct {
	Common     CommonOpts
	Parentfork bool
	Url        string
	ForgeName  string
}

type ForkOpts struct {
	Common    CommonOpts
	Url       string
	ForgeName string
}

type Forge interface {
	InitForges(config *gitconfig.ForgeConfig) error
	Clone(opts CloneOpts) error
	Fork(opts ForkOpts) error
}

type LocalRepoOps interface {
	CreateLocalRepo(name string, bare bool, url string) (*git.Repository, error)
	CreateRemote(name string, url string) (*git.Remote, error)
	CreateForgeConfig() error
}

type ForgeObj struct {
	dir   string
	Lrepo *git.Repository
}

//
// LocalRepOps Implementation
//
func (f *ForgeObj) CreateLocalRepo(name string, bare bool, url string) (*git.Repository, error) {
	err := os.Mkdir("./"+name, 0755)
	if err != nil {
		return nil, err
	}
	lrepo, clonerr := git.PlainClone("./"+name, bare, &git.CloneOptions{
		URL: url})
	if clonerr != nil {
		return nil, clonerr
	}
	f.Lrepo = lrepo
	f.dir = name
	return lrepo, nil
}

func (f *ForgeObj) CreateRemote(name string, url string) (*git.Remote, error) {

	rConfig := &config.RemoteConfig{
		Name: name,
		URLs: []string{url},
	}

	logging.Forgelog.Printf("Adding remote %s\n", name)
	lremote, remerr := f.Lrepo.CreateRemote(rConfig)
	if remerr != nil {
		return nil, remerr
	}
	return lremote, nil
}

func (f *ForgeObj) CreateForgeConfig(ftype string, child string, parent string) error {
	cfg, ferr := gitconfig.NewForgeConfig("./" + f.dir + "/.git/config")
	if ferr != nil {
		return ferr
	}
	cerr := cfg.AddForgeRemoteSection(ftype, child, parent)
	if cerr != nil {
		return cerr
	}
	defer cfg.CommitConfig()
	return nil
}
