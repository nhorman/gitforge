package cmds

import ()

type TestData struct {
	Testargs []string
	Deps     []string
	Tested   bool
}

type CmdData struct {
	Cmd      func() error
	Testinfo *TestData
}

var Subcmds map[string]CmdData = make(map[string]CmdData, 0)

func RegisterCmd(cmd string, ifunc func() error, testdata *TestData) {
	Subcmds[cmd] = CmdData{ifunc, testdata}
}
