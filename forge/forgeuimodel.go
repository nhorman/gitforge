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

type CommitCommentData struct {
	Author  string `json:"author"`
	Content string `json:"content"`
	Path    string `json:"path"`
	Offset  int    `json:"offset"`
}

type Commit struct {
	Hash     string              `json:"hash"`
	Comments []CommitCommentData `json:"comments"`
}

type PR struct {
	Unread       bool         `json:"unread"`
	PrId         int64        `json:"prid"`
	CurrentToken string       `json:"currenttoken"`
	Title        string       `json:"title"`
	PullSpec     PrSpec       `json:"prspec"`
	Discussions  []Discussion `json:"discussions"`
	Commits      []Commit     `json:"commits"`
}

type UpdateResult int

const (
	UPDATE_CURRENT  = iota //Means that the current cached pr is up to date
	UPDATE_REPULL   = iota //Means that the current cached pr is being updated by the model
	UPDATE_FINISHED = iota //Means that the model is done with all updated (this is always sent in the non failure case)
	UPDATE_FAILED   = iota //Means that the update failed
)

type UpdatedPR struct {
	Result UpdateResult
	Pr     *PR
}

type ForgeUIModel interface {
	GetAllPrTitles() ([]PrTitle, error)
	GetPr(idstring string) (*PR, error)
	RefreshPr(pr *PR) (chan *UpdatedPR, error)
}
