package gitconfig

import (
	"fmt"
	"gopkg.in/ini.v1"
)

type forgeconfig interface {
	AddForge(name string, forgetype string, cloneUrl string, apiUrl string, user string, pass string) error
	DelForge(name string) error
	CommitConfig() error
}

type ForgeConfig struct {
	path string
	cfg  *ini.File
}

func NewForgeConfig(path string) (*ForgeConfig, error) {
	cfg, err := ini.Load(path)
	if err != nil {
		return nil, fmt.Errorf("Unable to open %s: %s\n", path, err)
	}
	return &ForgeConfig{
		path: path,
		cfg:  cfg,
	}, nil
}

func (f *ForgeConfig) AddForge(name string, forgetype string, cloneUrl string, apiUrl string, user string, pass string) error {
	var err error
	var sec *ini.Section

	sec, err = f.cfg.GetSection("forge \"" + name + "\"")
	if err == nil {
		// We found the section, so this is a duplicate
		return fmt.Errorf("Forge %s can't be created: %s\n", name, err)
	}

	sec, err = f.cfg.NewSection("forge \"" + name + "\"")
	if err != nil {
		return fmt.Errorf("Forge %s can't be created: %s\n", name, err)
	}

	// Now add our keys
	sec.NewKey("forgetype", forgetype)
	sec.NewKey("cloneurl", cloneUrl)
	sec.NewKey("apiurl", apiUrl)
	sec.NewKey("user", user)
	sec.NewKey("pass", pass)

	return nil
}

func (f *ForgeConfig) DelForge(name string) error {
	return nil
}

func (f *ForgeConfig) CommitConfig() error {
	return f.cfg.SaveTo(f.path)
}
