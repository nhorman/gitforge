package cmds

import (
	"flag"
	//"git-forge/config"
	//"git-forge/configset"
	"fmt"
	"git-forge/forge"
	"git-forge/log"
	//"os"
)

var createPrDeps = TestData{[]string{"createpr", "-sbranch", "testbranch", "-tbranch", "master", "--title", "Test Merge"}, []string{"clone"}, false}

func init() {
	RegisterCmd("createpr", CreatePrForgeCmd, &createPrDeps)
}

func CreatePrusage() {
	logging.Forgelog.Printf("Usage: git forge createpr [options]\n")
	logging.Forgelog.Printf("Description: create a pull request from the specified branch of your repo to the parent repo target branch\n")
	logging.Forgelog.Printf("Options:\n")
	flag.PrintDefaults()
}

func CreatePrForgeCmd() error {

	helpopt := flag.Bool("help", false, "display help for createpr command")
	sbranch := flag.String("sbranch", "", "source branch to request merge for")
	tbranch := flag.String("tbranch", "", "target branch to merge to")
	title := flag.String("title", "", "title of PR")
	description := flag.String("description", "", "description of pr")

	flag.Parse()

	if *helpopt == true {
		CreatePrusage()
		return nil
	}

	if *sbranch == "" || *tbranch == "" {
		return fmt.Errorf("Both source branch and target branch options are requried\n")
	}

	if *title == "" {
		return fmt.Errorf("PR Requires a Title\n")
	}

	myforge, err := AllocateForge()
	if err != nil {
		return fmt.Errorf("Unable to allocate forge: %s\n", err)
	}

	propts := forge.CreatePrOpts{
		Sbranch:     *sbranch,
		Tbranch:     *tbranch,
		Title:       *title,
		Description: *description,
	}

	prerr := myforge.CreatePr(propts)
	if prerr != nil {
		return prerr
	}

	return nil
}
