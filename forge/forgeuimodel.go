package forge

import ()

type PrTitle struct {
	Title string
	PrId  int64
}

//
// Define the structs that represent a
// Pull Request
//
type PrRemote struct {
	URL        string `json:"url"`
	BranchName string `json:"branch"`
}

type PrSpec struct {
	Source PrRemote `json:"source"`
	Target PrRemote `json:"target"`
}

type DiscussionType int

const (
	GENERAL DiscussionType = iota
	INLINE  DiscussionType = iota
)

type InlineCommentData struct {
	Path   string `json:"path"`
	Offset int    `json:"offset"`
}

type Discussion struct {
	Id       int               `json:"id"`
	ParentId int               `json:"parentid"`
	Type     DiscussionType    `json:"type"`
	Inline   InlineCommentData `json:"inlinedata"`
	Author   string            `json:"author"`
	Content  string            `json::"content"`
}

type PR struct {
	PrId        int64        `json:"prid"`
	Title       string       `json:"title"`
	PullSpec    PrSpec       `json:"prspec"`
	Discussions []Discussion `json:"discussions"`
}

type ForgeUIModel interface {
	GetAllPrTitles() ([]PrTitle, error)
	GetPr(idstring string) (*PR, error)
}
