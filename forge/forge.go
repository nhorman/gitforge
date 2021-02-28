package forge

import (
	"fmt"
	"git-forge/log"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/config"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
	"os"
)

type ForgeConfig struct {
	Name         string
	Type         string
	User         string
	Pass         string
	CloneBaseUrl string
	ApiBaseUrl   string
}

type CloneOpts struct {
	Parentfork bool
	Url        string
}

type ForkOpts struct {
	Url       string
	ForgeName string
}

type CreatePrOpts struct {
	Sbranch     string
	Tbranch     string
	Title       string
	Description string
}

type ForgeRemoteInfo struct {
	Name string
	Url  string
}

type ForgeLocalConfig struct {
	Type   string
	Child  ForgeRemoteInfo
	Parent ForgeRemoteInfo
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
	Fetch(location string, refspec string) error
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
	refspec := "refs/heads/" + sbranch + ":refs/heads/" + sbranch

	options := &git.PushOptions{
		RemoteName: remote,
		RefSpecs:   []config.RefSpec{config.RefSpec(refspec)},
		Prune:      false,
	}

	verr := options.Validate()
	if verr != nil {
		return verr
	}

	ret := f.Lrepo.Push(options)

	if ret != nil {
		// Mask already up to date errors
		if ret == git.NoErrAlreadyUpToDate {
			ret = nil
		}
	}
	return ret
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

func (f *ForgeObj) Fetch(location string, refspec string, user string, pass string) error {

	copts := &config.RemoteConfig{
		Name: "_tmpfetchremote",
		URLs: []string{location},
	}

	remote, cerr := f.Lrepo.CreateRemote(copts)
	if cerr != nil {
		return fmt.Errorf("Unable to create rmote: %s", cerr)
	}
	defer f.Lrepo.DeleteRemote("_tmpfetchremote")

	fopts := &git.FetchOptions{
		RemoteName: location,
		RefSpecs:   []config.RefSpec{config.RefSpec(refspec)},
		Force:      true,
		Auth:       &http.BasicAuth{user, pass},
	}

	err := fopts.Validate()
	if err != nil {
		return fmt.Errorf("Unable to fetch: %s", err)
	}
	ret := remote.Fetch(fopts)
	if ret != nil {
		if ret == git.NoErrAlreadyUpToDate {
			return nil
		}
	}
	return ret
}
