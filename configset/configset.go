package gitconfigset

import (
	"fmt"
	"gopkg.in/ini.v1"
	"os"
	"path/filepath"
	"strings"

	"git-forge/forge"
)

type ForgeConfigOps interface {
	AddForge(name string, forgetype string, cloneUrl string, apiUrl string, user string, pass string) error
	DelForge(name string) error
	AddForgeRemoteSection(string, string, string) error
	GetForgeRemoteSection() (string, string, string, error)
	GetRemoteUrl(string) (string, error)
	CommitConfig() error
}

type ForgeCfg struct {
	exists   bool
	modified bool
	path     string
	cfg      *ini.File
}

type ForgeConfigSet struct {
	Local  ForgeCfg
	Global ForgeCfg
}

func (f *ForgeConfigSet) AddForge(name string, forgetype string, cloneUrl string, apiUrl string, user string, pass string) error {
	var err error
	var sec *ini.Section

	sec, err = f.Global.cfg.GetSection("forge \"" + name + "\"")
	if err == nil {
		// We found the section, so this is a duplicate
		return fmt.Errorf("Forge %s can't be created: Already Exists\n", name)
	}

	sec, err = f.Global.cfg.NewSection("forge \"" + name + "\"")
	if err != nil {
		return fmt.Errorf("Forge %s can't be created: %s\n", name, err)
	}

	// Now add our keys
	sec.NewKey("forgetype", forgetype)
	sec.NewKey("cloneurl", cloneUrl)
	sec.NewKey("apiurl", apiUrl)
	sec.NewKey("user", user)
	sec.NewKey("pass", pass)
	f.Global.modified = true
	return nil
}

func (f *ForgeConfigSet) DelForge(name string) error {
	return f.Global.cfg.DeleteSectionWithIndex("forge \""+name+"\"", 0)
}

func (f *ForgeConfigSet) AddForgeRemoteSection(forgetype string, child string, parent string) error {
	var sec *ini.Section
	var err error

	sec, err = f.Local.cfg.GetSection("forge")
	if err != nil {
		sec, err = f.Local.cfg.NewSection("forge")
		if err != nil {
			return fmt.Errorf("Unable to create a forge section in git configuration\n")
		}
	}

	sec.NewKey("forgetype", forgetype)
	sec.NewKey("childremote", child)
	sec.NewKey("parentremote", parent)
	f.Local.modified = true
	return nil
}

func (f *ForgeConfigSet) GetForgeRemoteSection() (*forge.ForgeLocalConfig, error) {
	sec, serr := f.Local.cfg.GetSection("forge")
	if serr != nil {
		return nil, serr
	}

	childremote, cerr := sec.GetKey("childremote")
	parentremote, perr := sec.GetKey("parentremote")
	forgetype, terr := sec.GetKey("forgetype")
	if cerr != nil {
		return nil, fmt.Errorf("Unable to get config keys for child forge remotes: %s\n", cerr)
	}

	if perr != nil {
		return nil, fmt.Errorf("Unable to get config keys for parent forge remotes: %s\n", perr)
	}

	if terr != nil {
		return nil, fmt.Errorf("Unable to get config keys for forge type: %s\n", terr)
	}

	childurlsec, err1 := f.Local.cfg.GetSection("remote \"" + childremote.String() + "\"")
	if err1 != nil {
		return nil, err1
	}
	parenturlsec, err2 := f.Local.cfg.GetSection("remote \"" + parentremote.String() + "\"")
	if err2 != nil {
		return nil, err2
	}

	childurl, err3 := childurlsec.GetKey("url")
	parenturl, err4 := parenturlsec.GetKey("url")

	if err3 != nil || err4 != nil {
		return nil, fmt.Errorf("Unable to get urls for forge remotes\n")
	}

	ret := &forge.ForgeLocalConfig{
		Type: forgetype.String(),
		Child: forge.ForgeRemoteInfo{
			Name: childremote.String(),
			Url:  childurl.String(),
		},
		Parent: forge.ForgeRemoteInfo{
			Name: parentremote.String(),
			Url:  parenturl.String(),
		},
	}
	return ret, nil
}

func (f *ForgeConfigSet) ConfigFromUrl(url string) (*forge.ForgeConfig, error) {
	secs := f.Global.cfg.Sections()

	for _, sec := range secs {
		if sec.HasKey("forgetype") == false {
			continue
		}
		cloneurl := sec.Key("cloneurl").String()
		if strings.HasPrefix(url, cloneurl) == true {
			cfg := &forge.ForgeConfig{
				Type:         sec.Key("forgetype").String(),
				User:         sec.Key("user").String(),
				Pass:         sec.Key("pass").String(),
				CloneBaseUrl: sec.Key("cloneurl").String(),
				ApiBaseUrl:   sec.Key("apiurl").String(),
			}
			return cfg, nil
		}
	}
	return nil, fmt.Errorf("No forge config for url %s\n", url)
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

func loadConfigs(c *ForgeConfigSet) error {
	var err error
	if c.Local.exists == true {
		c.Local.cfg, err = ini.Load(c.Local.path)
		if err != nil {
			c.Local.exists = false
			return err
		}
	}
	if c.Global.exists == true {
		c.Global.cfg, err = ini.Load(c.Global.path)
		if err != nil {
			c.Global.exists = false
			return nil
		}
	}
	return nil
}

func NewForgeConfigSet() (*ForgeConfigSet, error) {

	gpath, err := filepath.Abs(os.Getenv("HOME") + "/.gitconfig")
	if err != nil {
		return nil, err
	}

	lpath, err := findTopLevelGitDir(".")

	global := ForgeCfg{false, false, gpath, nil}
	global.exists = true

	local := ForgeCfg{false, false, lpath + "/.git/config", nil}
	if err == nil {
		local.exists = true
	}

	fcs := &ForgeConfigSet{local, global}
	err = loadConfigs(fcs)
	if err != nil {
		return nil, err
	}
	return fcs, nil

}

func NewForgeConfigSetInDir(dirname string) (*ForgeConfigSet, error) {
	cfg, err := NewForgeConfigSet()
	if err != nil {
		return nil, err
	}

	newlocalpath, err2 := filepath.Abs("./" + dirname + "/.git/config")
	if err2 != nil {
		return nil, err
	}

	cfg.Local = ForgeCfg{false, false, newlocalpath, nil}
	cfg.Local.exists = true
	ini, err2 := ini.Load(cfg.Local.path)
	if err2 != nil {
		cfg.Local.exists = false
		return nil, err2
	}
	cfg.Local.cfg = ini
	return cfg, nil
}

func (f *ForgeConfigSet) CommitConfig() error {
	var err1, err2 error
	if f.Local.modified == true {
		err1 = f.Local.cfg.SaveTo(f.Local.path)
	} else {
		err1 = nil
	}

	if f.Global.modified == true {
		err2 = f.Global.cfg.SaveTo(f.Global.path)
	} else {
		err2 = nil
	}

	if err1 != nil {
		return err1
	}
	if err2 != nil {
		return err2
	}
	return nil
}
