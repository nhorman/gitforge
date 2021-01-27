package forge

import ()

type CloneOpts struct {
	Parentfork bool
	Url        string
}

type ForkOpts struct {
	Url string
}

type Forge interface {
	Clone(opts CloneOpts) error
	Fork(opts ForkOpts) error
}
