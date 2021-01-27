package bitbucketforge

import (
	"fmt"
)

type BitBucketForge struct {
}

func NewBitBucketForge() *BitBucketForge {
	return &BitBucketForge{}

}

func (f *BitBucketForge) Clone(parentFork bool, url string) error {
	fmt.Printf("This appears to be a bitbucket forge\n")

	return nil
}
