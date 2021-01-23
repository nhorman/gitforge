package cmds

import (
	"flag"
	"fmt"
	"git-forge/config"
)

func usage() {
	fmt.Printf("Usage: git forge addforge [options]\n")
	fmt.Printf("Description: Add a forge type to the global gitconfig\n")
	fmt.Printf("Options:\n")
	fmt.Printf("\t -h : This message\n")
}

func AddForgeCmd() error {

	helpopt := flag.Bool("help", false, "display help for addforge command")
	flag.Parse()

	if *helpopt == true {
		usage()
		return nil
	}

	forgeconfig, err := gitconfig.NewForgeConfig("~/.gitconfig")
	if err != nil {
		return fmt.Errorf("Failed to open .gitconfig: %s\n", err)
	}

	return nil
}
