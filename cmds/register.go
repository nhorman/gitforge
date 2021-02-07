package cmds

import ()

type CmdData struct {
	Cmd      func() error
	testargs []string
}

var Subcmds map[string]CmdData = make(map[string]CmdData, 0)

func RegisterCmd(cmd string, ifunc func() error, testargs []string) {
	Subcmds[cmd] = CmdData{ifunc, testargs}
}
