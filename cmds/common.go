package cmds

import (
	"fmt"
	"git-forge/config"
	"git-forge/forge"
	"git-forge/log"
)

var forgetypes map[string]func() forge.Forge = make(map[string]func() forge.Forge, 0)

func RegisterForgeType(name string, ifunc func() forge.Forge) error {
	if _, ok := forgetypes[name]; ok {
		return fmt.Errorf("%s is already registered as a forge type\n", name)
	}
	forgetypes[name] = ifunc
	return nil
}

func PrintForgeTypes() {
	logging.Forgelog.Printf("Available Forge Types:\n")
	for key, _ := range forgetypes {
		logging.Forgelog.Printf("%s\n", key)
	}
}

func AllocateForgeFromUrl(url string) (forge.Forge, error) {

	var forge forge.Forge = nil

	ftype, err := gitconfig.LookupForgeType(url)

	if err != nil {
		return nil, err
	}

	ifunc, ok := forgetypes[ftype]
	if ok == true {
		forge = ifunc()
	} else {
		return nil, fmt.Errorf("No such forge type: %s\n", ftype)
	}
	return forge, nil
}
