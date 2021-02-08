package cmds

import (
	"flag"
	"git-forge/config"
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

	fconfig, err := gitconfig.GetForgeConfigFromUrl(os.Getenv("HOME")+"/.gitconfig", url)
	if err != nil {
		return err
	}

	user, pass, cerr := fconfig.GetCreds()
	if cerr != nil {
		return cerr
	}

	opts := forge.ForkOpts{
		Common: forge.CommonOpts{
			User: user,
			Pass: pass,
		},
		Url: url,
	}
	forge, err := AllocateForgeFromUrl(url)
	if err != nil {
		return err
	}

	clonerr := forge.Fork(opts)

	return clonerr
}
