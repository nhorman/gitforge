package dummyforge

import (
	"git-forge/config"
	"git-forge/forge"
)

type DummyForge struct {
}

func NewDummyForge() forge.Forge {
	return &DummyForge{}

}

func (f *DummyForge) InitForges(config *gitconfig.ForgeConfig) error {

	// We want to register the standard forges for bitbucket org as both
	// a git@ prefix and an https prefix
	config.AddForge("dummy-ssh",
		"dummy",
		"git@dummy.org",
		"api.dummy.org/2.0",
		"UNUSEDUSER",
		"UNUSEDPAS")

	return nil
}

func (f *DummyForge) Clone(opts forge.CloneOpts) error {
	return nil
}

func (f *DummyForge) Fork(opts forge.ForkOpts) error {
	return nil
}
