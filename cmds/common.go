package cmds

import (
	"fmt"
	"git-forge/config"
	"git-forge/forge"
	"git-forge/forge/bitbucket"
	"os"
)

func LookupForgeType(url string) (string, error) {

	gitconfigpath := os.Getenv("HOME") + "/.gitconfig"

	forgeconfig, err := gitconfig.NewForgeConfig(gitconfigpath)
	if err != nil {
		return "", fmt.Errorf("Lookup forge config failed: %s\n", err)
	}

	return forgeconfig.LookupForge(url)
}

func AllocateForgeFromUrl(url string) (forge.Forge, error) {

	var forge forge.Forge = nil

	ftype, err := LookupForgeType(url)

	if err != nil {
		return nil, err
	}

	switch ftype {
	case "bitbucket":
		// Allocate a bitbucket Forge
		forge = bitbucketforge.NewBitBucketForge()
	default:
		return nil, fmt.Errorf("This build does not support forge type %s\n", ftype)
	}

	return forge, nil
}
