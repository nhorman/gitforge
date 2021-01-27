package forge

import ()

type CommonOpts struct {
	User string
	Pass string
}

type CloneOpts struct {
	Common     CommonOpts
	Parentfork bool
	Url        string
}

type ForkOpts struct {
	Common CommonOpts
	Url    string
}

type Forge interface {
	Clone(opts CloneOpts) error
	Fork(opts ForkOpts) error
}
