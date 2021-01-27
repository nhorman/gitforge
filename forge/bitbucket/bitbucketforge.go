package bitbucketforge

import (
	"git-forge/forge"
	"git-forge/log"
	//"github.com/ktrysmt/go-bitbucket"
	"gopkg.in/src-d/go-git.v4"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type BitBucketForge struct {
}

func NewBitBucketForge() *BitBucketForge {
	return &BitBucketForge{}

}

func (f *BitBucketForge) Clone(opts forge.CloneOpts) error {
	logging.Forgelog.Printf("%s appears to be a bitbucket forge\n", opts.Url)

	dirname := strings.TrimSuffix(path.Base(opts.Url), filepath.Ext(path.Base(opts.Url)))
	err := os.Mkdir("./"+dirname, 0755)
	if err != nil {
		return err
	}

	_, clonerr := git.PlainClone("./"+dirname, false, &git.CloneOptions{
		URL: opts.Url})
	if clonerr != nil {
		return clonerr
	}

	return nil
}

func (f *BitBucketForge) Fork(opts forge.ForkOpts) error {
	return nil
}
