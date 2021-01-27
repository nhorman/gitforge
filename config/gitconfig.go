package gitconfig

import (
	"fmt"
	"gopkg.in/ini.v1"
	"strings"
)

type forgeconfig interface {
	AddForge(name string, forgetype string, cloneUrl string, apiUrl string, user string, pass string) error
	DelForge(name string) error
	LookupForge(url string) (string, error)
	CommitConfig() error
}

type ForgeConfig struct {
	path string
	cfg  *ini.File
	sec  *ini.Section
}

func NewForgeConfig(path string) (*ForgeConfig, error) {
	cfg, err := ini.Load(path)
	if err != nil {
		return nil, fmt.Errorf("Unable to open %s: %s\n", path, err)
	}
	return &ForgeConfig{
		path: path,
		cfg:  cfg,
		sec:  nil,
	}, nil
}

func GetForgeConfig(path string, name string) (*ForgeConfig, error) {
	cfg, err := ini.Load(path)
	if err != nil {
		return nil, fmt.Errorf("Unable to open ~/.gitconfig: %s\n", err)
	}

	secs := cfg.Sections()

	for _, sec := range secs {
		if sec.HasKey("forgetype") == false {
			continue
		}
		fname := sec.Key("name").String()
		if name == fname {
			return &ForgeConfig{
				path: "~/.gitconfig",
				cfg:  cfg,
				sec:  sec,
			}, nil
		}
	}

	return nil, fmt.Errorf("Unable to find Forge named %s\n", name)

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
	f.sec = sec
	return nil
}

func (f *ForgeConfig) LookupForge(url string) (string, error) {
	secs := f.cfg.Sections()

	for _, sec := range secs {
		if sec.HasKey("forgetype") == false {
			continue
		}
		cloneurl := sec.Key("cloneurl").String()
		if strings.HasPrefix(url, cloneurl) == true {
			return sec.Key("forgetype").String(), nil
		}
	}

	return "", fmt.Errorf("Unable to locate forge for url %s\n", url)
}

func (f *ForgeConfig) DelForge() error {
	name := f.sec.Key("name").String()
	f.cfg.DeleteSection("forge \"" + name + "\"")
	return nil
}

func (f *ForgeConfig) CommitConfig() error {
	return f.cfg.SaveTo(f.path)
}
