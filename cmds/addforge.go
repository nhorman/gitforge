package cmds

import (
	"flag"
	"fmt"
	"git-forge/configset"
	"git-forge/log"
)

var addForgeDeps = TestData{[]string{"addforge", "-name", "dummy-http"}, []string{"initconfig"}, false}

func init() {
	RegisterCmd("addforge", AddForgeCmd, &addForgeDeps)
}

func Addusage() {
	logging.Forgelog.Printf("Usage: git forge addforge [options]\n")
	logging.Forgelog.Printf("Description: Add a forge type to the global gitconfig\n")
	logging.Forgelog.Printf("Options:\n")
	flag.PrintDefaults()
	PrintForgeTypes()
}

func AddForgeCmd() error {

	helpopt := flag.Bool("help", false, "display help for addforge command")
	nameopt := flag.String("name", "", "Name of the forge to add")
	typeopt := flag.String("type", "", "Type of forge (bitbucket, github, etc)")
	clone := flag.String("cloneurl", "", "Base url this forge clones from")
	api := flag.String("apiurl", "", "Base url this forge uses to access REST apis")
	user := flag.String("user", "", "User name this forge uses for personal access")
	pass := flag.String("pass", "", "Password used to access this forge")

	flag.Parse()

	if *helpopt == true {
		Addusage()
		return nil
	}

	forgeconfig, err := gitconfigset.NewForgeConfigSet()
	if err != nil {
		return fmt.Errorf("Create forge config failed: %s\n", err)
	}
	defer forgeconfig.CommitConfig()

	ferr := forgeconfig.AddForge(*nameopt, *typeopt, *clone, *api, *user, *pass)
	if ferr != nil {
		return fmt.Errorf("Failed to add forge: %s\n", ferr)
	}

	return nil
}
