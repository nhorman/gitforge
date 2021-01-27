package bitbucketforge

import (
	"git-forge/log"
)

type BitBucketForge struct {
}

func NewBitBucketForge() *BitBucketForge {
	return &BitBucketForge{}

}

func (f *BitBucketForge) Clone(parentFork bool, url string) error {
	logging.Forgelog.Printf("%s appears to be a bitbucket forge\n", url)

	return nil
}
