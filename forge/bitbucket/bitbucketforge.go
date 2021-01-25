package bitbucketforge

import ()

type BitBucketForge struct {
}

func NewBitBucketForge() *BitBucketForge {
	return &BitBucketForge{}

}

func (f *BitBucketForge) Clone(createFork bool, attachFork bool, url string) error {
	return nil
}
