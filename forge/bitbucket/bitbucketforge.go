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
	logging.Forgelog.Printf("This appears to be a bitbucket forge\n")

	return nil
}
