package forge

import (
	"git-forge/config"
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
