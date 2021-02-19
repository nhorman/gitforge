package forge

import (
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/config"
	"os"

	"git-forge/log"
)

type CloneOpts struct {
	Parentfork bool
	Url        string
}

type ForkOpts struct {
	Url       string
	ForgeName string
}

type CreatePrOpts struct {
	Remote      string
	Sbranch     string
	Tbranch     string
	Title       string
	Description string
}

type ForgeConfig struct {
	Name         string
	Type         string
	User         string
	Pass         string
	CloneBaseUrl string
	ApiBaseUrl   string
}

type Forge interface {
	InitForges() error
	Clone(opts CloneOpts) error
	Fork(opts ForkOpts) error
	CreatePr(opts CreatePrOpts) error
}

type LocalRepoOps interface {
	CreateLocalRepo(name string, bare bool, url string) (*git.Repository, error)
	OpenLocalRepo() (*git.Repository, error)
	Push(remote string, sbranch string, tbranch string) error
	CreateRemote(name string, url string) (*git.Remote, error)
	CreateForgeConfig() error
	GetForgeConfig() (ForgeConfig, error)
}

type ForgeObj struct {
	dir   string
	Lrepo *git.Repository
}

type RemoteInfo struct {
	Name string
	Url  string
}

type ForgeInfo struct {
	Parent RemoteInfo
	Child  RemoteInfo
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

func (f *ForgeObj) OpenLocalRepo() (*git.Repository, error) {
	lrepo, err := git.PlainOpenWithOptions("./", &git.PlainOpenOptions{DetectDotGit: true})
	if err != nil {
		return nil, err
	}
	f.dir = "./"
	f.Lrepo = lrepo
	return lrepo, err
}

func (f *ForgeObj) Push(remote string, sbranch string, tbranch string) error {
	refspec := sbranch + ":" + tbranch

	return f.Lrepo.Push(&git.PushOptions{
		RemoteName: remote,
		RefSpecs:   []config.RefSpec{config.RefSpec(refspec)},
		Prune:      false,
	})
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
