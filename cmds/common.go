package cmds

import (
	"fmt"
	"git-forge/configset"
	"git-forge/forge"
	"git-forge/log"
)

var forgetypes map[string]func(*forge.ForgeConfig) forge.Forge = make(map[string]func(*forge.ForgeConfig) forge.Forge, 0)

func RegisterForgeType(name string, ifunc func(*forge.ForgeConfig) forge.Forge) error {
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

	forgeconfig, err := gitconfigset.NewForgeConfigSet()
	if err != nil {
		return nil, err
	}

	fconfig, err := forgeconfig.ConfigFromUrl(url)

	if err != nil {
		return nil, fmt.Errorf("No forge for url %s\n", url)
	}

	ifunc, ok := forgetypes[fconfig.Type]
	if ok == true {
		forge = ifunc(fconfig)
	} else {
		return nil, fmt.Errorf("No such forge type: %s\n", fconfig.Type)
	}
	return forge, nil
}
