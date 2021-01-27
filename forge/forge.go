package forge

import ()

type Forge interface {
	Clone(parentfork bool, url string) error
}
