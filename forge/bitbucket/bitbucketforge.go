package bitbucketforge

import (
	"git-forge/log"
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

func (f *BitBucketForge) Clone(parentFork bool, url string) error {
	logging.Forgelog.Printf("%s appears to be a bitbucket forge\n", url)

	dirname := strings.TrimSuffix(path.Base(url), filepath.Ext(path.Base(url)))
	err := os.Mkdir("./"+dirname, 0755)
	if err != nil {
		return err
	}

	_, clonerr := git.PlainClone("./"+dirname, false, &git.CloneOptions{
		URL: url})
	if clonerr != nil {
		return clonerr
	}

	return nil
}
