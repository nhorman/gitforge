package gitconfig

import ()

type ForgeType int

const (
	BITBUCKET ForgeType = iota
	MAX                 = iota
)

type forgeconfig interface {
	AddForge(name string, forgetype ForgeType, cloneUrl string, apiUrl string, user string, pass string) error
	DelForge(name string) error
}

type ForgeConfig struct {
	path string
}

func NewForgeConfig(path string) (*ForgeConfig, error) {
	return &ForgeConfig{
		path: path,
	}, nil
}

func (*ForgeConfig) AddForge(name string, forgetype ForgeType, cloneUrl string, apiUrl string, user string, pass string) error {
	return nil
}

func (ForgeConfig) DelForge(name string) error {
	return nil
}
