package cmds

import (
	"flag"
	"fmt"
	"git-forge/config"
	"git-forge/log"
	"os"
)

func Delusage() {
	logging.Forgelog.Printf("Usage: git forge delforge [options]\n")
	logging.Forgelog.Printf("Description: Remove a forge type to the global gitconfig\n")
	logging.Forgelog.Printf("Options:\n")
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

	forgeconfig, err := gitconfig.GetForgeConfig(gitconfigpath, *nameopt)
	if err != nil {
		return fmt.Errorf("Create forge config failed: %s\n", err)
	}
	defer forgeconfig.CommitConfig()

	ferr := forgeconfig.DelForge()
	if ferr != nil {
		return fmt.Errorf("Failed to add forge: %s\n", ferr)
	}

	return nil
}
