package cmds

import (
	"flag"
	"fmt"
	"os"
	"testing"

	"git-forge/forge/dummy"
	//"git-forge/log"
)

// This sets up our dummy forge driver and implicitly tests
// the addforge command
func TestMain(m *testing.M) {

	os.RemoveAll(".test")
	os.Mkdir(".test", 0755)
	f, serr := os.OpenFile("./.test/.gitconfig", os.O_CREATE, 0666)
	if serr != nil {
		fmt.Printf("Unable to setup test: %s\n", serr)
		os.Exit(1)
	}
	f.Close()

	cwd, wderr := os.Getwd()
	if wderr != nil {
		fmt.Printf("unable to set HOME directory, can't run tests\n")
		os.Exit(1)
	}

	os.Setenv("HOME", cwd+"/.test")
	os.Chdir(".test")

	// hand register the dummy forge driver
	rerr := RegisterForgeType("dummy", dummyforge.NewDummyForge)
	if rerr != nil {
		fmt.Printf("Unable to setup test: %s\n", rerr)
		os.Exit(1)
	}

	os.Args = []string{"initconfig"}
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	//logging.SuppressLog()
	terr := ForgeInitCmd()
	if terr != nil {
		fmt.Printf("Testing setup fails: %s\n", terr)
		os.Exit(1)
	}

	ret := m.Run()
	os.Chdir("../")
	os.RemoveAll("./.test")
	os.Exit(ret)
}

func TestCmds(t *testing.T) {
	var allcommandscomplete bool
	var skip bool

	t.Logf("Testing commands\n")
	// we run the tests until we've gone through them all
	// which may require several inner loop iterations
	// to satisfy dependencies
	for {
		allcommandscomplete = true
		for n, c := range Subcmds {
			skip = false

			// skip test that are complete
			if c.Testinfo.Tested == true {
				continue
			}
			if len(c.Testinfo.Testargs) == 0 {
				continue
			}
			// Make sure our test deps are resolved
			for _, d := range c.Testinfo.Deps {
				if Subcmds[d].Testinfo.Tested == false {
					allcommandscomplete = false
					skip = true
					break
				}
			}

			if skip == true {
				continue
			}

			t.Logf("Testing %s cmd\n", n)
			os.Args = c.Testinfo.Testargs
			flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
			terr := c.Cmd()
			if terr != nil {
				t.Errorf("%s Fails: %s\n", n, terr)
			}
			c.Testinfo.Tested = true
		}
		if allcommandscomplete == true {
			break
		}
	}
}
