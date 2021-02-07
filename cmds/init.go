package cmds

import (
	"flag"
	"git-forge/config"
	"git-forge/log"
	"os"
)

func init() {
	RegisterCmd("initconfig", ForgeInitCmd, []string{"initconfig"})
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

	config, cerr := gitconfig.NewForgeConfig(os.Getenv("HOME") + "/.gitconfig")
	if cerr != nil {
		return cerr
	}

	defer config.CommitConfig()

	for k, f := range forgetypes {
		logging.Forgelog.Printf("Registering forge instances for %s type\n", k)
		forge := f()
		ferr := forge.InitForges(config)
		if ferr != nil {
			logging.Forgelog.Printf("Failed to configure %s\n", k)
		}
	}

	logging.Forgelog.Printf("Forges configured, make sure to edit your ~/.gitconfig file to add your username and password where appropriate\n")

	return nil
}
