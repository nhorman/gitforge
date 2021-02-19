package cmds

import (
	"flag"
	"git-forge/forge"
	"git-forge/log"
	"os"
)

var cloneDeps = TestData{[]string{"clone", "-getparent", "git@dummy.org:childtestuser/testrepo.git"}, []string{"initconfig"}, false}

func init() {
	RegisterCmd("clone", CloneForgeCmd, &cloneDeps)
}

func Cloneusage() {
	logging.Forgelog.Printf("Usage: git forge clone [options] <repo url>\n")
	logging.Forgelog.Printf("Description: clone a git repo, opionally setting up remotes for forks\n")
	logging.Forgelog.Printf("Options:\n")
	flag.PrintDefaults()
}

func CloneForgeCmd() error {

	helpopt := flag.Bool("help", false, "display help for clone command")
	parentopt := flag.Bool("getparent", false, "Find the parent of this repo, and add a remote for it")
	flag.Parse()

	if *helpopt == true {
		Cloneusage()
		return nil
	}
	// last argument must be the url
	url := os.Args[len(os.Args)-1]

	cloneopts := forge.CloneOpts{
		Parentfork: *parentopt,
		Url:        url,
	}

	forge, err2 := AllocateForgeFromUrl(url)
	if err2 != nil {
		return err2
	}

	clonerr := forge.Clone(cloneopts)

	return clonerr
}
