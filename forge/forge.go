package forge

import ()

type Forge interface {
	Clone(createFork bool, attachFork bool, url string) error
}
