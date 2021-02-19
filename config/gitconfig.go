package gitconfig

import (
	"fmt"
	"gopkg.in/ini.v1"
	"os"
	"path/filepath"
	"strings"
)

type forgeconfig interface {
	AddForge(name string, forgetype string, cloneUrl string, apiUrl string, user string, pass string) error
	DelForge(name string) error
	LookupForgeType(url string) (string, error)
	LookupForgeName(url string) (string, error)
	GetCreds() (string, string, error)
	AddForgeRemoteSection(string, string, string) error
	GetForgeRemoteSection() (string, string, error)
	GetRemoteUrl(string) (string, error)
	CommitConfig() error
}

type ForgeConfig struct {
	path string
	cfg  *ini.File
	sec  *ini.Section
}

func findTopLevelGitDir(workingDir string) (string, error) {
	dir, err := filepath.Abs(workingDir)
	if err != nil {
		return "", err
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, ".git")); err == nil {
			return dir, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("No git tree found\n")
		}
		dir = parent
	}
	return "", fmt.Errorf("No Git tree found\n")
}

func LookupForgeType(url string) (string, error) {

	gitconfigpath := os.Getenv("HOME") + "/.gitconfig"

	forgeconfig, err := NewForgeConfig(gitconfigpath)
	if err != nil {
		return "", fmt.Errorf("Lookup forge config failed: %s\n", err)
	}

	return forgeconfig.LookupForgeType(url)
}

func LookupForgeName(url string) (string, error) {

	gitconfigpath := os.Getenv("HOME") + "/.gitconfig"

	forgeconfig, err := NewForgeConfig(gitconfigpath)
	if err != nil {
		return "", fmt.Errorf("Lookup forge config failed: %s\n", err)
	}

	return forgeconfig.LookupForgeName(url)
}

func NewForgeConfig(path string) (*ForgeConfig, error) {
	abspath, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("Unalbe to open %s: %s\n", abspath, err)
	}
	cfg, err := ini.Load(abspath)
	if err != nil {
		return nil, fmt.Errorf("Unable to open %s: %s\n", abspath, err)
	}
	return &ForgeConfig{
		path: abspath,
		cfg:  cfg,
		sec:  nil,
	}, nil
}

func NewLocalForgeConfig() (*ForgeConfig, error) {
	gitdir, err := findTopLevelGitDir("./")

	if err != nil {
		return nil, err
	}
	return NewForgeConfig(gitdir + "/.git/config")
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
		fname := strings.Trim(strings.SplitAfter(sec.Name(), " ")[1], "\"")
		if name == fname {
			return &ForgeConfig{
				path: os.Getenv("HOME") + "/.gitconfig",
				cfg:  cfg,
				sec:  sec,
			}, nil
		}
	}

	return nil, fmt.Errorf("Unable to find Forge named %s\n", name)

}

func GetForgeConfigFromUrl(path string, url string) (*ForgeConfig, error) {
	forge, err := NewForgeConfig(path)
	if err != nil {
		return nil, err
	}
	defer forge.CommitConfig()

	fname, err2 := forge.LookupForgeName(url)
	if err2 != nil {
		return nil, err2
	}
	return GetForgeConfig(path, fname)
}

func (f *ForgeConfig) AddForge(name string, forgetype string, cloneUrl string, apiUrl string, user string, pass string) error {
	var err error
	var sec *ini.Section

	sec, err = f.cfg.GetSection("forge \"" + name + "\"")
	if err == nil {
		// We found the section, so this is a duplicate
		return fmt.Errorf("Forge %s can't be created: Already Exists\n", name)
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

func (f *ForgeConfig) LookupForgeType(url string) (string, error) {
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

func (f *ForgeConfig) LookupForgeName(url string) (string, error) {
	secs := f.cfg.Sections()

	for _, sec := range secs {
		if sec.HasKey("forgetype") == false {
			continue
		}
		cloneurl := sec.Key("cloneurl").String()
		if strings.HasPrefix(url, cloneurl) == true {
			secname := sec.Name()
			secnameparts := strings.SplitAfter(secname, " ")
			return strings.Trim(secnameparts[1], "\""), nil
		}
	}

	return "", fmt.Errorf("Unable to locate forge for url %s\n", url)
}

func (f *ForgeConfig) DelForge() error {
	if f.sec == nil {
		return fmt.Errorf("No such Forge configured by that name\n")
	}
	name := strings.Trim(strings.SplitAfter(f.sec.Name(), " ")[1], "\"")
	fmt.Printf("Removing %s from %s\n", name, f.path)
	err := f.cfg.DeleteSectionWithIndex("forge \""+name+"\"", 0)
	if err != nil {
		return err
	}
	return nil
}

func (f *ForgeConfig) GetCreds() (string, string, error) {
	if f.sec == nil {
		return "", "", fmt.Errorf("No section specified for this config\n")
	}
	return f.sec.Key("user").String(), f.sec.Key("pass").String(), nil
}

func (f *ForgeConfig) AddForgeRemoteSection(forgetype string, child string, parent string) error {
	var sec *ini.Section
	var err error

	sec, err = f.cfg.GetSection("forge")
	if err != nil {
		sec, err = f.cfg.NewSection("forge")
		if err != nil {
			return fmt.Errorf("Unable to create a forge section in git configuration\n")
		}
	}

	sec.NewKey("forgetype", forgetype)
	sec.NewKey("childremote", child)
	sec.NewKey("parentremote", parent)

	return nil
}

func (f *ForgeConfig) GetForgeRemoteSection() (string, string, error) {
	sec, serr := f.cfg.GetSection("forge")
	if serr != nil {
		return "", "", serr
	}

	childremote, cerr := sec.GetKey("childremote")
	parentremote, perr := sec.GetKey("parentermote")
	if cerr != nil || perr != nil {
		return "", "", fmt.Errorf("Unable to get config keys for forge remotes\n")
	}

	return childremote.String(), parentremote.String(), nil
}

func (f *ForgeConfig) GetRemoteUrl(name string) (string, error) {
	sec, err := f.cfg.GetSection("remote \"" + name + "\"")
	if err != nil {
		return "", nil
	}
	url, uerr := sec.GetKey("url")
	if uerr != nil {
		return "", uerr
	}
	return url.String(), nil
}

func (f *ForgeConfig) CommitConfig() error {
	return f.cfg.SaveTo(f.path)
}
