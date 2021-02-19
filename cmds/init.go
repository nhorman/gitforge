package cmds

import (
	"flag"
	"git-forge/log"
)

var initConfigDeps = TestData{[]string{}, []string{}, true}

func init() {
	RegisterCmd("initconfig", ForgeInitCmd, &initConfigDeps)
}

func initusage() {
	logging.Forgelog.Printf("Usage: git forge init\n")
	logging.Forgelog.Printf("Description: initalize global git config with standard forge instances\n")
}

func ForgeInitCmd() error {

	helpopt := flag.Bool("help", false, "display help for fork command")
	flag.Parse()

	if *helpopt == true {
		initusage()
		return nil
	}

	for k, f := range forgetypes {
		logging.Forgelog.Printf("Registering forge instances for %s type\n", k)
		forge := f(nil)
		ferr := forge.InitForges()
		if ferr != nil {
			logging.Forgelog.Printf("Failed to configure %s: %s\n", k, ferr)
		}
	}

	logging.Forgelog.Printf("Forges configured, make sure to edit your ~/.gitconfig file to add your username and password where appropriate\n")

	return nil
}
