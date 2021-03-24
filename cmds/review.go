package cmds

import (
	"flag"
	//"git-forge/config"
	//"git-forge/configset"
	"fmt"
	"git-forge/forge"
	"git-forge/log"
	"git-forge/ui"
	//"os"
)

var ReviewDeps = TestData{[]string{"review"}, []string{"createpr"}, false}

func init() {
	RegisterCmd("review", ReviewCmd, &ReviewDeps)
}

func ReviewUsage() {
	logging.Forgelog.Printf("Usage: git forge review\n")
	logging.Forgelog.Printf("Interactively view PRs and Issues\n")
	logging.Forgelog.Printf("Options:\n")
	flag.PrintDefaults()
}

func ReviewCmd() error {

	helpopt := flag.Bool("help", false, "display help for createpr command")
	logfileopt := flag.String("logfile", "", "Redirect log messages to file")

	flag.Parse()

	if *helpopt == true {
		CreatePrusage()
		return nil
	}

	if *logfileopt != "" {
		logging.LogToFile(*logfileopt)
	}

	myforge, err := AllocateForge()
	if err != nil {
		return fmt.Errorf("Unable to allocate forge: %s\n", err)
	}

	return forgeui.RunUi(myforge.(forge.ForgeUIModel))
}
