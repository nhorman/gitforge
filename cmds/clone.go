package cmds

import (
	"flag"
	"git-forge/config"
	"git-forge/forge"
	"git-forge/log"
	"os"
)

func init() {
	RegisterCmd("clone", CloneForgeCmd, []string{"clone", "-getparent", "git@dummy.org:childtestuser/testrepo.git"})
}

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
	url := os.Args[len(os.Args)-1]

	fconfig, err := gitconfig.GetForgeConfigFromUrl(os.Getenv("HOME")+"/.gitconfig", url)
	if err != nil {
		return err
	}

	user, pass, cerr := fconfig.GetCreds()
	if cerr != nil {
		return cerr
	}

	fname, ferr := gitconfig.LookupForgeName(url)
	if ferr != nil {
		return ferr
	}
	cloneopts := forge.CloneOpts{
		Common: forge.CommonOpts{
			User: user,
			Pass: pass,
		},
		Parentfork: *parentopt,
		Url:        url,
		ForgeName:  fname,
	}

	forge, err2 := AllocateForgeFromUrl(url)
	if err != nil {
		return err2
	}

	clonerr := forge.Clone(cloneopts)

	return clonerr
}
