package dummyforge

import (
	"git-forge/configset"
	"git-forge/forge"
)

type DummyForge struct {
	cfg *forge.ForgeConfig
}

func NewDummyForge(cfg *forge.ForgeConfig) forge.Forge {
	return &DummyForge{cfg}

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
		"git@dummy.org",
		"api.dummy.org/2.0",
		"UNUSEDUSER",
		"UNUSEDPAS")

	return err
}

func (f *DummyForge) Clone(opts forge.CloneOpts) error {
	return nil
}

func (f *DummyForge) Fork(opts forge.ForkOpts) error {
	return nil
}

func (f *DummyForge) CreatePr(opts forge.CreatePrOpts) error {
	return nil
}
