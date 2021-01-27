package forge

import ()

type CommonOpts struct {
	user string
	pass string
}

type CloneOpts struct {
	CommonOpts
	Parentfork bool
	Url        string
}

type ForkOpts struct {
	CommonOpts
	Url string
}

type Forge interface {
	Clone(opts CloneOpts) error
	Fork(opts ForkOpts) error
}
