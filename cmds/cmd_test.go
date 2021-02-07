package cmds

import (
	"flag"
	"fmt"
	"os"
	"testing"

	_ "git-forge/forge/dummy"
	//"git-forge/log"
)

// This sets up our dummy forge driver and implicitly tests
// the addforge command
func TestMain(m *testing.M) {
	//logging.SuppressLog()
	terr := ForgeInitCmd()
	if terr != nil {
		fmt.Printf("Testing setup fails: %s\n", terr)
		os.Exit(1)
	}
}

func TestCmds(t *testing.T) {

	cwd, wderr := os.Getwd()
	if wderr != nil {
		t.Errorf("unable to set HOME directory, can't run tests\n")
	}

	os.Setenv("HOME", cwd)
	t.Log("Testing commands\n")
	for n, c := range Subcmds {
		t.Logf("Testing %s cmd\n", n)
		os.Args = c.testargs
		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
		terr := c.Cmd()
		if terr != nil {
			t.Errorf("%s Fails: %s\n", n, terr)
		}
	}
}
