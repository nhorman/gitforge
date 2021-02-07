package cmds

import (
	"flag"
	"fmt"
	"os"
	"testing"

	_ "git-forge/forge/dummy"
	"git-forge/log"
)

// This sets up our dummy forge driver and implicitly tests
// the addforge command
func TestMain(m *testing.M) {
	logging.SuppressLog()
	terr := ForgeInitCmd()
	if terr != nil {
		fmt.Printf("Testing setup fails: %s\n", terr)
		os.Exit(1)
	}
}

func TestForkCmd(t *testing.T) {
	return
}

func TestCloneCmd(t *testing.T) {
	return
}

func TestDelForgeCmd(t *testing.T) {

	t.Log("Testing delforge cmd\n")
	// This should have been registered by the TestInitCmd test above
	os.Args = []string{"delforge", "-name", "dummy-ssh"}
	// Need to reset the CommandLine flag set
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	terr := DelForgeCmd()
	if terr != nil {
		t.Errorf("DelForgeCmd Fails: %s\n", terr)
	}
}
