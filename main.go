package main

import (
	"fmt"
	"git-forge/cmds"
	"os"
)

var subcmds = map[string]func() error{
	"addforge": cmds.AddForgeCmd,
}

func usage() error {
	fmt.Printf("%s <cmd> [options]\n", os.Args[0])
	fmt.Printf("cmds:\n")
	fmt.Printf("\thelp\n")
	for key, _ := range subcmds {
		fmt.Printf("\t%s\n", key)
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
	cmd, found := subcmds[cmdname]
	if found != true {
		usage()
		os.Exit(1)
	}

	// Advance out argument list so that the stupid go flag parser doesnt
	// get all confused
	os.Args = os.Args[1:]
	err := cmd()
	if err != nil {
		fmt.Printf("%s failed: %s\n", cmdname, err)
		os.Exit(1)
	}
	os.Exit(0)
}
