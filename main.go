package main

import (
	"git-forge/cmds"
	"git-forge/log"
	"os"

	// imports for forge registrations
	_ "git-forge/forge/bitbucket"
)

func usage() error {
	logging.Forgelog.Printf("%s <cmd> [options]\n", os.Args[0])
	logging.Forgelog.Printf("cmds:\n")
	logging.Forgelog.Printf("\thelp\n")
	for key, _ := range cmds.Subcmds {
		logging.Forgelog.Printf("\t%s\n", key)
	}
	return nil
}

func main() {
	// get the subcommand, which will always be os.Aarg[1]
	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}

	if os.Args[1] == "help" {
		usage()
		os.Exit(0)
	}

	cmdname := os.Args[1]
	cmd, found := cmds.Subcmds[cmdname]
	if found != true {
		usage()
		os.Exit(1)
	}

	// Advance out argument list so that the stupid go flag parser doesnt
	// get all confused
	os.Args = os.Args[1:]
	err := cmd.Cmd()
	if err != nil {
		logging.Forgelog.Printf("%s failed: %s\n", cmdname, err)
		os.Exit(1)
	}
	os.Exit(0)
}
