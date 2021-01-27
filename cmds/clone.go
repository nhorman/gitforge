package cmds

import (
	"flag"
	"git-forge/log"
	"os"
)

func Cloneusage() {
	logging.Forgelog.Printf("Usage: git forge clone [options] <repo url>\n")
	logging.Forgelog.Printf("Description: clone a git repo, opionally setting up remotes for forks\n")
	logging.Forgelog.Printf("Options:\n")
	flag.PrintDefaults()
}

func CloneForgeCmd() error {

	helpopt := flag.Bool("help", false, "display help for addforge command")
	parentopt := flag.Bool("getparent", false, "Find the parent of this repo, and add a remote for it")
	flag.Parse()

	if *helpopt == true {
		Cloneusage()
		return nil
	}
	// last argument must be the url
	url := os.Args[len(os.Args)]

	forge, err := AllocateForgeFromUrl(url)
	if err != nil {
		return err
	}

	clonerr := forge.Clone(*parentopt, url)

	return clonerr
}
