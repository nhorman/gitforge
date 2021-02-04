package cmds

import ()

var Subcmds map[string]func() error = make(map[string]func() error, 0)

func RegisterCmd(cmd string, ifunc func() error) {
	Subcmds[cmd] = ifunc
}
