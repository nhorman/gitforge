package cmds

import (
	"flag"
	"git-forge/forge"
	"git-forge/log"
	"os"
)

var forkDeps = TestData{[]string{"fork", "git@dummy.org:testuser/testrepo.git"}, []string{"initconfig"}, false}

func init() {
	RegisterCmd("fork", ForkForgeCmd, &forkDeps)
}

func Forkusage() {
	logging.Forgelog.Printf("Usage: git forge fork [options] <repo url>\n")
	logging.Forgelog.Printf("Description: fork a git repo on the forge to your private namespace\n")
	logging.Forgelog.Printf("Options:\n")
	flag.PrintDefaults()
}

func ForkForgeCmd() error {

	helpopt := flag.Bool("help", false, "display help for fork command")
	flag.Parse()

	if *helpopt == true {
		Cloneusage()
		return nil
	}
	// last argument must be the url
	url := os.Args[len(os.Args)-1]

	opts := forge.ForkOpts{
		Url: url,
	}
	forge, err := AllocateForgeFromUrl(url)
	if err != nil {
		return err
	}

	clonerr := forge.Fork(opts)

	return clonerr
}
