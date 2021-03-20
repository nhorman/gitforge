package dummyforge

import (
	"fmt"
	"git-forge/configset"
	"git-forge/forge"
	"git-forge/log"
	"os"
	"path/filepath"
)

type DummyForge struct {
	forge *forge.ForgeObj
	cfg   *forge.ForgeConfig
}

func NewDummyForge(cfg *forge.ForgeConfig) forge.Forge {
	return &DummyForge{&forge.ForgeObj{}, cfg}

}

func (f *DummyForge) InitForges() error {

	config, cerr := gitconfigset.NewForgeConfigSet()
	if cerr != nil {
		return cerr
	}
	defer config.CommitConfig()

	// We want to register the standard forges for bitbucket org as both
	// a git@ prefix and an https prefix
	err := config.AddForge("dummy-ssh",
		"dummy",
		"file:///",
		"api.dummy.org/2.0",
		"UNUSEDUSER",
		"UNUSEDPAS")

	return err
}

func (f *DummyForge) Clone(opts forge.CloneOpts) error {
	logging.Forgelog.Printf("%s appears to be a dummy forge\n", opts.Url)

	dirname := "testrepo"

	wd, _ := os.Getwd()

	dir, _ := filepath.Abs(wd + "../../..")

	opts.Url = "file:///" + dir

	// Start by cloning the repository requested
	_, clonerr := f.forge.CreateLocalRepo(dirname, false, opts.Url)
	if clonerr != nil {
		return fmt.Errorf("Unable to clone %s: %s\n", opts.Url, clonerr)
	}

	if opts.Parentfork == true {
		_, remerr := f.forge.CreateRemote("origin-parent", opts.Url)
		if remerr != nil {
			return remerr
		}

		cfg, err := gitconfigset.NewForgeConfigSetInDir(dirname)
		if err != nil {
			return err
		}
		defer cfg.CommitConfig()

		ferr := cfg.AddForgeRemoteSection(f.cfg.Type, "origin", "origin-parent")
		if ferr != nil {
			return ferr
		}
	}

	// dummy command side effects, after we finish the clone, we need to cd
	// into the directory for subsequent commands
	os.Chdir("testrepo")
	return nil
}

func (f *DummyForge) Fork(opts forge.ForkOpts) error {
	return nil
}

func (f *DummyForge) CreatePr(opts forge.CreatePrOpts) error {
	return nil
}

func (f *DummyForge) GetAllPrTitles() ([]forge.PrTitle, error) {
	return nil, nil
}

func (f *DummyForge) GetPr(idstring string) (*forge.PR, error) {
	return nil, nil
}

func (f *DummyForge) RefreshPr(pr *forge.PR) (chan *forge.UpdatedPR, error) {
	return nil, nil
}
