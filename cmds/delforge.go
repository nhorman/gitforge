package cmds

import (
	"flag"
	"fmt"
	"git-forge/config"
	"os"
)

func Delusage() {
	fmt.Printf("Usage: git forge delforge [options]\n")
	fmt.Printf("Description: Remove a forge type to the global gitconfig\n")
	fmt.Printf("Options:\n")
	flag.PrintDefaults()
}

func DelForgeCmd() error {

	helpopt := flag.Bool("help", false, "display help for addforge command")
	nameopt := flag.String("name", "", "Name of the forge to delete")
	flag.Parse()

	gitconfigpath := os.Getenv("HOME") + "/.gitconfig"
	if *helpopt == true {
		Delusage()
		return nil
	}

	forgeconfig, err := gitconfig.NewForgeConfig(gitconfigpath)
	if err != nil {
		return fmt.Errorf("Create forge config failed: %s\n", err)
	}
	defer forgeconfig.CommitConfig()

	ferr := forgeconfig.DelForge(*nameopt)
	if ferr != nil {
		return fmt.Errorf("Failed to add forge: %s\n", ferr)
	}

	return nil
}
